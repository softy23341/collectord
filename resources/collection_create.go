package resource

import (
	"errors"
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/collections"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// CreateCollection TBD
func (c *Collection) CreateCollection(params collections.PostCollectionsParams, principal interface{}) middleware.Responder {
	var (
		userContext   = principal.(*auth.UserContext)
		logger        = userContext.Logger(params.HTTPRequest)
		DBM           = c.Context.DBM
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return collections.NewPostCollectionsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		successResponse = func(objectID int64) middleware.Responder {
			return collections.NewPostCollectionsOK().WithPayload(&models.ANewCollection{
				ID: &objectID,
			})
		}
		forbiddenResponse     = collections.NewPostCollectionsForbidden
		unprocessableResponse = collections.NewPostCollectionsUnprocessableEntity

		collectionParams = params.RNewCollection
		userID, rootID   = userContext.User.ID, *params.RNewCollection.RootID
		inputRights      = params.RNewCollection.Collection.Rights
	)

	logger.Debug("create collection")

	// check access rights
	ok, err := accessChecker.IsUserRelatedToRoots(userContext.User.ID, []int64{rootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, rootID)
		logger.Error("user cant create entity here", "err", err)
		return forbiddenResponse()
	}

	collectionID, err := DBM.GetCollectionIDByUserUniqID(userID, *collectionParams.ClientUniqID)
	if err != nil {
		logger.Error("GetCollectionIDByUserUniqID", "err", err)
		return errorResponse(500, err.Error())
	}
	if collectionID != nil {
		return successResponse(*collectionID)
	}

	collection, err := models.NewDtoCollection(collectionParams.Collection)
	if err != nil {
		logger.Error("NewDtoCollection", "err", err)
		return errorResponse(422, err.Error())
	}
	collection.RootID = rootID
	collection.UserID = &userID
	collection.UserUniqID = collectionParams.ClientUniqID

	// check media
	if collection.ImageMediaID != nil {
		medias, err := DBM.GetMediasByIDs([]int64{*collection.ImageMediaID})
		if err != nil || len(medias) != 1 {
			if err == nil {
				err = errors.New("can't find image")
			}
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500, err.Error())
		}
		accessChecker.IsUserOwnerOfMedias(userContext.User.ID, medias)
		if !ok {
			err := fmt.Errorf("user is not owner of medias: %d => %v",
				userContext.User.ID, medias.GetIDs())
			logger.Error("user cant use these medias", "err", err)
			return forbiddenResponse()
		}

	}

	// check groups
	collectionGroupsIDs := collectionParams.Collection.GroupsIds
	collectionGroupsList, err := DBM.GetGroupsByIDs(collectionGroupsIDs)
	if err != nil {
		logger.Error("GetGroupsByIDs", "error", err)
		return errorResponse(500, err.Error())
	} else if len(collectionGroupsList) != len(collectionGroupsIDs) {
		err := errors.New("cant find all groups")
		logger.Error("cant find all groups", "error", err)
		return errorResponse(500, err.Error())
	} else if !util.Int64Slice(collectionGroupsList.RootsIDs()).All(rootID) {
		err := errors.New("cant mix roots entities")
		logger.Error("cant attach collection to groups", "err", err)
		return forbiddenResponse()
	}

	// check rights to groups
	ok, err = accessChecker.
		HasUserRightsForGroups(
			userID,
			dto.RightEntityLevelWrite,
			collectionGroupsIDs,
		)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user cant touch groups: %d => %v", userID, collectionGroupsIDs)
		logger.Error("user cant create collections here", "err", err)
		return forbiddenResponse()
	}

	rootUsers, err := DBM.GetUsersByRootID(rootID)
	if err != nil {
		logger.Error("cant get users by root", "err", err)
		return errorResponse(500, err.Error())
	}

	rootOwner, err := DBM.GetRootOwner(rootID)
	if err != nil {
		logger.Error("cant get root owner", "err", err)
		return errorResponse(500, err.Error())
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// save collection
	if err := tx.CreateCollection(collection); err != nil {
		logger.Error("CreateCollection", "err", err)
		return errorResponse(500, err.Error())
	}

	// prepare rights
	rights := make(dto.UserEntityRightList, 0, len(inputRights))
	id2user := rootUsers.IDToUser()

	for _, inputRight := range inputRights {
		if _, found := id2user[*inputRight.UserID]; !found {
			logger.Error("user is not in the root", "userID", inputRight.UserID)
			return forbiddenResponse()
		}

		if !dto.IsRightEntityLevelValid(*inputRight.Level) {
			logger.Error("invalid entity type", "type", *inputRight.Level)
			return unprocessableResponse()
		}

		if *inputRight.UserID == rootOwner.ID || *inputRight.UserID == userID {
			continue
		}

		inputLevel := dto.RightEntityLevel(*inputRight.Level)
		rights = append(rights, &dto.UserEntityRight{
			UserID:     *inputRight.UserID,
			EntityType: dto.RightEntityTypeCollection,
			EntityID:   collection.ID,
			Level:      inputLevel,
			RootID:     rootID,
		})
	}

	rights = append(rights, &dto.UserEntityRight{
		UserID:     rootOwner.ID,
		EntityType: dto.RightEntityTypeCollection,
		EntityID:   collection.ID,
		Level:      dto.RightEntityLevelAdmin,
		RootID:     rootID,
	})
	if userID != rootOwner.ID {
		rights = append(rights, &dto.UserEntityRight{
			UserID:     userID,
			EntityType: dto.RightEntityTypeCollection,
			EntityID:   collection.ID,
			Level:      dto.RightEntityLevelAdmin,
			RootID:     rootID,
		})
	}

	// set rights
	for _, right := range rights {
		if err := tx.PutUserRight(right); err != nil {
			logger.Error("cant set Right", "err", err)
			return errorResponse(500, err.Error())
		}
	}

	// create refs
	if err := tx.CreateCollectionGroupsRefs(collection.ID, collectionGroupsIDs); err != nil {
		logger.Error("CreateCollectionGroupsRefs", "err", err)
		return errorResponse(500, err.Error())
	}

	newCollectionEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.CollectionEntityType,
				ID:   collection.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newCollectionEvent); err != nil {
		logger.Error("new collection event", "error", err)
		return errorResponse(500, err.Error())
	}

	// commit
	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	c.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		c.Context.EventSender.Send(newCollectionEvent)
	}))

	return successResponse(collection.ID)
}
