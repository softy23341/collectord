package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/badges"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// DeleteBadge TBD
func (b *Badge) DeleteBadge(params badges.DeleteBadgeIDDeleteParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("delete badge")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return badges.NewDeleteBadgeIDDeleteDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return badges.NewDeleteBadgeIDDeleteNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	var deletedBadge *dto.Badge
	if badgesList, err := b.Context.DBM.GetBadgesByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find badge err", "err", err)
		return errorResponse(500, err.Error())
	} else if len(badgesList) != 1 {
		err := fmt.Errorf("cant find badge with id: %d", params.ID)
		logger.Error("cant find badge", "err", err)
		return notFoundResponse(err.Error())
	} else {
		deletedBadge = badgesList[0]
	}

	// check basic rights
	if deletedBadge.RootID == nil {
		err := fmt.Errorf("user cant delete this entity")
		logger.Error("user cant allow be here", "err", err)
		return badges.NewDeleteBadgeIDDeleteForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(b.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*deletedBadge.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *deletedBadge.RootID)
		logger.Error("user cant allow be here", "err", err)
		return badges.NewDeleteBadgeIDDeleteForbidden()
	}

	tx, err := b.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.DeleteBadges([]int64{deletedBadge.ID}); err != nil {
		logger.Error("cant delete badges", "err", err)
		return errorResponse(500, err.Error())
	}

	deletedBadgeEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedEntity: &dto.EntityRef{
				Typo: dto.BadgeEntityType,
				ID:   deletedBadge.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, deletedBadgeEvent); err != nil {
		logger.Error("deleted badge event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	b.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		b.Context.EventSender.Send(deletedBadgeEvent)

		query := &services.ScrollSearchQuery{
			RootID: deletedBadge.RootID,
			Filters: &dto.ObjectSearchFilters{
				Badges: []int64{deletedBadge.ID},
			},
		}
		err := b.Context.SearchClient.ScrollThrought(query, b.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant scroll th", "err", err)
		}
	}))

	return badges.NewDeleteBadgeIDDeleteNoContent()
}
