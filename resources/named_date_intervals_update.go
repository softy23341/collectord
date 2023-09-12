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
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// UpdateNamedDateInterval TBD
func (n *NamedDateInterval) UpdateNamedDateInterval(params nameddateintervals.PostNamedDateIntervalIDUpdateParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create named date interval")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return nameddateintervals.NewPostNamedDateIntervalIDUpdateDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return nameddateintervals.NewPostNamedDateIntervalIDUpdateConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return nameddateintervals.NewPostNamedDateIntervalIDUpdateNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	updateParams := params.RUpdateNamedDateInterval

	// get nameddateintervals
	var updatedNamedDateInterval *dto.NamedDateInterval
	if namedDateIntervalList, err := n.Context.DBM.GetNamedDateIntervalsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant get named date interval", "err", err)
		return errorResponse(500, err.Error())
	} else if len(namedDateIntervalList) != 1 {
		err := fmt.Errorf("cant find named date interval: %d", params.ID)
		logger.Error("cant find cant find named date interval", "err", err)
		return notFoundResponse(err.Error())
	} else {
		updatedNamedDateInterval = namedDateIntervalList[0]
	}

	// check basic rights
	if updatedNamedDateInterval.RootID == nil {
		err := fmt.Errorf("user cant update this actor")
		logger.Error("user cant update this actor", "err", err)
		return actors.NewPostActorIDUpdateForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(n.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*updatedNamedDateInterval.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *updatedNamedDateInterval.RootID)
		logger.Error("user cant allow be here", "err", err)
		return nameddateintervals.NewPostNamedDateIntervalIDUpdateForbidden()
	}

	if len(updateParams.Name) != 0 {
		newNamedDateIntervalNormalName := util.NormalizeString(updateParams.Name)
		if namedDateIntervalList, err := n.Context.DBM.GetNamedDayeIntervalsByNormalNames(*updatedNamedDateInterval.RootID, []string{newNamedDateIntervalNormalName}); err != nil {
			logger.Error("cant get named date interval", "err", err)
			return errorResponse(500, err.Error())
		} else if len(namedDateIntervalList) != 0 && namedDateIntervalList[0].ID != updatedNamedDateInterval.ID {
			err := fmt.Errorf("named date interval already present; normal name: %s", newNamedDateIntervalNormalName)
			logger.Error("named date interval already present", "err", err)
			return alreadyPresentResponse(err.Error())
		}

		updatedNamedDateInterval.Name = updateParams.Name
		updatedNamedDateInterval.NormalName = newNamedDateIntervalNormalName
	}

	if from := updateParams.From; from != nil {
		updatedNamedDateInterval.ProductionDateIntervalFrom = *from
	}

	if to := updateParams.To; to != nil {
		updatedNamedDateInterval.ProductionDateIntervalTo = *to
	}

	tx, err := n.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.UpdateNamedDateInterval(updatedNamedDateInterval); err != nil {
		logger.Error("cant update named date interval", "err", err)
		return errorResponse(500, err.Error())
	}

	editedNamedDateIntervalEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.NamedIntervalEntityType,
				ID:   updatedNamedDateInterval.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, editedNamedDateIntervalEvent); err != nil {
		logger.Error("edit named Date Interval event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	n.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		n.Context.EventSender.Send(editedNamedDateIntervalEvent)

		query := &services.ScrollSearchQuery{
			RootID: updatedNamedDateInterval.RootID,
			Filters: &dto.ObjectSearchFilters{
				ProductionDateIntervalID: &updatedNamedDateInterval.ID,
			},
		}
		err := n.Context.SearchClient.ScrollThrought(query, n.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant scroll th", "err", err)
		}
	}))

	return nameddateintervals.NewPostNamedDateIntervalIDUpdateNoContent()
}
