package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/actors"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// DeleteActor TBD
func (a *Actor) DeleteActor(params actors.DeleteActorIDDeleteParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("delete actor")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return actors.NewDeleteActorIDDeleteDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return actors.NewDeleteActorIDDeleteNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	// get actor to delete
	var deletedActor *dto.Actor
	if actorsList, err := a.Context.DBM.GetActorsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find actor err", "err", err)
		return errorResponse(500, err.Error())
	} else if len(actorsList) != 1 {
		err := fmt.Errorf("cant find actor with id: %d", params.ID)
		logger.Error("cant find actor", "err", err)
		return notFoundResponse(err.Error())
	} else {
		deletedActor = actorsList[0]
	}

	// check basic rights
	if deletedActor.RootID == nil {
		err := fmt.Errorf("user cant delete this entity")
		logger.Error("user cant allow be here", "err", err)
		return actors.NewDeleteActorIDDeleteForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(a.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*deletedActor.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *deletedActor.RootID)
		logger.Error("user cant allow be here", "err", err)
		return actors.NewDeleteActorIDDeleteForbidden()
	}

	tx, err := a.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.DeleteActors([]int64{deletedActor.ID}); err != nil {
		logger.Error("cant delete actors", "err", err)
		return errorResponse(500, err.Error())
	}

	deletedActorEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedEntity: &dto.EntityRef{
				Typo: dto.ActorEntityType,
				ID:   deletedActor.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, deletedActorEvent); err != nil {
		logger.Error("deleted actor event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	a.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		a.Context.EventSender.Send(deletedActorEvent)

		query := &services.ScrollSearchQuery{
			RootID: deletedActor.RootID,
			Filters: &dto.ObjectSearchFilters{
				Actors: []int64{deletedActor.ID},
			},
		}

		err := a.Context.SearchClient.ScrollThrought(query, a.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant scroll th", "err", err)
		}
	}))

	return actors.NewDeleteActorIDDeleteNoContent()
}
