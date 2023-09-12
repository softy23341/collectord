package resource

import (
	"html/template"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/util"
)

// Invite TBD
type Invite struct {
	Context   Context
	Templates *template.Template
}

// NewTeamInviteMail TBD
func NewTeamInviteMail(templates *template.Template, from, to *dto.User) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.invite.new.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.invite.new.title",
				map[string]interface{}{"SenderName": from.FullName()}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "invite_body", props)
}

// NewTeamInviteMail TBD
func InviteWasAcceptedMail(templates *template.Template, from, to *dto.User) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.invite.accepted.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.invite.accepted.title",
				map[string]interface{}{"SenderName": from.FullName()}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "invite_body", props)
}

// NewTeamInviteMail TBD
func InviteWasRejectedMail(templates *template.Template, from, to *dto.User) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.invite.rejected.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.invite.rejected.title",
				map[string]interface{}{"SenderName": from.FullName()}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "invite_body", props)
}

// NewTeamInviteMail TBD
func InviteWasCancelledMail(templates *template.Template, from, to *dto.User) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.invite.cancelled.subject")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.invite.cancelled.title",
				map[string]interface{}{"SenderName": from.FullName()}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "invite_body", props)
}
