package services

import (
	"errors"
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/dto"
	openapierrors "github.com/go-openapi/errors"
	"gopkg.in/inconshreveable/log15.v2"
)

// Autheticator TBD
type Autheticator interface {
	Auth(token string) (*auth.UserContext, error)
}

// DBAutheticator TBD
type DBAutheticator struct {
	DBM dal.TrManager
	Log log15.Logger
}

// Auth TBD
func (d *DBAutheticator) Auth(token string) (*auth.UserContext, error) {
	if len(token) == 0 {
		return nil, errors.New("empty auth token")
	}

	// TODO make one request with join and manual mapping
	var user *dto.User
	if users, err := d.DBM.GetUsersBySessionTokens([]string{token}); err != nil {
		d.Log.Error("cant get users by tokens", "err", err)
		return nil, err
	} else if len(users) != 1 {
		d.Log.Error("cant get users by tokens", "err", "cant find user")
		return nil, openapierrors.Unauthenticated("cant find user by token")
	} else {
		user = users[0]
	}

	var session *dto.UserSession
	if sessions, err := d.DBM.GetSessionsBySessionTokens([]string{token}); err != nil {
		d.Log.Error("cant get session by tokens", "err", err)
		return nil, fmt.Errorf("can't find session with auth token: %s", token)
	} else if len(sessions) != 1 {
		d.Log.Error("cant get session by tokens", "err", "cant find session")
		return nil, openapierrors.Unauthenticated("cant find session by token")
	} else {
		session = sessions[0]
	}

	return &auth.UserContext{
		User:    user,
		Session: session,
	}, nil
}
