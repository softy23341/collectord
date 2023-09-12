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

// UpdateMaterial TBD
func (m *Material) UpdateMaterial(params materials.PostMaterialIDUpdateParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create material")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return materials.NewPostMaterialIDUpdateDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	alreadyPresentResponse := func(msg string) middleware.Responder {
		code := int32(409)
		return materials.NewPostMaterialIDUpdateConflict().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return materials.NewPostMaterialIDUpdateNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	updateParams := params.RUpdateMaterial

	// get materials
	var updatedMaterial *dto.Material
	if materialList, err := m.Context.DBM.GetMaterialsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant get materials", "err", err)
		return errorResponse(500, err.Error())
	} else if len(materialList) != 1 {
		err := fmt.Errorf("cant find material: %d", params.ID)
		logger.Error("cant find material", "err", err)
		return notFoundResponse(err.Error())
	} else {
		updatedMaterial = materialList[0]
	}

	// check basic rights
	if updatedMaterial.RootID == nil {
		err := fmt.Errorf("user cant update this entity")
		logger.Error("user cant update this entity", "err", err)
		return materials.NewPostMaterialIDUpdateForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*updatedMaterial.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *updatedMaterial.RootID)
		logger.Error("user cant allow be here", "err", err)
		return materials.NewPostMaterialIDUpdateForbidden()
	}

	// by normal name
	newMaterialNormalName := util.NormalizeString(*updateParams.Name)
	if materialList, err := m.Context.DBM.GetMaterialsByNormalNames(*updatedMaterial.RootID, []string{newMaterialNormalName}); err != nil {
		logger.Error("cant get materials", "err", err)
		return errorResponse(500, err.Error())
	} else if len(materialList) != 0 && materialList[0].ID != updatedMaterial.ID {
		err := fmt.Errorf("badge already present; normal name: %s", newMaterialNormalName)
		logger.Error("material already present", "err", err)
		return alreadyPresentResponse(err.Error())
	}

	updatedMaterial.Name = *updateParams.Name
	updatedMaterial.NormalName = newMaterialNormalName

	tx, err := m.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.UpdateMaterial(updatedMaterial); err != nil {
		logger.Error("cant update material", "err", err)
		return errorResponse(500, err.Error())
	}

	editedMaterialEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.MaterialEntityType,
				ID:   updatedMaterial.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, editedMaterialEvent); err != nil {
		logger.Error("edit material event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	m.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		m.Context.EventSender.Send(editedMaterialEvent)
	}))

	return materials.NewPostMaterialIDUpdateNoContent()
}
