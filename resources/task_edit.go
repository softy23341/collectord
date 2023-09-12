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

// EditTask TBD
func (t *Task) EditTask(params tasks.PutTaskIDParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("edit task")

	var (
		DBM       = t.Context.DBM
		inputTask = params.REditTask.Task
		taskID    = params.ID

		changedArchive   = false
		newArchiveStatus bool

		newStatus     dto.TaskStatus
		changedStatus = false

		newAssignedUserID   *int64
		oldUserID           *int64
		changedAssignedUser = false
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewPutTaskIDDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenErrorResponse := tasks.NewPutTaskIDForbidden

	var task *dto.Task
	changedFields := dto.TaskAllFields.Del(dto.TaskFieldID)

	if tasksList, err := DBM.GetTasksByIDs([]int64{taskID}); err != nil {
		logger.Error("cant get task by id", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(tasksList) != 1 {
		err := fmt.Errorf("cant find the task: %d", taskID)
		logger.Error("cant get task by id", "err", err)
		return tasks.NewPutTaskIDNotFound()
	} else {
		task = tasksList[0]
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

	if assignedUserID := inputTask.AssignedUserID; assignedUserID != nil {
		rootsList, err := DBM.GetMainUserRoot(userContext.User.ID)
		if err != nil {
			logger.Error("cant create task", "err", err)
			return defaultErrorResponse(500, err.Error())
		}

		ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
			IsUserInRoots(*assignedUserID, rootsList.IDs())
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user is not in root ids: %d => %v", *assignedUserID, rootsList.IDs())
			logger.Error("user cant edit task here", "err", err)
			return forbiddenErrorResponse()
		}
	}

	if util.PInt64HasChnaged(task.AssignedUserID, inputTask.AssignedUserID) {
		changedAssignedUser = true
		oldUserID = task.AssignedUserID
		newAssignedUserID = inputTask.AssignedUserID
	}

	// status
	if !dto.IsTaskStatusValid(inputTask.Status) {
		logger.Error("cant parse status", "status", inputTask.Status)
		return tasks.NewPutTaskIDUnprocessableEntity()
	}

	changedStatus = task.Status != dto.TaskStatus(inputTask.Status)
	newStatus = dto.TaskStatus(inputTask.Status)
	task.Status = dto.TaskStatus(inputTask.Status)

	// deadline
	deadline, err := models.NewDtoDate(inputTask.Deadline)
	if err != nil {
		logger.Error("cant parse deadline", "err", err)
		return tasks.NewPutTaskIDUnprocessableEntity()
	}
	task.Deadline = deadline

	task.Title = inputTask.Title
	task.Description = inputTask.Description
	task.CreatorUserID = userContext.User.ID
	task.AssignedUserID = inputTask.AssignedUserID

	changedArchive = task.Archive != inputTask.Archive
	newArchiveStatus = inputTask.Archive

	task.Archive = inputTask.Archive

	if !task.IsValid() {
		logger.Error("task is not valid")
		return tasks.NewPutTaskIDUnprocessableEntity()
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return defaultErrorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// update task
	if task, err = tx.UpdateTask(taskID, task, changedFields); err != nil {
		logger.Error("Update task", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// users
	taskUsersIDs := make([]int64, 0, 3)
	taskUsersIDs = append(taskUsersIDs, userContext.User.ID)
	if assignedUserID := inputTask.AssignedUserID; assignedUserID != nil && userContext.User.ID != *assignedUserID {
		taskUsersIDs = append(taskUsersIDs, *assignedUserID)
	}
	if oldUserID != nil && !util.Int64Slice(taskUsersIDs).ToSet().Contains(*oldUserID) {
		taskUsersIDs = append(taskUsersIDs, *oldUserID)
	}

	taskUserList, err := tx.GetUsersByIDs(taskUsersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskUserList) != len(taskUsersIDs) {
		err := errors.New("cant find all users")
		logger.Error("getUsersByIds", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// objects
	if err := tx.DeleteTaskObjectsRefs(task.ID); err != nil {
		logger.Error("DeleteTaskObjectsRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	taskObjectsList, err := tx.GetObjectsByIDs(inputTask.ObjectsIds)
	if err != nil {
		logger.Error("GetObjectsListByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskObjectsList) != len(inputTask.ObjectsIds) {
		err := errors.New("cant find all objects")
		logger.Error("GetObjectsListByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// check users access rights
	for _, userID := range taskUsersIDs {
		ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
			HasUserRightsForObjects(
				userID,
				dto.RightEntityLevelRead,
				taskObjectsList.GetIDs(),
			)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user cant touch objects: %d => %v", userID, taskObjectsList.GetIDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}
	}

	if err := tx.CreateTaskObjectsRefs(task.ID, inputTask.ObjectsIds); err != nil {
		logger.Error("CreateTaskObjectsRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// collections
	if err := tx.DeleteTaskCollectionsRefs(task.ID); err != nil {
		logger.Error("DeleteTaskCollectionsRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	taskCollectionsList, err := tx.GetCollectionsByIDs(inputTask.CollectionsIds)
	if err != nil {
		logger.Error("GetCollectionsByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskCollectionsList) != len(inputTask.CollectionsIds) {
		err := errors.New("cant find all collections")
		logger.Error("cant find all collections", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	for _, userID := range taskUsersIDs {
		ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
			HasUserRightsForCollections(
				userID,
				dto.RightEntityLevelRead,
				taskCollectionsList.GetIDs(),
			)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user cant touch collections: %d => %v", userID, taskCollectionsList.GetIDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}
	}

	if err := tx.CreateTaskCollectionsRefs(task.ID, inputTask.CollectionsIds); err != nil {
		logger.Error("CreateTaskCollectionsRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// groups
	if err := tx.DeleteTaskGroupsRefs(task.ID); err != nil {
		logger.Error("DeleteTaskGroupsRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	taskGroupsList, err := tx.GetGroupsByIDs(inputTask.GroupsIds)
	if err != nil {
		logger.Error("GetGroupsByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskGroupsList) != len(inputTask.GroupsIds) {
		err := errors.New("cant find all groups")
		logger.Error("cant find all groups", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	for _, userID := range taskUsersIDs {
		ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
			HasUserRightsForGroups(
				userID,
				dto.RightEntityLevelRead,
				taskGroupsList.GetIDs(),
			)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return defaultErrorResponse(500, err.Error())
		} else if !ok {
			err := fmt.Errorf("user cant touch groups: %d => %v", userID, taskGroupsList.GetIDs())
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}
	}

	if err := tx.CreateTaskGroupsRefs(task.ID, inputTask.GroupsIds); err != nil {
		logger.Error("CreateTaskGroupsRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// medias
	if err := tx.DeleteTaskMediasRefs(task.ID); err != nil {
		logger.Error("DeleteTaskMediasRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	taskMediasList, err := tx.GetMediasByIDs(inputTask.MediasIds)
	if err != nil {
		logger.Error("GetMediasByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskMediasList) != len(inputTask.MediasIds) {
		logger.Error("cant find all medias", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	ok = NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
		IsUserOwnerOfMedias(userContext.User.ID, taskMediasList)
	if !ok {
		err := fmt.Errorf("user is not owner of medias: %d => %v",
			userContext.User.ID, taskMediasList.GetIDs())
		logger.Error("user cant create task here", "err", err)
		return forbiddenErrorResponse()
	}

	if err := tx.CreateTaskMediasRefs(task.ID, inputTask.MediasIds); err != nil {
		logger.Error("CreateTaskMediasRefs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// fetch additional medias
	additionalMediasIDs := taskCollectionsList.GetMediasIDs()
	additionalMediasIDs = append(additionalMediasIDs, taskUserList.GetAvatarsMediaIDs()...)

	additionalMedias, err := tx.GetMediasByIDs(additionalMediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskMediasList) != len(inputTask.MediasIds) {
		logger.Error("cant find all medias", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	systemUser, err := tx.GetSystemUser()
	if err != nil {
		logger.Error("get system user", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	var (
		events dto.EventList
		jobs   []delayedjob.Job
	)

	if changedAssignedUser && oldUserID != nil {
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

	if true || changedArchive || changedStatus || changedAssignedUser {
		usersIDsForNotify := util.Int64Slice(task.GetUsersIDs()).Delete(userContext.User.ID)

		if changedAssignedUser && newAssignedUserID != nil {
			usersIDsForNotify = util.Int64Slice(usersIDsForNotify).Delete(*newAssignedUserID)

			newTaskFor := &dto.Message{
				UserID:     systemUser.ID,
				UserUniqID: util.NextUniqID(),
				PeerID:     *newAssignedUserID,
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
							NewAssignedUserID:   newAssignedUserID,
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

			assignedUser, found := taskUserList.IDToUser()[*newAssignedUserID]
			if !found {
				logger.Error("something is wrong cant find user")
				return defaultErrorResponse(500, "")
			}

			jobs = append(jobs, delayedjob.NewJob(delayedjob.Immideate, func() {
				subject, body := TaskIsOnYouNowMail(t.Templates, userContext.User, assignedUser, task)
				err := t.Context.MailClient.Send(&services.Mail{
					To:      []string{assignedUser.Email},
					From:    services.SystemMailFrom,
					Subject: subject,
					Body:    body,
				})
				if err != nil {
					logger.Error("cant send email", "err", err)
				}
			}))
		}

		for _, userID := range usersIDsForNotify {
			taskMsg := &dto.ServiceMessageTask{
				TaskID:    &task.ID,
				TaskTitle: task.Title,

				ActorUserID:   userContext.User.ID,
				ActorUserName: userContext.User.FullName(),

				WasChanged:          true,
				ChangedAssignedUser: changedAssignedUser,
				ChangedStatus:       changedStatus,
				ChangedArchive:      changedArchive,
			}
			if changedAssignedUser {
				taskMsg.NewAssignedUserID = newAssignedUserID
			}
			if changedStatus {
				taskMsg.NewStatus = newStatus
			}
			if changedArchive {
				taskMsg.NewArchiveStatus = newArchiveStatus
			}
			newTaskFor := &dto.Message{
				UserID:     systemUser.ID,
				UserUniqID: util.NextUniqID(),
				PeerID:     userID,
				PeerType:   dto.PeerTypeUser,
				Typo:       dto.MessageTypeService,
				MessageExtra: dto.MessageExtra{
					Service: &dto.ServiceMessage{
						Type: dto.ServiceMessageTypeTask,
						Task: taskMsg,
					},
				},
			}

			result, err := t.Context.MessengerClient.
				NewMessageSender(tx, newTaskFor, &services.MessageInfo{}).
				Send()

			if err != nil {
				logger.Error("cant send messages", "err", err)
				return defaultErrorResponse(500, err.Error())
			}

			jobs = append(jobs, result.Jobs...)
			events = append(events, result.Events...)

			nUser, found := taskUserList.IDToUser()[userID]
			if !found {
				logger.Error("something is wrong cant find user")
				return defaultErrorResponse(500, "")
			}

			jobs = append(jobs, delayedjob.NewJob(delayedjob.Immideate, func() {
				subject, body := TaskWasChangedMail(t.Templates, userContext.User, nUser, task)
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
	}

	// end commit
	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	t.Context.EventSender.Send(events...)
	t.Context.JobPool.Enqueue(jobs...)

	allCollections := taskCollectionsList
	allGroups := taskGroupsList

	allMedias := taskMediasList
	allMedias = append(allMedias, additionalMedias...)

	//
	return tasks.NewPutTaskIDOK().WithPayload(&models.AGetTask{
		Collections:    models.NewModelCollectionList(allCollections),
		Groups:         models.NewModelGroupList(allGroups),
		ObjectsPreview: models.NewModelObjectPreviewList(taskObjectsList.Preview()),

		Medias: models.NewModelMediaList(allMedias),

		Task: models.NewModelTask(task).
			WithCollectionsIDs(taskCollectionsList.GetIDs()).
			WithGroupsIDs(taskGroupsList.GetIDs()).
			WithObjectsIDs(taskObjectsList.GetIDs()).
			WithMediasIDs(taskMediasList.GetIDs()),
		Users: models.NewModelUserList(taskUserList),
	})
}
