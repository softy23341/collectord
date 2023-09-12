package resource

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	sauth "git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/auth/storage"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/auth"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// AuthCookieKey TBD
const (
	AuthCookieKey = "auth-token"
	AuthHeaderKey = "auth-token"
)

// Auth TBD
type (
	Auth struct {
		Context           Context
		EmailTokenStorage storage.Storage
		Templates         *template.Template
	}

	cookieAuthRespond struct {
		responder middleware.Responder
		authToken string
	}
)

func newCookieAuthRespond(responder middleware.Responder, authToken string) *cookieAuthRespond {
	return &cookieAuthRespond{
		responder: responder,
		authToken: authToken,
	}
}

// WriteResponse TBD
func (c *cookieAuthRespond) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	cookie := http.Cookie{
		Name:     AuthCookieKey,
		Value:    c.authToken,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(10 * 365 * 24 * time.Hour),
	}
	http.SetCookie(rw, &cookie)
	c.responder.WriteResponse(rw, p)
}

// Login TBD
func (a *Auth) Login(params auth.PostAuthLoginParams) middleware.Responder {
	logger := sauth.RequestLogger(params.HTTPRequest)
	logger.Debug("Login")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthLoginDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthLoginForbidden().WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check email
	loginParams := params.RLogin
	users, err := a.Context.DBM.GetUsersByEmail([]string{*loginParams.Email})
	if err != nil {
		logger.Error("cant check users present GetUsersByEmail", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(users) != 1 {
		logger.Error("can't find user")
		return forbiddenResponse(403, "user or password wrong")
	}
	user := users[0]

	// check password
	//password := strings.ToLower(*loginParams.Password)
	password := *loginParams.Password
	if !user.IsEqualEPass(password) {
		logger.Error("cant crypt password", "err", err)
		return forbiddenResponse(403, "user or password wrong")
	}

	if user.EmailVerified == false {
		c := int32(428)
		logger.Error("email verification", "err", "email is not confirmed")
		return auth.NewPostAuthLoginPreconditionRequired().WithPayload(
			&models.Error{Code: &c, Message: "confirmed email is required"},
		)
	}

	// auth token
	token, err := util.GenerateToken()
	if err != nil {
		logger.Error("cant create user", "user", fmt.Sprintf("%+v", user))
		return errorResponse(500, err.Error())
	}

	session := &dto.UserSession{
		UserID:    user.ID,
		AuthToken: token,
	}
	if err := a.Context.DBM.CreateUserSession(session); err != nil {
		logger.Error("error", "err", err.Error())
		return errorResponse(500, err.Error())
	}

	// save user locale
	headerLocale := strings.Split(strings.Split(params.HTTPRequest.Header.Get("Accept-Language"), ",")[0], "-")[0]
	if len(headerLocale) != 0 {
		user.Locale = headerLocale
		if err := a.Context.DBM.UpdateUser(user); err != nil {
			logger.Error("update user error", "err", err.Error())
			return errorResponse(500, err.Error())
		}
	}

	return newCookieAuthRespond(auth.NewPostAuthLoginOK().WithPayload(&models.ALogin{
		AuthToken: &models.AuthToken{
			Token: &token,
		},
	}), token)
}

// Logout TBD
func (a *Auth) Logout(params auth.PostAuthLogoutParams, principal interface{}) middleware.Responder {
	userContext := principal.(*sauth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("logout")

	errorResponse := func(code int) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthLogoutDefault(code).WithPayload(
			&models.Error{Code: &c},
		)
	}

	err := a.Context.DBM.DestroyUserSessionByTokens([]string{userContext.AToken()})
	if err != nil {
		logger.Error("DestroyUserSession", "err", err)
		return errorResponse(500)
	}

	return newCookieAuthRespond(auth.NewPostAuthLogoutNoContent(), "")
}

// NewConfirmEmailMail TBD
func NewConfirmEmailMail(templates *template.Template, to *dto.User, confirmURL string) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.auth.confirm_email.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title":  T("api.auth.confirm_email.title"),
			"Link":   confirmURL,
			"Button": T("api.auth.confirm_email.button"),
			"Footer": T("api.auth.confirm_email.footer.info"),
		},
	}

	return subject, util.Parse(templates, "auth_body", props)
}

// NewResetPasswordMail TBD
func NewResetPasswordMail(templates *template.Template, to *dto.User, confirmURL string) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.auth.reset_password.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title":  T("api.auth.reset_password.title"),
			"Link":   confirmURL,
			"Button": T("api.auth.reset_password.button"),
			"Footer": T("api.auth.reset_password.footer.info"),
		},
	}

	return subject, util.Parse(templates, "auth_body", props)
}
