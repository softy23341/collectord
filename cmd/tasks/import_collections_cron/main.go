package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	cleaver "git.softndit.com/collector/backend/cleaver/client"
	"git.softndit.com/collector/backend/config"
	"git.softndit.com/collector/backend/dal"
	dalpg "git.softndit.com/collector/backend/dal/pg"
	"git.softndit.com/collector/backend/dto"
	resource "git.softndit.com/collector/backend/resources"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/jackc/pgx"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli"
	"gopkg.in/inconshreveable/log15.v2"
)

var cliFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config",
		Value: "collectord.toml",
	},
	cli.StringFlag{
		Name:  "email",
		Value: "sir_winston_churchill@gmail.com",
	},
}

type ImporterContext struct {
	DBM           dal.TrManager
	ProdDBM       dal.Manager
	Logger        log15.Logger
	CleaverClient cleaver.ConnectClient
	FileStorage   services.FileStorage
}

func main() {
	app := cli.NewApp()
	app.Flags = cliFlags

	app.Action = func(c *cli.Context) {
		configPath := c.String("config")

		mainConfig := &config.Config{}
		if err := mainConfig.ReadConfig(configPath); err != nil {
			panic(err)
		}

		instanceHolder, err := config.NewInstanceHolder(mainConfig)
		if err != nil {
			panic(err)
		}

		serverConfig := mainConfig.ServerConfig

		//logger
		logger := log15.New("app", "import collections")

		// get prod db connection
		prodDBM, err := dalpg.NewManager(&dalpg.ManagerContext{
			Log: logger,
			PoolConfig: &pgx.ConnPoolConfig{
				ConnConfig: pgx.ConnConfig{
					Host:     "localhost",
					Port:     5434,
					Database: "collector_prod",
					User:     "collector_app",
					Password: "collectordevpassstrong",
				},
				MaxConnections: 10,
			},
		})
		if err != nil {
			panic(err)
		}

		// get db connection
		pgConfig, err := instanceHolder.GetPgConfig(serverConfig.DBM.PgRef)
		if err != nil {
			panic(err)
		}
		DBM, err := dalpg.NewManager(&dalpg.ManagerContext{
			PoolConfig: pgConfig,
		})
		if err != nil {
			panic(err)
		}

		// cleaver
		rabbitConfig, err := instanceHolder.GetRabbitConfig(serverConfig.CleaverClient.RabbitRef)
		if err != nil {
			panic(err)
		}
		cleaverClient := cleaver.NewRMQClient(rabbitConfig.URL)
		if err := cleaverClient.Connect(); err != nil {
			panic(err)
		}

		// file storage
		filerCfg, err := instanceHolder.GetGenericFilerConfig(serverConfig.FilerClient.FilerRef)
		if err != nil {
			panic(err)
		}
		builder, err := services.GetTomlFilerBuilder(filerCfg.Type)
		if err != nil {
			panic(err)
		}
		storage, err := builder(&services.TomlFileStorageContext{
			Log:    logger.New("service", "filer"),
			Config: filerCfg.Config,
		})
		if err != nil {
			panic(err)
		}

		rand.Seed(time.Now().Unix())

		ci := &ImporterContext{
			DBM:           DBM,
			ProdDBM:       prodDBM,
			Logger:        logger,
			FileStorage:   storage,
			CleaverClient: cleaverClient,
		}

		//init job
		gocron.Every(1).Days().Do(func() {
			ci.runImportCollections(c.String("email"))
			ci.runImportObjects(c.String("email"))
		})

		//start all the pending jobs
		<-gocron.Start()
	}

	app.Run(os.Args)
}

func containsCollection(s dto.CollectionList, e *dto.Collection) bool {
	for _, a := range s {
		if *a.UserUniqID == *e.UserUniqID {
			return true
		}
	}
	return false
}

func containsObjectStatus(s dto.ObjectStatusList, e *dto.ObjectStatus) (*dto.ObjectStatus, bool) {
	for _, a := range s {
		if a.Name == e.Name {
			return a, true
		}
	}
	return nil, false
}

func getUserByEmail(dbm dal.Manager, email string) (*dto.User, error) {
	user, err := dbm.GetUserByEmailWithCustomFields(email, dto.UserFieldsList{
		dto.UserFieldID,
		dto.UserFieldFirstName,
		dto.UserFieldLastName,
		dto.UserFieldEmail,
	})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("no users found")
	}
	return user, nil
}

func (ci *ImporterContext) runImportCollections(email string) {
	ci.Logger.Info("run collections import", "at", time.Now().Format(time.RFC3339))

	produser, err := getUserByEmail(ci.ProdDBM, email)
	if err != nil {
		ci.Logger.Error("collections import", "get prod user", err.Error())
		return
	}

	user, err := getUserByEmail(ci.DBM, email)
	if err != nil {
		ci.Logger.Error("collections import", "get user error", err.Error())
		return
	}

	var root *dto.Root
	if roots, err := ci.DBM.GetRootsByUserID(user.ID); err != nil {
		ci.Logger.Error("collections import", "get root", err.Error())
		return
	} else if len(roots) != 1 {
		ci.Logger.Error("collections import", "get root", "no root found")
		return
	} else {
		root = roots[0]
	}

	prodCollectionsList, err := ci.ProdDBM.GetCollectionsByUserIDsWithCustomFields([]int64{produser.ID},
		dto.CollectionFieldsList{
			dto.CollectionFieldID,
			dto.CollectionFieldName,
			dto.CollectionFieldRootID,
			dto.CollectionFieldUserID,
			dto.CollectionFieldUserUniqID,
			dto.CollectionFieldCreationTime,
			dto.CollectionFieldTypo,
			dto.CollectionFieldDescription,
			dto.CollectionFieldImageMediaID,
		})
	if err != nil {
		ci.Logger.Error("collections import", "get prod collections", err.Error())
		return
	}

	collectionsList, err := ci.DBM.GetCollectionsByUserIDsWithCustomFields([]int64{user.ID}, dto.CollectionAllFields)
	if err != nil {
		ci.Logger.Error("collections import", "get collections", err.Error())
		return
	}

	ci.Logger.Info("collections_list", "length", len(collectionsList))
	ci.Logger.Info("prod_collections_list", "length", len(prodCollectionsList))

	var newCollections dto.CollectionList
	for _, pd := range prodCollectionsList {
		if !containsCollection(collectionsList, pd) {
			newCollections = append(newCollections, pd)
		}
	}

	ci.Logger.Info("selected collections", "length", len(newCollections))

	for _, collection := range newCollections {
		// update new collection data;
		collection.CreationTime = time.Now()
		collection.UserID = &user.ID
		collection.RootID = root.ID
		// save collection
		if err := ci.saveCollection(root.ID, collection); err != nil {
			ci.Logger.Error("collections import", "save collection", err.Error())
			return
		} else {
			ci.Logger.Info("new collection created", "ok", collection.ID)
		}
	}

	ci.Logger.Info("stop collections import", "at", time.Now().Format(time.RFC3339))
}

func (ci *ImporterContext) runImportObjects(email string) {
	ci.Logger.Info("run objects import", "at", time.Now().Format(time.RFC3339))

	produser, err := getUserByEmail(ci.ProdDBM, email)
	if err != nil {
		ci.Logger.Error("objects import", "get prod user", err.Error())
		return
	}

	user, err := getUserByEmail(ci.DBM, email)
	if err != nil {
		ci.Logger.Error("objects import", "get user", err.Error())
		return
	}

	var root *dto.Root
	if roots, err := ci.DBM.GetRootsByUserID(user.ID); err != nil {
		ci.Logger.Error("collections import", "get root", err.Error())
		return
	} else if len(roots) != 1 {
		ci.Logger.Error("collections import", "get root", "no root found")
		return
	} else {
		root = roots[0]
	}

	prodObjectsList, err := ci.ProdDBM.GetObjectsByUserIDsWithCustomFields([]int64{produser.ID}, dto.ObjectAllFields)
	if err != nil {
		ci.Logger.Error("objects import", "get objects", err.Error())
		return
	}

	for _, object := range prodObjectsList {
		// check is object exist
		if objectID, _ := ci.DBM.GetObjectIDByUserUniqID(user.ID, object.UserUniqID); objectID == nil {
			// update new user;
			object.UserID = user.ID
			// update new date interval id;
			if object.ProductionDateIntervalID != nil {
				if nid, err := ci.getProductionDateIntervalID(*object.ProductionDateIntervalID, root.ID); err == nil {
					object.ProductionDateIntervalID = nid
				} else {
					object.ProductionDateIntervalID = nil
				}
			}
			// get new collection id;
			if nid, err := ci.getCollectionID(user.ID, object.CollectionID); err == nil || nid != nil {
				object.CollectionID = *nid
				// save collection
				if err := ci.saveObject(root.ID, object); err != nil {
					ci.Logger.Info("objects import", "save object", err.Error())
					return
				} else {
					ci.Logger.Info("objects import", "new object", object.ID)
				}
			} else {
				ci.Logger.Error("objects import", "get new collection", err)
			}
		}
	}

	ci.Logger.Info("stop objects import", "at", time.Now().Format(time.RFC3339))
}

func (ci *ImporterContext) getCollectionID(userId, prodCollectionID int64) (*int64, error) {
	var prodCollection *dto.Collection
	if collections, err := ci.ProdDBM.GetCollectionsByIDsWithCustomFields([]int64{prodCollectionID}, dto.CollectionFieldsList{
		dto.CollectionFieldID, dto.CollectionFieldUserUniqID}); err != nil {
		return nil, err
	} else if len(collections) != 1 || collections[0].UserUniqID == nil {
		return nil, errors.New(fmt.Sprintf("no collections found: prod collection_id: %v", prodCollectionID))
	} else {
		prodCollection = collections[0]
	}
	return ci.DBM.GetCollectionIDByUserUniqID(userId, *prodCollection.UserUniqID)
}

func (ci *ImporterContext) getProductionDateIntervalID(prodDateIntervalID int64, newRootID int64) (*int64, error) {
	prodIntervals, err := ci.ProdDBM.GetNamedDateIntervalsByIDs([]int64{prodDateIntervalID})
	if err != nil {
		return nil, err
	}
	if len(prodIntervals) != 1 {
		return nil, errors.New("no intervals found")
	}
	prodInterval := prodIntervals[0]

	intervals, err := ci.DBM.GetNamedDayeIntervalsByNormalNames(newRootID, []string{prodInterval.NormalName})
	if err != nil {
		return nil, err
	}
	if len(intervals) == 1 {
		interval := intervals[0]
		return &interval.ID, nil
	} else {
		prodInterval.RootID = &newRootID
		if err := ci.DBM.CreateNamedDateInterval(prodInterval); err != nil {
			return nil, err
		}
		return &prodInterval.ID, nil
	}
}

func (ci *ImporterContext) getObjectStatusRefID(prodStatusID int64, userID int64) (*int64, error) {
	prodStatuses, err := ci.ProdDBM.GetObjectStatusByIDs([]int64{prodStatusID})
	if err != nil {
		return nil, err
	}
	if len(prodStatuses) != 1 {
		return nil, errors.New("no statuses found")
	}
	prodStatus := prodStatuses[0]
	statuses, err := ci.DBM.GetObjectStatuses()
	if err != nil {
		return nil, err
	}
	if o, found := containsObjectStatus(statuses, prodStatus); found {
		return &o.ID, nil
	}
	return nil, nil
}

func (ci *ImporterContext) saveCollection(rootID int64, collection *dto.Collection) error {
	tx, err := ci.DBM.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// save collection image;
	if collection.ImageMediaID != nil {
		if collectionMedias, _ := ci.ProdDBM.GetMediasByIDs([]int64{*collection.ImageMediaID}); len(collectionMedias) == 1 {

			collectionMedia := collectionMedias[0]
			collectionMedia.UserID = collection.UserID
			collectionMedia.UserUniqID = NewInt64(util.NextUniqID())

			if media, err := ci.saveMedia(tx, collectionMedia); err != nil {
				ci.Logger.Error("collections import", "save media", err.Error())
				return err
			} else if media == nil {
				collection.ImageMediaID = nil
			} else {
				collection.ImageMediaID = &media.ID
			}
		}
	}

	// save collection
	if err := tx.CreateCollection(collection); err != nil {
		return err
	}

	// set rights
	right := &dto.UserEntityRight{
		UserID:     *collection.UserID,
		EntityType: dto.RightEntityTypeCollection,
		EntityID:   collection.ID,
		Level:      dto.RightEntityLevelAdmin,
		RootID:     rootID,
	}
	if err := tx.PutUserRight(right); err != nil {
		return err
	}

	// commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ci *ImporterContext) saveObject(rootID int64, object *dto.Object) error {
	prodObjectID := object.ID

	tx, err := ci.DBM.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// save object;
	if err := tx.CreateObject(object); err != nil {
		return err
	}

	// save object medias;
	var objectMediasIDs []int64
	if prodMediaRefs, err := ci.ProdDBM.GetObjectsMediaRefs([]int64{prodObjectID}); err != nil {
		return err
	} else if len(prodMediaRefs) != 0 {
		if medias, err := ci.ProdDBM.GetMediasByIDs(prodMediaRefs.OrderedMediasIDs()); err != nil {
			return err
		} else {
			for _, media := range medias {
				media.UserID = &object.UserID
				media.RootID = &rootID
				if media, err := ci.saveMedia(tx, media); err != nil {
					ci.Logger.Error("objects import", "save media", err.Error(), "media_id", media.ID)
					return err
				} else if media != nil {
					objectMediasIDs = append(objectMediasIDs, media.ID)
				}
			}
		}
	}

	// save object badges;
	var objectBadgesIDs []int64
	if prodObjectBadges, err := ci.ProdDBM.GetObjectsBadgeRefs([]int64{prodObjectID}); err != nil {
		return err
	} else if len(prodObjectBadges) != 0 {
		if objectBadges, err := ci.ProdDBM.GetBadgesByIDs(prodObjectBadges.BadgesIDs()); err != nil {
			return err
		} else {
			for _, b := range objectBadges {
				b.RootID = &rootID
				if badge, err := tx.GetOrCreateBadgeByNormalNameAndColor(b); err != nil {
					return err
				} else if badge != nil {
					objectBadgesIDs = append(objectBadgesIDs, badge.ID)
				}
			}
		}
	}

	// save object actors;
	var objectActorsIDs []int64
	if prodObjectActors, err := ci.ProdDBM.GetObjectsActorRefs([]int64{prodObjectID}); err != nil {
		return err
	} else if len(prodObjectActors) != 0 {
		if objectActors, err := ci.ProdDBM.GetActorsByIDs(prodObjectActors.ActorsIDs()); err != nil {
			return err
		} else {
			for _, a := range objectActors {
				a.RootID = &rootID
				if actor, err := tx.GetOrCreateActorByNormalName(a); err != nil {
					return err
				} else if actor != nil {
					objectActorsIDs = append(objectActorsIDs, actor.ID)
				}
			}
		}
	}

	// save object actors;
	var objectMaterialsIDs []int64
	if prodObjectMaterials, err := ci.ProdDBM.GetMaterialRefsByObjectsIDs([]int64{prodObjectID}); err != nil {
		return err
	} else if len(prodObjectMaterials) != 0 {
		if objectMaterials, err := ci.ProdDBM.GetMaterialsByIDs(prodObjectMaterials.MaterialsIDs()); err != nil {
			return err
		} else {
			for _, m := range objectMaterials {
				m.RootID = &rootID
				if material, err := tx.GetOrCreateMaterialByNormalName(m); err != nil {
					return err
				} else if material != nil {
					objectMaterialsIDs = append(objectMaterialsIDs, material.ID)
				}
			}
		}
	}

	// save object status;
	if prodObjectStatuses, err := ci.ProdDBM.GetCurrentObjectsStatusesRefs([]int64{prodObjectID}); err != nil {
		return err
	} else if len(prodObjectStatuses) != 0 {
		for _, st := range prodObjectStatuses {
			st.ObjectID = object.ID

			if pid, err := ci.getObjectStatusRefID(st.ObjectStatusID, object.UserID); err == nil && pid != nil {
				st.ObjectStatusID = *pid
				if err := tx.CreateObjectStatusRef(st); err != nil {
					return err
				}
			}
		}
	}

	// seve refs
	if len(objectMediasIDs) != 0 {
		if err = tx.CreateObjectMediasRefs(object.ID, objectMediasIDs); err != nil {
			ci.Logger.Error("objects import", "save medias ref", err.Error())
			return err
		}
	}
	if len(objectBadgesIDs) != 0 {
		if err = tx.CreateObjectBadgesRefs(object.ID, objectBadgesIDs); err != nil {
			ci.Logger.Error("objects import", "save badges ref", err.Error())
			return err
		}
	}
	if len(objectActorsIDs) != 0 {
		if err = tx.CreateObjectActorsRefs(object.ID, objectActorsIDs); err != nil {
			ci.Logger.Error("objects import", "save actors ref", err.Error())
			return err
		}
	}
	if len(objectMaterialsIDs) != 0 {
		if err = tx.CreateObjectMaterialsRefs(object.ID, objectMaterialsIDs); err != nil {
			ci.Logger.Error("objects import", "save materials ref", err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		ci.Logger.Error("objects import", "commit transaction", err.Error())
		return err
	}

	return nil
}

func (ci *ImporterContext) saveMedia(tx dal.TxManager, imageMedia *dto.Media) (*dto.Media, error) {
	var fileName string
	var fileUrl string

	if imageMedia.Document != nil {
		fileName = imageMedia.Document.Name
		fileUrl = imageMedia.Document.URI
	} else if imageMedia.Photo != nil {
		fileName = imageMedia.Photo.Name
		var byteSize int64 = 0
		for _, u := range imageMedia.Photo.Variants {
			if u.ByteSize > byteSize {
				byteSize = u.ByteSize
				fileUrl = u.URI
			}
		}
	} else {
		return nil, errors.New("cannot get file name")
	}

	response, err := http.Get(fmt.Sprintf("https://prod.clr.su%v", fileUrl))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	mediaSaver := &resource.MediaSaver{
		DBM:           tx,
		Logger:        ci.Logger,
		Storage:       ci.FileStorage,
		CleaverClient: ci.CleaverClient,
	}

	inputFileData := &resource.InputFileData{
		OriginalFileName: fileName,
		Typo:             imageMedia.Type,
		UserID:           imageMedia.UserID,
		UniqID:           *imageMedia.UserUniqID,
		Content:          response.Body,
	}

	media, err := mediaSaver.SaveMedia(inputFileData)
	if err != nil {
		ci.Logger.Error("media saver", "err", err, "media_id", imageMedia.ID)
		return nil, err
	}

	ci.Logger.Info("import media", "new media", media.ID)
	return media, nil
}

func (ci *ImporterContext) runUpdateObjectValuations(email string) {
	user, err := getUserByEmail(ci.DBM, email)
	if err != nil {
		ci.Logger.Error("valuations import", "get user", err.Error())
		return
	}

	objectsList, err := ci.DBM.GetObjectsByUserIDsWithCustomFields([]int64{user.ID}, dto.ObjectAllFields)
	if err != nil {
		ci.Logger.Error("valuations import", "get objects", err.Error())
		return
	}

	valuations, err := ci.DBM.GetObjectValuations([]int64{551})
	if err != nil {
		ci.Logger.Error("valuations import", "get test valuation", err.Error())
		return
	}
	if len(valuations) == 0 {
		ci.Logger.Error("valuations import", "get test valuation", "not found")
	}
	valuation := valuations[0]

	for _, object := range objectsList {
		valuations, err := ci.DBM.GetObjectValuations([]int64{object.ID})
		if err != nil {
			ci.Logger.Error("valuations import", "get valuation", err.Error())
			return
		}
		if len(valuations) == 0 {
			valuation.ObjectID = object.ID
			if _, err := ci.DBM.CreateValuation(valuation); err != nil {
				ci.Logger.Error("valuations import", "save valuation", err.Error())
			} else {
				ci.Logger.Error("valuations import", "get valuation", valuation.ID)
			}
		}
	}
}

func (ci *ImporterContext) runUpdateObjectLocation(email string) {

	tx, err := ci.DBM.BeginTx()
	if err != nil {
		ci.Logger.Error("object origin_location import", "tx", err.Error())
		return
	}
	defer tx.Rollback()

	user, err := getUserByEmail(ci.DBM, email)
	if err != nil {
		ci.Logger.Error("object origin_location import", "get user", err.Error())
		return
	}

	var root *dto.Root
	if roots, err := ci.DBM.GetRootsByUserID(user.ID); err != nil {
		ci.Logger.Error("object origin_location import", "get root", err.Error())
		return
	} else if len(roots) != 1 {
		ci.Logger.Error("object origin_location import", "get root", "no root found")
		return
	} else {
		root = roots[0]
	}

	produser, err := getUserByEmail(ci.ProdDBM, email)
	if err != nil {
		ci.Logger.Error("collections import", "get prod user", err.Error())
		return
	}

	prodObjectsList, err := ci.ProdDBM.GetObjectsByUserIDsWithCustomFields([]int64{produser.ID}, dto.ObjectAllFields)
	if err != nil {
		ci.Logger.Error("objects import", "get objects", err.Error())
		return
	}

	for _, prodobject := range prodObjectsList {
		if objectID, _ := ci.DBM.GetObjectIDByUserUniqID(user.ID, prodobject.UserUniqID); objectID != nil {
			var objectLocationsIDs []int64
			if prodLocationBadges, err := ci.ProdDBM.GetObjectsOriginLocationRefs([]int64{prodobject.ID}); err != nil {
				ci.Logger.Error("object origin_location import", "get location", err.Error())
				return
			} else if len(prodLocationBadges) != 0 {
				if objectLocations, err := ci.ProdDBM.GetOriginLocationsByIDs(prodLocationBadges.OriginLocationsIDs()); err != nil {
					ci.Logger.Error("object origin_location import", "get location", err.Error())
					return
				} else {
					for _, l := range objectLocations {
						l.RootID = &root.ID
						if loc, err := tx.GetOrCreateOriginLocationByNormalName(l); err != nil {
							ci.Logger.Error("object origin_location import", "get or create location", err.Error())
							return
						} else if loc != nil {
							objectLocationsIDs = append(objectLocationsIDs, loc.ID)
						}
					}
				}
			}

			if len(objectLocationsIDs) != 0 {
				if err = tx.CreateObjectOriginLocationsRefs(*objectID, objectLocationsIDs); err != nil {
					ci.Logger.Error("object origin_location import", "save locations ref", err.Error())
					return
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		ci.Logger.Error("object origin_location import", "commit transaction", err.Error())
		return
	}

	ci.Logger.Info("object origin_location import", "status", "ok")
}

func NewInt64(v int64) *int64 {
	return &v
}

//func (ci *ImporterContext) runImportObjectsStatuses() {
//	statuses, _ := ci.ProdDBM.GetObjectStatuses()
//	for _, status := range statuses {
//		if status.ImageMediaID != nil {
//			if medias, err := ci.ProdDBM.GetMediasByIDs([]int64{*status.ImageMediaID}); err == nil && len(medias) == 1 {
//				media := medias[0]
//
//				media.UserID = nil
//				media.RootID = nil
//
//				nextID := util.NextUniqID()
//				media.UserUniqID = &nextID
//
//				if m, err := ci.saveMedia(media); err != nil {
//					ci.Logger.Error("object media import", "save media", err.Error())
//					return
//				} else if m != nil {
//					status.ImageMediaID = &m.ID
//				}
//			}
//
//			if err := ci.DBM.CreateObjectStatus(status); err != nil {
//				ci.Logger.Error("object statuses import", "save status", err.Error())
//				return
//			} else {
//				ci.Logger.Info("object status created", "ok", status.ID)
//			}
//		}
//	}
//}
