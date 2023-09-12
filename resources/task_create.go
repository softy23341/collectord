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

// CreateTask TBD
func (t *Task) CreateTask(params tasks.PostTaskParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("create task")

	var (
		DBM       = t.Context.DBM
		inputTask = params.RCreateTask.Task
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewPostTaskDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenErrorResponse := tasks.NewPostTaskForbidden

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
			logger.Error("user cant create task here", "err", err)
			return forbiddenErrorResponse()
		}
	}

	// construct task
	task := &dto.Task{
		Title:          inputTask.Title,
		Description:    inputTask.Description,
		CreatorUserID:  userContext.User.ID,
		AssignedUserID: inputTask.AssignedUserID,

		Archive:    inputTask.Archive,
		UserUniqID: *params.RCreateTask.ClientUniqID,
	}

	// status
	if !dto.IsTaskStatusValid(inputTask.Status) {
		logger.Error("invalid status", "status", inputTask.Status)
		return tasks.NewPostTaskUnprocessableEntity()
	}
	task.Status = dto.TaskStatus(inputTask.Status)

	// deadline
	deadline, err := models.NewDtoDate(inputTask.Deadline)
	if err != nil {
		logger.Error("invalid deadline", "deadline err", err.Error())
		return tasks.NewPostTaskUnprocessableEntity()
	}
	task.Deadline = deadline

	if !task.IsValid() {
		logger.Error("task is not valid")
		return tasks.NewPostTaskUnprocessableEntity()
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return defaultErrorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// create task
	if err := tx.CreateTask(task); err != nil {
		logger.Error("CreateTask", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// users
	taskUsersIDs := make([]int64, 0, 2)
	taskUsersIDs = append(taskUsersIDs, userContext.User.ID)
	if assignedUserID := inputTask.AssignedUserID; assignedUserID != nil && userContext.User.ID != *assignedUserID {
		taskUsersIDs = append(taskUsersIDs, *assignedUserID)
	}

	taskUserList, err := tx.GetUsersByIDs(taskUsersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskUserList) != len(util.Int64Slice(taskUsersIDs).Uniq()) {
		err := errors.New("cant find all users")
		logger.Error("getUsersByIds", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// objects
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
	taskMediasList, err := tx.GetMediasByIDs(inputTask.MediasIds)
	if err != nil {
		logger.Error("GetMediasByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(taskMediasList) != len(inputTask.MediasIds) {
		logger.Error("cant find all medias", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	ok := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
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

	var (
		events dto.EventList
		jobs   []delayedjob.Job
	)

	if task.AssignedUserID != nil {
		systemUser, err := tx.GetSystemUser()
		if err != nil {
			logger.Error("get system user", "err", err)
			return defaultErrorResponse(500, err.Error())
		}

		newTaskFor := &dto.Message{
			UserID:     systemUser.ID,
			UserUniqID: util.NextUniqID(),
			PeerID:     *task.AssignedUserID,
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

						NewTask: true,
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

		assignedUser, found := taskUserList.IDToUser()[*task.AssignedUserID]
		if !found {
			logger.Error("something is wrong cant find user")
			return defaultErrorResponse(500, "")
		}

		jobs = append(jobs, delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := NewTaskOnYouMail(t.Templates, userContext.User, assignedUser, task)
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

	return tasks.NewPostTaskOK().WithPayload(&models.AGetTask{
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
