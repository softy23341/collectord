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
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// UpdateBadge TBD
func (b *Badge) UpdateBadge(params badges.PostBadgeIDUpdateParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("update badge")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return badges.NewPostBadgeIDUpdateDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return badges.NewPostBadgeIDUpdateConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return badges.NewPostBadgeIDUpdateNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	updateParams := params.RUpdateBadge
	inputUpdateBadge := models.NewDtoInputBadge(updateParams.Badge)

	var updatedBadge *dto.Badge
	if badgesList, err := b.Context.DBM.GetBadgesByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find badge err", "err", err)
		return errorResponse(500, err.Error())
	} else if len(badgesList) != 1 {
		err := fmt.Errorf("cant find badge with id: %d", params.ID)
		logger.Error("cant find badge", "err", err)
		return notFoundResponse(err.Error())
	} else {
		updatedBadge = badgesList[0]
	}

	// check basic rights
	if updatedBadge.RootID == nil {
		err := fmt.Errorf("user cant update this entity")
		logger.Error("user cant update this entity", "err", err)
		return badges.NewPostBadgeIDUpdateForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(b.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*updatedBadge.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *updatedBadge.RootID)
		logger.Error("user cant allow be here", "err", err)
		return badges.NewPostBadgeIDUpdateForbidden()
	}

	// TODO: technically dangerous place; but there are uniq constraints in sql
	existBadges, err := b.Context.DBM.GetBadgesByNormalNamesOrColors(
		*updatedBadge.RootID,
		[]string{util.NormalizeString(inputUpdateBadge.Name)},
		[]string{inputUpdateBadge.Color},
	)
	if err != nil {
		logger.Error("cant get badges by normal names", "err", err)
		return errorResponse(500, err.Error())
	} else if len(existBadges) != 0 && existBadges[0].ID != params.ID {
		err := fmt.Errorf("badge already present; normal name: %s", inputUpdateBadge.Name)
		logger.Error("cant create badge", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	updatedBadge.Name = inputUpdateBadge.Name
	updatedBadge.NormalName = util.NormalizeString(inputUpdateBadge.Name)
	updatedBadge.Color = inputUpdateBadge.Color

	tx, err := b.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.UpdateBadge(updatedBadge); err != nil {
		logger.Error("cant update badge", "err", err)
		return errorResponse(500, err.Error())
	}

	editedBadgeEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.BadgeEntityType,
				ID:   updatedBadge.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, editedBadgeEvent); err != nil {
		logger.Error("save edited badge event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	b.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		b.Context.EventSender.Send(editedBadgeEvent)
	}))

	return badges.NewPostBadgeIDUpdateNoContent()
}
