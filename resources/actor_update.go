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
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// UpdateActor TBD
func (a *Actor) UpdateActor(params actors.PostActorIDUpdateParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create actor")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return actors.NewPostActorIDUpdateDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return actors.NewPostActorIDUpdateConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return actors.NewPostActorIDUpdateNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	updateParams := params.RUpdateActor

	// get actors
	var updatedActor *dto.Actor
	if actorList, err := a.Context.DBM.GetActorsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant get actors", "err", err)
		return errorResponse(500, err.Error())
	} else if len(actorList) != 1 {
		err := fmt.Errorf("cant find actor: %d", params.ID)
		logger.Error("cant find actor", "err", err)
		return notFoundResponse(err.Error())
	} else {
		updatedActor = actorList[0]
	}

	// check basic rights
	if updatedActor.RootID == nil {
		err := fmt.Errorf("user cant update this entity")
		logger.Error("user cant update this entity", "err", err)
		return actors.NewPostActorIDUpdateForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(a.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*updatedActor.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *updatedActor.RootID)
		logger.Error("user cant allow be here", "err", err)
		return actors.NewPostActorIDUpdateForbidden()
	}

	// by normal name
	newActorNormalName := util.NormalizeString(*updateParams.Name)
	if actorList, err := a.Context.DBM.GetActorsByNormalNames(*updatedActor.RootID, []string{newActorNormalName}); err != nil {
		logger.Error("cant get actors", "err", err)
		return errorResponse(500, err.Error())
	} else if len(actorList) != 0 && actorList[0].ID != updatedActor.ID {
		err := fmt.Errorf("badge already present; normal name: %s", newActorNormalName)
		logger.Error("actor already present", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	updatedActor.Name = *updateParams.Name
	updatedActor.NormalName = newActorNormalName

	tx, err := a.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.UpdateActor(updatedActor); err != nil {
		logger.Error("cant update actor", "err", err)
		return errorResponse(500, err.Error())
	}

	editedActorEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.ActorEntityType,
				ID:   updatedActor.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, editedActorEvent); err != nil {
		logger.Error("edited actor event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	a.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		a.Context.EventSender.Send(editedActorEvent)

		query := &services.ScrollSearchQuery{
			RootID: updatedActor.RootID,
			Filters: &dto.ObjectSearchFilters{
				Actors: []int64{updatedActor.ID},
			},
		}
		err := a.Context.SearchClient.ScrollThrought(query, a.Context.ReindexObjects)
		if err != nil {
			logger.Error("reindex objects with author", "err", err)
		}
	}))

	return actors.NewPostActorIDUpdateNoContent()
}
