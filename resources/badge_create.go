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

// CreateBadge TBD
func (b *Badge) CreateBadge(params badges.PostBadgeNewParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("new badge")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return badges.NewPostBadgeNewDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return badges.NewPostBadgeNewConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	newBadgeParams := params.RCreateBadge
	inputBadge := models.NewDtoInputBadge(newBadgeParams.Badge)

	newBadge := &dto.Badge{
		RootID:     newBadgeParams.RootID,
		Color:      inputBadge.Color,
		Name:       inputBadge.Name,
		NormalName: util.NormalizeString(inputBadge.Name),
	}

	// check access rights
	ok, err := NewAccessRightsChecker(b.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*newBadgeParams.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *newBadgeParams.RootID)
		logger.Error("user cant create entity here", "err", err)
		return badges.NewPostBadgeNewForbidden()
	}

	// TODO: technically dangerous place; but there are uniq constraints in sql
	existBadges, err := b.Context.DBM.GetBadgesByNormalNamesOrColors(
		*newBadge.RootID,
		[]string{newBadge.NormalName},
		[]string{newBadge.Color},
	)

	if err != nil {
		logger.Error("cant get badges by normal names", "err", err)
		return errorResponse(500, err.Error())
	} else if len(existBadges) != 0 {
		err := fmt.Errorf("badge already present; normal name: %s", newBadge.NormalName)
		logger.Error("cant create badge", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	tx, err := b.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.CreateBadge(newBadge); err != nil {
		logger.Error("cant create badge", "err", err)
		return errorResponse(500, err.Error())
	}

	newBadgeEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.BadgeEntityType,
				ID:   newBadge.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newBadgeEvent); err != nil {
		logger.Error("new badge event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	b.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		b.Context.EventSender.Send(newBadgeEvent)
	}))

	return badges.NewPostBadgeNewOK().WithPayload(&models.ACreateBadge{
		ID: &newBadge.ID,
	})
}
