package resource

import (
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

// ChangeStatus TBD
func (t *Task) ChangeStatus(params tasks.PostTaskIDChangeStatusParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("change status")

	var (
		DBM       = t.Context.DBM
		newStatus = params.RChangeTaskStatus.NewStatus
		taskID    = params.ID
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewPostTaskIDChangeStatusDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenErrorResponse := tasks.NewPostTaskIDChangeStatusForbidden

	// TODO check authorization
	var task *dto.Task
	changedFields := dto.TaskFieldsList{dto.TaskFieldStatus}

	if tasksList, err := DBM.GetTasksByIDs([]int64{taskID}); err != nil {
		logger.Error("cant get task by id", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(tasksList) != 1 {
		err := fmt.Errorf("cant find the task: %d", taskID)
		logger.Error("cant get task by id", "err", err)
		return tasks.NewPostTaskIDChangeStatusNotFound()
	} else {
		task = tasksList[0]
	}

	// check task access
	accessChecker := NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	ok, err := accessChecker.IsUserRelatedToTask(userContext.User.ID, task)
	if err != nil {
		logger.Error("IsUserRelatedToTask", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to task: %d => %d", userContext.User.ID, task.ID)
		logger.Error("user cant change task status", "err", err)
		return forbiddenErrorResponse()
	}

	// get task users
	taskUserList, err := DBM.GetUsersByIDs(task.GetUsersIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// status
	if !dto.IsTaskStatusValid(newStatus) {
		return tasks.NewPostTaskIDChangeStatusUnprocessableEntity()
	}

	if task.Status == dto.TaskStatus(newStatus) {
		logger.Debug("task is already in", "status", newStatus)
		return tasks.NewPostTaskIDChangeStatusNoContent()
	}

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
						TaskID:        &task.ID,
						TaskTitle:     task.Title,
						ChangedStatus: true,
						NewStatus:     dto.TaskStatus(newStatus),
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
			subject, body := TaskChangedStatusMail(t.Templates, userContext.User, nUser, task)
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

	task.Status = dto.TaskStatus(newStatus)

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

	return tasks.NewPostTaskIDChangeStatusNoContent()
}
