package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/originlocations"
	"github.com/go-openapi/runtime/middleware"
)

// GetOriginLocationsForRoot TBD
func (a *OriginLocation) GetOriginLocationsForRoot(params originlocations.GetOriginLocationGetParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get originLocations")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return originlocations.NewGetOriginLocationGetDefault(code).WithPayload(
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
		return originlocations.NewGetOriginLocationGetForbidden()
	}

	originLocationsList, err := a.Context.DBM.GetOriginLocationsByRootID(params.RootID)
	if err != nil {
		logger.Error("cant get originLocations list", "err", err)
		return errorResponse(500, err.Error())
	}

	return originlocations.NewGetOriginLocationGetOK().WithPayload(&models.AOriginLocations{
		OriginLocations: models.NewModelOriginLocationList(originLocationsList),
	})
}
