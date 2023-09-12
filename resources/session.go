package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/session"
	"github.com/go-openapi/runtime/middleware"
)

// Session TBD
type Session struct {
	Context Context
}

// SetDeviceToken TBD
func (s *Session) SetDeviceToken(params session.PostSessionRegisterDeviceTokenParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("set device token")

	p := params.RRegisterDeviceToken

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return session.NewPostSessionRegisterDeviceTokenDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	tx, err := s.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("cant start transaction", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// delete exist
	if err := tx.DeleteSessionDevices([]string{userContext.AToken()}); err != nil {
		logger.Error("cant delete sessions", "err", err)
		return errorResponse(500, err.Error())
	}

	userSessionDevice := &dto.UserSessionDevice{
		SessionID: userContext.Session.ID,
		Typo:      dto.UserSessionDeviceType(*p.Typo),
		Token:     *p.Token,
		Sandbox:   p.Sandbox,
	}

	// set new
	if err := tx.SetSessionDeviceByAuthToken(userSessionDevice); err != nil {
		logger.Error("cant set new device token", "err", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("cant commit transaction", "err", err)
		return errorResponse(500, err.Error())
	}

	return session.NewPostSessionRegisterDeviceTokenNoContent()
}
