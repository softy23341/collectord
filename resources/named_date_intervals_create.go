package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/nameddateintervals"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// CreateNamedDateInterval TBD
func (n *NamedDateInterval) CreateNamedDateInterval(params nameddateintervals.PostNamedDateIntervalNewParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create nameddateinterval")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return nameddateintervals.NewPostNamedDateIntervalNewDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return nameddateintervals.NewPostNamedDateIntervalNewConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	createParams := params.RCreateNamedDateInterval
	namedDateIntervalNormalName := util.NormalizeString(*createParams.Name)

	if namedDateIntervalList, err := n.Context.DBM.GetNamedDayeIntervalsByNormalNames(*createParams.RootID, []string{namedDateIntervalNormalName}); err != nil {
		logger.Error("cant get named date interval by normal name", "err", err)
		return errorResponse(500, err.Error())
	} else if len(namedDateIntervalList) != 0 {
		err := fmt.Errorf("named date interval already present with name: %s", namedDateIntervalNormalName)
		logger.Error("cant create nameddateinterval", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	// check access rights
	ok, err := NewAccessRightsChecker(n.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*createParams.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *createParams.RootID)
		logger.Error("user cant create entity here", "err", err)
		return nameddateintervals.NewPostNamedDateIntervalNewForbidden()
	}

	newNamedDateInterval := &dto.NamedDateInterval{
		RootID:                     createParams.RootID,
		Name:                       *createParams.Name,
		NormalName:                 namedDateIntervalNormalName,
		ProductionDateIntervalFrom: *createParams.From,
		ProductionDateIntervalTo:   *createParams.To,
	}

	tx, err := n.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.CreateNamedDateInterval(newNamedDateInterval); err != nil {
		logger.Error("cant create named date interval", "err", err)
		return errorResponse(500, err.Error())
	}

	newNamedDateIntervalEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.NamedIntervalEntityType,
				ID:   newNamedDateInterval.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newNamedDateIntervalEvent); err != nil {
		logger.Error("new named date interval event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	n.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		n.Context.EventSender.Send(newNamedDateIntervalEvent)
	}))

	return nameddateintervals.NewPostNamedDateIntervalNewOK().WithPayload(&models.ACreateNamedDateInterval{
		ID: &newNamedDateInterval.ID,
	})
}
