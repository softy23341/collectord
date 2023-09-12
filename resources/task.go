package resource

import (
	"html/template"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/util"
)

// Task TBD
type Task struct {
	Context   Context
	Templates *template.Template
}

// NewTaskOnYouMail TBD
func NewTaskOnYouMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("New task for you")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.task.on_you.title",
				map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}

// TaskIsNotOnYouNowMail TBD
func TaskIsNotOnYouNowMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("The task is not on you now")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.task.not_on_you.title",
				map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}

// TaskIsOnYouNowMail TBD
func TaskIsOnYouNowMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("The task is not on you")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.task.assigned_on_you",
				map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}

// TaskChangedStatusMail TBD
func TaskChangedStatusMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("New task status")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.task.status_changed",
				map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}

// TaskChangedArchiveMail TBD
func TaskChangedArchiveMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	var subject string
	var body string
	T := util.GetTranslationFunc(to.Locale)

	if task.Archive {
		subject = T("The task now in archive")
		body = T("api.task.archived",
			map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title})
	} else {
		subject = T("The task now active")
		body = T("api.task.unarchived",
			map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title})
	}

	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title":  body,
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}

// TaskWasDeletedMail TBD
func TaskWasDeletedMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("The task was deleted")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.task.deleted",
				map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}

// TaskWasChangedMail TBD
func TaskWasChangedMail(templates *template.Template, from, to *dto.User, task *dto.Task) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("The task was changed")
	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title": T("api.task.changed",
				map[string]interface{}{"SenderName": from.FullName(), "TaskName": task.Title}),
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "task_body", props)
}
