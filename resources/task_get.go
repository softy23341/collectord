package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/tasks"
	"github.com/go-openapi/runtime/middleware"
)

// GetTaskByID TBD
func (t *Task) GetTaskByID(params tasks.GetTaskIDParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get tasks")

	var (
		DBM    = t.Context.DBM
		taskID = params.ID
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewGetTaskIDDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// Get task
	var task *dto.Task
	if tasksList, err := DBM.GetTasksByIDs([]int64{taskID}); err != nil {
		logger.Error("cant get task by id", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(tasksList) != 1 {
		err := fmt.Errorf("cant find the task: %d", taskID)
		logger.Error("cant get task by id", "err", err)
		return tasks.NewGetTaskIDNotFound()
	} else {
		task = tasksList[0]
	}

	// task users
	taskUsersIDs := task.GetUsersIDs()
	taskUserList, err := DBM.GetUsersByIDs(taskUsersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return defaultErrorResponse(500, err.Error())
	}

	// check access rights
	ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access rights")).
		IsUserRelatedToTask(userContext.User.ID, task)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to task ids: %d => %d", userContext.User.ID, task.ID)
		logger.Error("user get task", "err", err)
		return tasks.NewGetTaskIDForbidden()
	}

	// collections
	taskCollectionsList, err := DBM.GetTaskCollections(taskID)
	if err != nil {
		logger.Error("cant get task collections", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// groups
	taskGroupsList, err := DBM.GetTaskGroups(taskID)
	if err != nil {
		logger.Error("cant get task groups", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// objects
	taskObjectsList, err := DBM.GetTaskObjects(taskID)
	if err != nil {
		logger.Error("cant get task objects", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// medias
	taskMediasList, err := DBM.GetTaskMedias(taskID)
	if err != nil {
		logger.Error("cant get task medias", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	// additional
	objectCollections, err := DBM.GetCollectionsByIDs(taskObjectsList.GetCollectionsIDs())
	if err != nil {
		logger.Error("cant get objects collections", "err", err)
		return defaultErrorResponse(500, err.Error())
	}
	allCollections := append(taskCollectionsList, objectCollections...)
	allGroups := taskGroupsList

	additionalMediasIDs := append(
		allCollections.GetMediasIDs(),
		taskUserList.GetAvatarsMediaIDs()...,
	)

	additionalMedias, err := DBM.GetMediasByIDs(additionalMediasIDs)
	if err != nil {
		logger.Error("cant get additional medias by ids", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	allMedias := append(taskMediasList, additionalMedias...)

	return tasks.NewGetTaskIDOK().WithPayload(&models.AGetTask{
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
