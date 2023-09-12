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
	"github.com/go-openapi/runtime/middleware"
)

// DeleteMaterial TBD
func (m *Material) DeleteMaterial(params materials.DeleteMaterialIDDeleteParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("delete material")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return materials.NewDeleteMaterialIDDeleteDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return materials.NewDeleteMaterialIDDeleteNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	var deletedMaterial *dto.Material
	if materialsList, err := m.Context.DBM.GetMaterialsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find material err", "err", err)
		return errorResponse(500, err.Error())
	} else if len(materialsList) != 1 {
		err := fmt.Errorf("cant find material with id: %d", params.ID)
		logger.Error("cant find material", "err", err)
		return notFoundResponse(err.Error())
	} else {
		deletedMaterial = materialsList[0]
	}

	// check basic rights
	if deletedMaterial.RootID == nil {
		err := fmt.Errorf("user cant delete this entity")
		logger.Error("user cant allow be here", "err", err)
		return materials.NewDeleteMaterialIDDeleteForbidden()
	}

	// check access rights
	ok, err := NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{*deletedMaterial.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *deletedMaterial.RootID)
		logger.Error("user cant allow be here", "err", err)
		return materials.NewDeleteMaterialIDDeleteForbidden()
	}

	tx, err := m.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.DeleteMaterials([]int64{deletedMaterial.ID}); err != nil {
		logger.Error("cant delete materials", "err", err)
		return errorResponse(500, err.Error())
	}

	deletedMaterialEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedEntity: &dto.EntityRef{
				Typo: dto.MaterialEntityType,
				ID:   deletedMaterial.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, deletedMaterialEvent); err != nil {
		logger.Error("deleted material event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	m.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		m.Context.EventSender.Send(deletedMaterialEvent)

		query := &services.ScrollSearchQuery{
			RootID: deletedMaterial.RootID,
			Filters: &dto.ObjectSearchFilters{
				Materials: []int64{deletedMaterial.ID},
			},
		}
		err := m.Context.SearchClient.ScrollThrought(query, m.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant scroll th", "err", err)
		}
	}))

	return materials.NewDeleteMaterialIDDeleteNoContent()
}
