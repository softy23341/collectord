package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/actors"
	"git.softndit.com/collector/backend/restapi/operations/nameddateintervals"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// DeleteNamedDateInterval TBD
func (n *NamedDateInterval) DeleteNamedDateInterval(params nameddateintervals.DeleteNamedDateIntervalIDDeleteParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("delete nameddateinterval")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return nameddateintervals.NewDeleteNamedDateIntervalIDDeleteDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return nameddateintervals.NewDeleteNamedDateIntervalIDDeleteNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	var deletedNamedDateInterval *dto.NamedDateInterval
	if namedDateIntervalsList, err := n.Context.DBM.GetNamedDateIntervalsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find named date interval err", "err", err)
		return errorResponse(500, err.Error())
	} else if len(namedDateIntervalsList) != 1 {
		err := fmt.Errorf("cant find nameddateinterval with id: %d", params.ID)
		logger.Error("cant find named date interval", "err", err)
		return notFoundResponse(err.Error())
	} else {
		deletedNamedDateInterval = namedDateIntervalsList[0]
	}

	// check basic rights
	if deletedNamedDateInterval.RootID == nil {
		err := fmt.Errorf("user cant delete this entity")
		logger.Error("user cant allow be here", "err", err)
		return actors.NewDeleteActorIDDeleteForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(n.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*deletedNamedDateInterval.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *deletedNamedDateInterval.RootID)
		logger.Error("user cant allow be here", "err", err)
		return nameddateintervals.NewDeleteNamedDateIntervalIDDeleteForbidden()
	}

	tx, err := n.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.DeleteNamedDateIntervalsByIDs([]int64{deletedNamedDateInterval.ID}); err != nil {
		logger.Error("cant delete named date intervals", "err", err)
		return errorResponse(500, err.Error())
	}

	deletedNamedDateIntervalEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedEntity: &dto.EntityRef{
				Typo: dto.NamedIntervalEntityType,
				ID:   deletedNamedDateInterval.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, deletedNamedDateIntervalEvent); err != nil {
		logger.Error("deleted named DateInterval event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	n.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		n.Context.EventSender.Send(deletedNamedDateIntervalEvent)

		query := &services.ScrollSearchQuery{
			RootID: deletedNamedDateInterval.RootID,
			Filters: &dto.ObjectSearchFilters{
				ProductionDateIntervalID: &deletedNamedDateInterval.ID,
			},
		}
		err := n.Context.SearchClient.ScrollThrought(query, n.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant scroll th", "err", err)
		}
	}))

	return nameddateintervals.NewDeleteNamedDateIntervalIDDeleteNoContent()
}
