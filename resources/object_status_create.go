package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/objectstatus"
	"github.com/go-openapi/runtime/middleware"
)

// CreateObjectStatus TBD
func (o *ObjectStatus) CreateObjectStatus(params objectstatus.PostObjectStatusNewParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create object status")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return objectstatus.NewPostObjectStatusNewDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	return errorResponse(500, "deprecated")

	createParams := params.RCreateObjectStatus

	newObjectStatus := &dto.ObjectStatus{
		Name:         *createParams.Name,
		Description:  createParams.Description,
		ImageMediaID: createParams.MediaID,
	}

	if err := o.Context.DBM.CreateObjectStatus(newObjectStatus); err != nil {
		logger.Error("cant create object status", "err", err)
		return errorResponse(500, err.Error())
	}

	return objectstatus.NewPostObjectStatusNewOK().WithPayload(&models.ACreateObjectStatus{
		ID: &newObjectStatus.ID,
	})
}
