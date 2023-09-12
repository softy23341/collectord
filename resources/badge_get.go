package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/badges"
	"github.com/go-openapi/runtime/middleware"
)

// GetBadgesForRoot TBD
func (b *Badge) GetBadgesForRoot(params badges.GetBadgeGetParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get badges for root")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return badges.NewGetBadgeGetDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check access rights
	ok, err := NewAccessRightsChecker(b.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{params.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, params.RootID)
		logger.Error("user cant get entities here", "err", err)
		return badges.NewGetBadgeGetForbidden()
	}

	badgesList, err := b.Context.DBM.GetBadgesByRootID(params.RootID)
	if err != nil {
		logger.Error("get badges list", "err", err)
		return errorResponse(500, err.Error())
	}

	return badges.NewGetBadgeGetOK().WithPayload(&models.ABadges{
		Badges: models.NewModelBadgeList(badgesList),
	})
}
