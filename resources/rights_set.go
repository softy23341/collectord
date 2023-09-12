package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/rights"
	"github.com/go-openapi/runtime/middleware"
)

// SetRight TBD
func (uer *UserEntityRight) SetRight(params rights.PutRightParams, principal interface{}) middleware.Responder {
	var (
		userContext   = principal.(*auth.UserContext)
		logger        = userContext.Logger(params.HTTPRequest)
		DBM           = uer.Context.DBM
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		// responses
		defaultErrorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return rights.NewPutRightDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		unprocessableEntityResponse = func(msg string) middleware.Responder {
			c := int32(422)
			return rights.NewPutRightUnprocessableEntity().WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		forbiddenResponse = rights.NewPutRightForbidden

		inputRight   = params.RSetRight.Right
		entityRootID int64
	)

	logger.Debug("set rights")

	// validate
	if !dto.IsRightEntityTypeValid(inputRight.EntityType) {
		logger.Error("invalid entity type", "type", inputRight.EntityType)
		return unprocessableEntityResponse("invalid entity type")
	}
	if !dto.IsRightEntityLevelValid(inputRight.Level) {
		logger.Error("invalid right level", "level", inputRight.Level)
		return unprocessableEntityResponse("invalid right level")
	}

	if inputRight.UserID == userContext.User.ID {
		logger.Error("you cant set own rights")
		return forbiddenResponse()
	}

	// get target user
	var targetUser *dto.User
	if usersList, err := DBM.GetUsersByIDs([]int64{inputRight.UserID}); err != nil {
		logger.Error("cant get user by id", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(usersList) != 1 {
		logger.Error("cant find user by id", "err", err)
		return unprocessableEntityResponse("invalid user if")
	} else {
		targetUser = usersList[0]
	}

	// get user root
	var targetUserRoot *dto.Root
	if rootList, err := DBM.GetMainUserRoot(targetUser.ID); err != nil {
		logger.Error("cant GetMainUserRoot", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(rootList) != 1 {
		err := fmt.Errorf("cant find user main root for: %d", userContext.User.ID)
		logger.Error("cant find root", "err", err.Error())
		return defaultErrorResponse(500, err.Error())
	} else {
		targetUserRoot = rootList[0]
	}

	// check access
	switch dto.RightEntityType(inputRight.EntityType) {
	case dto.RightEntityTypeCollection:
		if collectionList, err := DBM.GetCollectionsByIDs([]int64{inputRight.EntityID}); err != nil {
			logger.Error("cant get collections", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if len(collectionList) != 1 {
			err := fmt.Errorf("cant find entity %d", inputRight.EntityID)
			logger.Error("cant find entity", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else {
			targetCollection := collectionList[0]
			entityRootID = targetCollection.RootID

			// check collection
			ok, err := accessChecker.HasUserRightsForCollections(
				userContext.User.ID,
				dto.RightEntityLevelAdmin,
				[]int64{targetCollection.ID},
			)
			if err != nil {
				logger.Error("cant check user rights for collection", "err", err)
				return defaultErrorResponse(500, err.Error())
			} else if !ok {
				logger.Error("user cant set rights on this collection 1",
					"collectionID", targetCollection.ID)
				return forbiddenResponse()
			}
		}
	default:
		logger.Error("invalid right entity type", "type", inputRight.EntityType)
		return unprocessableEntityResponse("invalid entity type")
	}

	if entityRootID == targetUserRoot.ID {
		logger.Error("user cant touch owner rights", "entityRootID", entityRootID)
		return forbiddenResponse()
	}

	// construct right
	right := &dto.UserEntityRight{
		UserID:     inputRight.UserID,
		EntityType: dto.RightEntityType(inputRight.EntityType),
		EntityID:   inputRight.EntityID,
		Level:      dto.RightEntityLevel(inputRight.Level),
		RootID:     entityRootID,
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return defaultErrorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// put right
	if err := tx.PutUserRight(right); err != nil {
		logger.Error("cant set Right", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// end commit
	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	return rights.NewPutRightOK().WithPayload(&models.ASetRight{
		Right: models.NewModelShortUserEntityRight(right.ToShort()),
	})
}
