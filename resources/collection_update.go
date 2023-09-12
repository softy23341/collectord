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

// UpdateCollection TBD
func (c *Collection) UpdateCollection(params collections.PutCollectionsIDParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)

		DBM           = c.Context.DBM
		logger        = userContext.Logger(params.HTTPRequest)
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return collections.NewPostCollectionsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		successResponse = func() middleware.Responder {
			return collections.NewPutCollectionsIDNoContent()
		}
		forbiddenResponse     = collections.NewPutCollectionsIDForbidden
		unprocessableResponse = collections.NewPutCollectionsIDUnprocessableEntity

		inputRights = params.REditCollection.Collection.Rights
	)

	logger.Debug("UpdateCollection")

	var collection *dto.Collection
	if collections, err := DBM.GetCollectionsByIDs([]int64{params.ID}); err != nil {
		logger.Error("GetCollectionsByIDs", "error", err)
		return errorResponse(500, err.Error())
	} else if len(collections) != 1 {
		err = errors.New("collection not found")
		logger.Error("GetCollectionsByIDs", "error", err)
		return errorResponse(404, err.Error())
	} else {
		collection = collections[0]
	}

	// check collection
	ok, err := accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{collection.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant update collection %d", collection.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	collectionParams := params.REditCollection.Collection

	// name
	if name := collectionParams.Name; name != "" {
		collection.Name = name
	}

	// description
	if description := collectionParams.Description; description != nil {
		collection.Description = *description
	}

	// media
	if imageID := params.REditCollection.Collection.ImageMediaID; imageID != nil {
		collection.ImageMediaID = imageID.Value
		if imageID.Value != nil {
			medias, err := DBM.GetMediasByIDs([]int64{*imageID.Value})
			if err != nil || len(medias) != 1 {
				if err == nil {
					err = errors.New("can't find image")
				}
				logger.Error("GetMediasByIDs", "error", err)
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
	}

	// public
	if isPublic := collectionParams.Public; isPublic != nil {
		collection.Public = *isPublic
	}

	// isAnonymous
	if isAnonymous := collectionParams.IsAnonymous; isAnonymous != nil {
		collection.IsAnonymous = *isAnonymous
	}

	// check groups
	collectionGroupsIDs := collectionParams.GroupsIds
	collectionGroupsList, err := DBM.GetGroupsByIDs(collectionGroupsIDs)
	if err != nil {
		logger.Error("GetGroupsByIDs", "error", err)
		return errorResponse(500, err.Error())
	} else if len(collectionGroupsList) != len(collectionGroupsIDs) {
		err := errors.New("cant find all groups")
		logger.Error("cant find all groups", "error", err)
		return errorResponse(500, err.Error())
	} else if !util.Int64Slice(collectionGroupsList.RootsIDs()).All(collection.RootID) {
		err := errors.New("cant mix roots entities")
		logger.Error("cant attach collection to groups", "err", err)
		return forbiddenResponse()
	}

	// check rights to groups
	ok, err = accessChecker.
		HasUserRightsForGroups(
			userContext.User.ID,
			dto.RightEntityLevelWrite,
			collectionGroupsIDs,
		)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user cant touch groups: %d => %v", userContext.User.ID, collectionGroupsIDs)
		logger.Error("user cant create collections here", "err", err)
		return forbiddenResponse()
	}

	rootUsers, err := DBM.GetUsersByRootID(collection.RootID)
	if err != nil {
		logger.Error("cant get users by root", "err", err)
		return errorResponse(500, err.Error())
	}

	rootOwner, err := DBM.GetRootOwner(collection.RootID)
	if err != nil {
		logger.Error("cant get root owner", "err", err)
		return errorResponse(500, err.Error())
	}

	if len(inputRights) > 0 {
		ok, err := accessChecker.HasUserRightsForCollections(
			userContext.User.ID,
			dto.RightEntityLevelAdmin,
			[]int64{collection.ID},
		)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return errorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("cant update collection rights %d", collection.ID)
			logger.Error("access denied", "err", err.Error())
			return forbiddenResponse()
		}
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

		if *inputRight.UserID == rootOwner.ID || *inputRight.UserID == userContext.User.ID {
			continue
		}

		inputLevel := dto.RightEntityLevel(*inputRight.Level)
		rights = append(rights, &dto.UserEntityRight{
			UserID:     *inputRight.UserID,
			EntityType: dto.RightEntityTypeCollection,
			EntityID:   collection.ID,
			Level:      inputLevel,
			RootID:     collection.RootID,
		})
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// set rights
	for _, right := range rights {
		if err := tx.PutUserRight(right); err != nil {
			logger.Error("cant set Right", "err", err)
			return errorResponse(500, err.Error())
		}
	}

	if err := tx.DeleteCollectionGroupsRefs(collection.ID); err != nil {
		logger.Error("cant delete collection groups refs", "err", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.CreateCollectionGroupsRefs(collection.ID, collectionGroupsIDs); err != nil {
		logger.Error("cant create collection groups refs", "err", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.UpdateCollection(collection); err != nil {
		logger.Error("update collection", "error", err)
		return errorResponse(500, err.Error())
	}

	editedCollectionEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.CollectionEntityType,
				ID:   collection.ID,
			},
		},
	}
	if _, err := services.EmplaceEvent(tx, editedCollectionEvent); err != nil {
		logger.Error("delete collection event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	c.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		c.Context.EventSender.Send(editedCollectionEvent)
	}))

	return successResponse()
}
