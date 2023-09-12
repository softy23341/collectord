package resource

import (
	"html/template"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/util"
)

// InviteByEmail TBD
type InviteByEmail struct {
	Context   Context
	Templates *template.Template
}

// NewTeamInviteMail TBD
func NewInviteByEmailMail(templates *template.Template, from, to *dto.User, confirmURL string) (string, string) {
	var locale string
	if len(to.Locale) != 0 {
		locale = to.Locale
	} else if len(from.Locale) != 0 {
		locale = from.Locale
	}

	T := util.GetTranslationFunc(locale)

	subject := T("api.invite_by_email.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.invite_by_email.title",
				map[string]interface{}{"SenderName": from.FullName()}),
			"Link":   confirmURL,
			"Button": T("api.invite_by_email.body.button"),
			"Footer": T("api.footer.info"),
		},
		"Html": map[interface{}]interface{}{
			"ExtraInfo": T("api.invite_by_email.body.extra_info"),
		},
	}

	return subject, util.Parse(templates, "invitebyemail_body", props)
}
