package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/originlocations"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// CreateOriginLocation TBD
func (o *OriginLocation) CreateOriginLocation(params originlocations.PostOriginLocationNewParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create originLocation")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return originlocations.NewPostOriginLocationNewDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return originlocations.NewPostOriginLocationNewConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	createParams := params.RCreateOriginLocation
	originLocationNormalName := util.NormalizeString(*createParams.Name)

	// check access rights
	ok, err := NewAccessRightsChecker(o.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*createParams.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *createParams.RootID)
		logger.Error("user cant create entity here", "err", err)
		return originlocations.NewPostOriginLocationNewForbidden()
	}

	if originLocationList, err := o.Context.DBM.GetOriginLocationsByNormalNames(*createParams.RootID, []string{originLocationNormalName}); err != nil {
		logger.Error("cant get originLocations by normal name", "err", err)
		return errorResponse(500, err.Error())
	} else if len(originLocationList) != 0 {
		err := fmt.Errorf("originLocation already present with name: %s", originLocationNormalName)
		logger.Error("cant create originLocation", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	newOriginLocation := &dto.OriginLocation{
		RootID:     createParams.RootID,
		Name:       *createParams.Name,
		NormalName: originLocationNormalName,
	}

	tx, err := o.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.CreateOriginLocation(newOriginLocation); err != nil {
		logger.Error("cant create originLocation", "err", err)
		return errorResponse(500, err.Error())
	}

	newLocationEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.OriginLocationEntityType,
				ID:   newOriginLocation.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newLocationEvent); err != nil {
		logger.Error("new location event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		o.Context.EventSender.Send(newLocationEvent)
	}))

	return originlocations.NewPostOriginLocationNewOK().WithPayload(&models.ACreateOriginLocation{
		ID: &newOriginLocation.ID,
	})
}
