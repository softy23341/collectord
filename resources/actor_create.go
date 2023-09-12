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

// CreateActor TBD
func (a *Actor) CreateActor(params actors.PostActorNewParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create actor")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return actors.NewPostActorNewDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return actors.NewPostActorNewConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	createParams := params.RCreateActor
	actorNormalName := util.NormalizeString(*createParams.Name)

	// check access rights
	ok, err := NewAccessRightsChecker(a.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*createParams.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *createParams.RootID)
		logger.Error("user cant create entity here", "err", err)
		return actors.NewPostActorNewForbidden()
	}

	// check presents
	if actorList, err := a.Context.DBM.GetActorsByNormalNames(*createParams.RootID, []string{actorNormalName}); err != nil {
		logger.Error("cant get actors by normal name", "err", err)
		return errorResponse(500, err.Error())
	} else if len(actorList) != 0 {
		err := fmt.Errorf("actor already present with name: %s", actorNormalName)
		logger.Error("cant create actor", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	newActor := &dto.Actor{
		RootID:     createParams.RootID,
		Name:       *createParams.Name,
		NormalName: actorNormalName,
	}

	tx, err := a.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.CreateActor(newActor); err != nil {
		logger.Error("cant create actor", "err", err)
		return errorResponse(500, err.Error())
	}

	newActorEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.ActorEntityType,
				ID:   newActor.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newActorEvent); err != nil {
		logger.Error("new actor event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	a.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		a.Context.EventSender.Send(newActorEvent)
	}))

	return actors.NewPostActorNewOK().WithPayload(&models.ACreateActor{
		ID: &newActor.ID,
	})
}
