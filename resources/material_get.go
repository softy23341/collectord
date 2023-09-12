package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/materials"
	"github.com/go-openapi/runtime/middleware"
)

// GetMaterialsForRoot TBD
func (m *Material) GetMaterialsForRoot(params materials.GetMaterialGetParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get materials")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return materials.NewGetMaterialGetDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check access rights
	ok, err := NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{params.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, params.RootID)
		logger.Error("user cant get entities here", "err", err)
		return materials.NewGetMaterialGetForbidden()
	}

	materialsList, err := m.Context.DBM.GetMaterialsByRootID(params.RootID)
	if err != nil {
		logger.Error("cant get materials list", "err", err)
		return errorResponse(500, err.Error())
	}

	return materials.NewGetMaterialGetOK().WithPayload(&models.AMaterials{
		Materials: models.NewModelMaterialList(materialsList),
	})
}
