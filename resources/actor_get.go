package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/actors"
	"github.com/go-openapi/runtime/middleware"
)

// GetActorsForRoot TBD
func (a *Actor) GetActorsForRoot(params actors.GetActorGetParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get actors")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return actors.NewGetActorGetDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check access rights
	ok, err := NewAccessRightsChecker(a.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{params.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, params.RootID)
		logger.Error("user cant get entities here", "err", err)
		return actors.NewGetActorGetForbidden()
	}

	// actor list
	actorsList, err := a.Context.DBM.GetActorsByRootID(params.RootID)
	if err != nil {
		logger.Error("cant get actors list", "err", err)
		return errorResponse(500, err.Error())
	}

	return actors.NewGetActorGetOK().WithPayload(&models.AActors{
		Actors: models.NewModelActorList(actorsList),
	})
}
