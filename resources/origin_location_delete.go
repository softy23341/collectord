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
	"github.com/go-openapi/runtime/middleware"
)

// DeleteOriginLocation TBD
func (o *OriginLocation) DeleteOriginLocation(params originlocations.DeleteOriginLocationIDDeleteParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("delete originLocation")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return originlocations.NewDeleteOriginLocationIDDeleteDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return originlocations.NewDeleteOriginLocationIDDeleteNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	var deletedOriginLocation *dto.OriginLocation
	if originLocationsList, err := o.Context.DBM.GetOriginLocationsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find originLocation err", "err", err)
		return errorResponse(500, err.Error())
	} else if len(originLocationsList) != 1 {
		err := fmt.Errorf("cant find originLocation with id: %d", params.ID)
		logger.Error("cant find originLocation", "err", err)
		return notFoundResponse(err.Error())
	} else {
		deletedOriginLocation = originLocationsList[0]
	}

	// check basic rights
	if deletedOriginLocation.RootID == nil {
		err := fmt.Errorf("user cant delete this actor")
		logger.Error("user cant allow be here", "err", err)
		return originlocations.NewDeleteOriginLocationIDDeleteForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(o.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*deletedOriginLocation.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *deletedOriginLocation.RootID)
		logger.Error("user cant allow be here", "err", err)
		return originlocations.NewDeleteOriginLocationIDDeleteForbidden()
	}

	tx, err := o.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := o.Context.DBM.DeleteOriginLocations([]int64{deletedOriginLocation.ID}); err != nil {
		logger.Error("cant delete originLocations", "err", err)
		return errorResponse(500, err.Error())
	}

	deletedLocationEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedEntity: &dto.EntityRef{
				Typo: dto.OriginLocationEntityType,
				ID:   deletedOriginLocation.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, deletedLocationEvent); err != nil {
		logger.Error("deleted location event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		o.Context.EventSender.Send(deletedLocationEvent)

		query := &services.ScrollSearchQuery{
			RootID: deletedOriginLocation.RootID,
			Filters: &dto.ObjectSearchFilters{
				OriginLocations: []int64{deletedOriginLocation.ID},
			},
		}
		err := o.Context.SearchClient.ScrollThrought(query, o.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant scroll th", "err", err)
		}
	}))

	return originlocations.NewDeleteOriginLocationIDDeleteNoContent()
}
