package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/materials"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// CreateMaterial TBD
func (m *Material) CreateMaterial(params materials.PostMaterialNewParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create material")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return materials.NewPostMaterialNewDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return materials.NewPostMaterialNewConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	createParams := params.RCreateMaterial
	materialNormalName := util.NormalizeString(*createParams.Name)

	if materialList, err := m.Context.DBM.GetMaterialsByNormalNames(*createParams.RootID, []string{materialNormalName}); err != nil {
		logger.Error("cant get materials by normal name", "err", err)
		return errorResponse(500, err.Error())
	} else if len(materialList) != 0 {
		err := fmt.Errorf("material already present with name: %s", materialNormalName)
		logger.Error("cant create material", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	// check access rights
	ok, err := NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*createParams.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *createParams.RootID)
		logger.Error("user cant create entity here", "err", err)
		return materials.NewPostMaterialNewForbidden()
	}

	newMaterial := &dto.Material{
		RootID:     createParams.RootID,
		Name:       *createParams.Name,
		NormalName: materialNormalName,
	}

	tx, err := m.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.CreateMaterial(newMaterial); err != nil {
		logger.Error("cant create material", "err", err)
		return errorResponse(500, err.Error())
	}

	newMaterialEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.MaterialEntityType,
				ID:   newMaterial.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newMaterialEvent); err != nil {
		logger.Error("new material event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	m.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		m.Context.EventSender.Send(newMaterialEvent)
	}))

	return materials.NewPostMaterialNewOK().WithPayload(&models.ACreateMaterial{
		ID: &newMaterial.ID,
	})
}
