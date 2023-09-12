package resource

import (
	"errors"
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/tasks"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// Archive TBD
func (t *Task) Archive(params tasks.PostTaskIDArchiveParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("archive")

	var (
		DBM     = t.Context.DBM
		archive = params.RArchiveTask.Archive
		taskID  = params.ID
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewPostTaskIDArchiveDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenErrorResponse := tasks.NewPostTaskIDArchiveForbidden

	var task *dto.Task
	changedFields := dto.TaskFieldsList{dto.TaskFieldArchive}

	if tasksList, err := DBM.GetTasksByIDs([]int64{taskID}); err != nil {
		logger.Error("cant get task by id", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(tasksList) != 1 {
		err := fmt.Errorf("cant find the task: %d", taskID)
		logger.Error("cant get task by id", "err", err)
		return tasks.NewPostTaskIDArchiveNotFound()
	} else {
		task = tasksList[0]
	}

	// get task users
	taskUserList, err := DBM.GetUsersByIDs(task.GetUsersIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// check task access
	ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
		IsUserTaskOwner(userContext.User.ID, task)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if !ok {
		err := errors.New("user is not task owner")
		logger.Error("user cant edit task", "err", err)
		return forbiddenErrorResponse()
	}

	if task.Archive == archive {
		logger.Debug("task already in", "archive", archive)
		return tasks.NewPostTaskIDArchiveNoContent()
	}
	task.Archive = archive

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("dbm.begintx fail", "err", err)
		return defaultErrorResponse(500, err.Error())
	}
	defer tx.Rollback()

	var (
		events dto.EventList
		jobs   []delayedjob.Job
	)

	systemUser, err := tx.GetSystemUser()
	if err != nil {
		logger.Error("get system user", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	usersIDsForNotify := util.Int64Slice(task.GetUsersIDs()).Delete(userContext.User.ID)
	for _, userID := range usersIDsForNotify {
		newTaskFor := &dto.Message{
			UserID:     systemUser.ID,
			UserUniqID: util.NextUniqID(),
			PeerID:     userID,
			PeerType:   dto.PeerTypeUser,
			Typo:       dto.MessageTypeService,
			MessageExtra: dto.MessageExtra{
				Service: &dto.ServiceMessage{
					Type: dto.ServiceMessageTypeTask,
					Task: &dto.ServiceMessageTask{
						TaskID:    &task.ID,
						TaskTitle: task.Title,

						ActorUserID:   userContext.User.ID,
						ActorUserName: userContext.User.FullName(),

						WasChanged:       true,
						ChangedArchive:   true,
						NewArchiveStatus: archive,
					},
				},
			},
		}

		result, err := t.Context.MessengerClient.
			NewMessageSender(tx, newTaskFor, &services.MessageInfo{}).
			Send()
		if err != nil {
			logger.Error("cant send messages", "error", err)
			return defaultErrorResponse(500, err.Error())
		}

		events = append(events, result.Events...)
		jobs = append(jobs, result.Jobs...)

		nUser, found := taskUserList.IDToUser()[userID]
		if !found {
			logger.Error("something is wrong cant find user")
			return defaultErrorResponse(500, "")
		}

		jobs = append(jobs, delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := TaskChangedArchiveMail(t.Templates, userContext.User, nUser, task)
			err := t.Context.MailClient.Send(&services.Mail{
				To:      []string{nUser.Email},
				From:    services.SystemMailFrom,
				Subject: subject,
				Body:    body,
			})
			if err != nil {
				logger.Error("cant send email", "err", err)
			}
		}))
	}

	// update task
	if _, err := tx.UpdateTask(taskID, task, changedFields); err != nil {
		logger.Error("update task", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	t.Context.EventSender.Send(events...)
	t.Context.JobPool.Enqueue(jobs...)

	return tasks.NewPostTaskIDArchiveNoContent()
}
