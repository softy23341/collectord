package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/events"
	"github.com/go-openapi/runtime/middleware"
)

// Event TBD
type Event struct {
	Context Context
}

// GetEvents TBD
func (e *Event) GetEvents(params events.PostEventParams, principal interface{}) middleware.Responder {
	eventContext := principal.(*auth.UserContext)
	logger := eventContext.Logger(params.HTTPRequest)
	logger.Debug("get events")

	errorResponse := func(code int) middleware.Responder {
		return events.NewPostEventDefault(code)
	}

	start, end := *params.RGetEvents.StartSeqNo, params.RGetEvents.EndSeqNo

	if start < 0 {
		start = 0
	}

	if end != nil && *end <= start {
		logger.Error("wrong input range")
		return errorResponse(422)
	}

	// dbm
	dbm := e.Context.DBM

	statuses := make([]dto.EventStatus, len(params.RGetEvents.Statuses))
	for i, s := range params.RGetEvents.Statuses {
		statuses[i] = dto.EventStatus(s)
	}

	eventsList, err := dbm.GetEvents(eventContext.User.ID, start, end, statuses)
	if err != nil {
		logger.Error("can't get evenets", "err", err)
		return errorResponse(500)
	}
	afterFilterEvents := make(dto.EventList, 0, len(eventsList))
	for _, event := range eventsList {
		if params.RGetEvents.Version == 0 && event.Type > 10 {
			continue
		}

		afterFilterEvents = append(afterFilterEvents, event)
	}

	return events.NewPostEventOK().WithPayload(&models.AGetEvents{
		Events: models.NewModelEventList(afterFilterEvents),
	})
}

// ConfirmEvents TBD
func (e *Event) ConfirmEvents(params events.PostEventConfirmParams, principal interface{}) middleware.Responder {
	eventContext := principal.(*auth.UserContext)
	logger := eventContext.Logger(params.HTTPRequest)
	logger.Debug("confirm events")

	errorResponse := func(code int) middleware.Responder {
		return events.NewPostEventConfirmDefault(code)
	}

	endSeqNo := *params.RConfirmEvent.EndSeqNo

	if endSeqNo < 1 {
		logger.Error("no sequences")
		return errorResponse(422)
	}

	err := e.Context.DBM.ConfirmEvents(eventContext.User.ID, endSeqNo)
	if err != nil {
		logger.Error("cant delete events", "err", err)
		return errorResponse(500)
	}

	return events.NewPostEventConfirmNoContent()
}
