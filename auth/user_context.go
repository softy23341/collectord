package auth

import (
	"net/http"

	"git.softndit.com/collector/backend/dto"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// UserContext TBD
type UserContext struct {
	User    *dto.User
	Session *dto.UserSession
}

// AToken TBD
func (u *UserContext) AToken() string {
	return u.Session.AuthToken
}

// Logger per-request logger
func (u *UserContext) Logger(HTTPRequest *http.Request) log15.Logger {
	return RequestLogger(HTTPRequest).
		New("token", u.AToken(), "user_id", u.User.ID)
}

// RequestLogger TBD
func RequestLogger(HTTPRequest *http.Request) log15.Logger {
	return HTTPRequest.Context().Value("logger").(log15.Logger)
}
