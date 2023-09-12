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

// AssignTo TBD
func (t *Task) AssignTo(params tasks.PostTaskIDAssignToParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("assign to")

	var (
		DBM       = t.Context.DBM
		newUserID = params.RAssignTaskTo.NewUserID
		oldUserID *int64
		taskID    = params.ID
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewPostTaskIDAssignToDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenErrorResponse := tasks.NewPostTaskIDAssignToForbidden

	var task *dto.Task
	changedFields := dto.TaskFieldsList{dto.TaskFieldAssignedUserID}

	if tasksList, err := DBM.GetTasksByIDs([]int64{taskID}); err != nil {
		logger.Error("cant get task by id", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(tasksList) != 1 {
		err := fmt.Errorf("cant find the task: %d", taskID)
		logger.Error("cant get task by id", "err", err)
		return tasks.NewPostTaskIDAssignToNotFound()
	} else {
		task = tasksList[0]
	}

	oldUserID = task.AssignedUserID

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

	if !util.PInt64HasChnaged(oldUserID, newUserID) {
		logger.Debug("there are no changes around")
		return tasks.NewPostTaskIDAssignToNoContent()
	}

	if newUserID != nil {
		rootsList, err := DBM.GetMainUserRoot(userContext.User.ID)
		if err != nil {
			logger.Error("cant create task", "err", err)
			return defaultErrorResponse(500, err.Error())
		}

		accessChecker := NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		ok, err = accessChecker.IsUserInRoots(*newUserID, rootsList.IDs())
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user is not in root ids: %d => %v", *newUserID, rootsList.IDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}

		// collections
		taskCollectionsList, err := DBM.GetTaskCollections(taskID)
		if err != nil {
			logger.Error("cant get task collections", "err", err)
			return defaultErrorResponse(500, err.Error())
		}
		ok, err = accessChecker.HasUserRightsForCollections(
			*newUserID,
			dto.RightEntityLevelRead,
			taskCollectionsList.GetIDs(),
		)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user cant touch collections: %d => %v", *newUserID, taskCollectionsList.GetIDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}

		// groups
		taskGroupsList, err := DBM.GetTaskGroups(taskID)
		if err != nil {
			logger.Error("cant get task groups", "err", err)
			return defaultErrorResponse(500, err.Error())
		}
		ok, err = accessChecker.HasUserRightsForGroups(
			*newUserID,
			dto.RightEntityLevelRead,
			taskGroupsList.GetIDs(),
		)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user cant touch groups: %d => %v", *newUserID, taskGroupsList.GetIDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}

		// objects
		taskObjectsList, err := DBM.GetTaskObjects(taskID)
		if err != nil {
			logger.Error("cant get task objects", "err", err)
			return defaultErrorResponse(500, err.Error())
		}
		ok, err = accessChecker.HasUserRightsForObjects(
			*newUserID,
			dto.RightEntityLevelRead,
			taskObjectsList.GetIDs(),
		)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user cant touch objects: %d => %v", *newUserID, taskObjectsList.GetIDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}
	}

	task.AssignedUserID = newUserID

	taskUsersIDs := make([]int64, 0, len(task.GetUsersIDs())+1)
	taskUsersIDs = append(taskUsersIDs, task.GetUsersIDs()...)
	if oldUserID != nil {
		taskUsersIDs = append(taskUsersIDs, *oldUserID)
	}

	taskUserList, err := DBM.GetUsersByIDs(taskUsersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
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

	if oldUserID != nil {
		newTaskForOldUser := &dto.Message{
			UserID:     systemUser.ID,
			UserUniqID: util.NextUniqID(),
			PeerID:     *oldUserID,
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

						WasChanged:          true,
						ChangedAssignedUser: true,
						NewAssignedUserID:   nil,
					},
				},
			},
		}

		result, err := t.Context.MessengerClient.
			NewMessageSender(tx, newTaskForOldUser, &services.MessageInfo{}).
			Send()

		if err != nil {
			logger.Error("cant send notifies", "err", err)
			return defaultErrorResponse(500, err.Error())
		}

		nUser, found := taskUserList.IDToUser()[*oldUserID]
		if !found {
			logger.Error("something is wrong cant find user")
			return defaultErrorResponse(500, "")
		}

		events = append(events, result.Events...)
		jobs = append(jobs, result.Jobs...)

		jobs = append(jobs, delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := TaskIsNotOnYouNowMail(t.Templates, userContext.User, nUser, task)
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

	if newUserID != nil {
		newTaskForOldUser := &dto.Message{
			UserID:     systemUser.ID,
			UserUniqID: util.NextUniqID(),
			PeerID:     *newUserID,
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

						WasChanged:          true,
						ChangedAssignedUser: true,
						NewAssignedUserID:   newUserID,
					},
				},
			},
		}

		result, err := t.Context.MessengerClient.
			NewMessageSender(tx, newTaskForOldUser, &services.MessageInfo{}).
			Send()

		if err != nil {
			logger.Error("cant send notifies", "err", err)
			return defaultErrorResponse(500, err.Error())
		}

		nUser, found := taskUserList.IDToUser()[*newUserID]
		if !found {
			logger.Error("something is wrong cant find user")
			return defaultErrorResponse(500, "")
		}

		events = append(events, result.Events...)
		jobs = append(jobs, result.Jobs...)

		jobs = append(jobs, delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := TaskIsOnYouNowMail(t.Templates, userContext.User, nUser, task)
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

	return tasks.NewPostTaskIDAssignToNoContent()
}
