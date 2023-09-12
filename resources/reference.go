package resource

import (
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/references"
	"github.com/go-openapi/runtime/middleware"
)

// Reference TBD
type Reference struct {
	Context Context
}

// GetReferences TBD
func (o *Reference) GetReferences(params references.GetReferencesParams) middleware.Responder {
	logger := o.Context.Log

	logger.Debug("Get reference request")

	errorResponse := func(code int) middleware.Responder {
		c := int32(code)
		return references.NewGetReferencesDefault(code).WithPayload(
			&models.Error{Code: &c},
		)
	}

	// currencies
	currencies, err := o.Context.DBM.GetCurrencies()
	if err != nil {
		logger.Error("get currencies from DB", "error", err)
		return errorResponse(500)
	}

	// object statuses
	objectStatuses, err := o.Context.DBM.GetObjectStatuses()
	if err != nil {
		logger.Error("GetObjectStatuses", "error", err)
		return errorResponse(500)
	}

	mediaIDs := objectStatuses.GetImageMediaIDs()

	medias, err := o.Context.DBM.GetMediasByIDs(mediaIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "error", err)
		return errorResponse(500)
	}

	return references.NewGetReferencesOK().WithPayload(&models.AReferences{
		Medias:         models.NewModelMediaList(medias),
		Currencies:     models.NewModelCurrencyList(currencies),
		ObjectStatuses: models.NewModelObjectStatusList(objectStatuses),
	})
}
