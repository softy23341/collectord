package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/actors"
	"git.softndit.com/collector/backend/restapi/operations/originlocations"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// UpdateOriginLocation TBD
func (o *OriginLocation) UpdateOriginLocation(params originlocations.PostOriginLocationIDUpdateParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create originLocation")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return originlocations.NewPostOriginLocationIDUpdateDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return originlocations.NewPostOriginLocationIDUpdateConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return originlocations.NewPostOriginLocationIDUpdateNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	updateParams := params.RUpdateOriginLocation

	// get originLocations
	var updatedOriginLocation *dto.OriginLocation
	if originLocationList, err := o.Context.DBM.GetOriginLocationsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant get originLocations", "err", err)
		return errorResponse(500, err.Error())
	} else if len(originLocationList) != 1 {
		err := fmt.Errorf("cant find originLocation: %d", params.ID)
		logger.Error("cant find originLocation", "err", err)
		return notFoundResponse(err.Error())
	} else {
		updatedOriginLocation = originLocationList[0]
	}

	// check basic rights
	if updatedOriginLocation.RootID == nil {
		err := fmt.Errorf("user cant update this entity")
		logger.Error("user cant update this entity", "err", err)
		return actors.NewPostActorIDUpdateForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(o.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*updatedOriginLocation.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *updatedOriginLocation.RootID)
		logger.Error("user cant allow be here", "err", err)
		return actors.NewPostActorIDUpdateForbidden()
	}

	// by normal name
	newOriginLocationNormalName := util.NormalizeString(*updateParams.Name)
	if originLocationList, err := o.Context.DBM.GetOriginLocationsByNormalNames(*updatedOriginLocation.RootID, []string{newOriginLocationNormalName}); err != nil {
		logger.Error("cant get originLocations", "err", err)
		return errorResponse(500, err.Error())
	} else if len(originLocationList) != 0 && originLocationList[0].ID != updatedOriginLocation.ID {
		err := fmt.Errorf("badge already present; normal name: %s", newOriginLocationNormalName)
		logger.Error("originLocation already present", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	updatedOriginLocation.Name = *updateParams.Name
	updatedOriginLocation.NormalName = newOriginLocationNormalName

	tx, err := o.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.UpdateOriginLocation(updatedOriginLocation); err != nil {
		logger.Error("cant update originLocation", "err", err)
		return errorResponse(500, err.Error())
	}

	editedLocationEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.OriginLocationEntityType,
				ID:   updatedOriginLocation.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, editedLocationEvent); err != nil {
		logger.Error("edit location event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		o.Context.EventSender.Send(editedLocationEvent)
	}))

	return originlocations.NewPostOriginLocationIDUpdateNoContent()
}
