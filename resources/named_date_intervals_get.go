package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/nameddateintervals"
	"github.com/go-openapi/runtime/middleware"
)

// GetNamedDatesIntervalsForRoot TBD
func (n *NamedDateInterval) GetNamedDatesIntervalsForRoot(params nameddateintervals.GetNamedDateIntervalGetParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get named date intervals")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return nameddateintervals.NewGetNamedDateIntervalGetDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check access rights
	ok, err := NewAccessRightsChecker(n.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{params.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, params.RootID)
		logger.Error("user cant create named date intervals here", "err", err)
		return nameddateintervals.NewGetNamedDateIntervalGetForbidden()
	}

	namedDateIntervalList, err := n.Context.DBM.GetNamedDateIntervalsForRoots([]int64{params.RootID})
	if err != nil {
		logger.Error("cant get named date intervals list", "err", err)
		return errorResponse(500, err.Error())
	}

	return nameddateintervals.NewGetNamedDateIntervalGetOK().WithPayload(&models.ANamedDateIntervals{
		DateIntervals: models.NewModelDateIntervalList(namedDateIntervalList),
	})
}
