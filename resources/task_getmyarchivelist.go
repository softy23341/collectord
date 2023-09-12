package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/tasks"
	"github.com/go-openapi/runtime/middleware"
)

// GetMyArchiveList TBD
func (t *Task) GetMyArchiveList(params tasks.PostTaskMyArchiveListParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get my archive list")

	var (
		DBM             = t.Context.DBM
		paramsPaginator = params.RMyArchiveList.Paginator
	)

	defaultErrorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return tasks.NewPostTaskMyArchiveListDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	paginator := dto.NewDefaultPagePaginator()
	if paramsPaginator != nil {
		paginator.Cnt = *paramsPaginator.Cnt
		paginator.Page = *paramsPaginator.Page
	}

	tasksList, err := DBM.GetUserRelatedTasksArchive(userContext.User.ID, paginator)
	if err != nil {
		logger.Error("cant get users task", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	task2medias, err := DBM.GetMediasCntByTaskID(tasksList.GetIDs())
	if err != nil {
		logger.Error("cant GetMediasCntByTaskID", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	usersList, err := DBM.GetUsersByIDs(tasksList.GetUsersIDs())
	if err != nil {
		logger.Error("cant get users by ids", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	mediasList, err := DBM.GetMediasByIDs(usersList.GetAvatarsMediaIDs())
	if err != nil {
		logger.Error("cant get medias by ids", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	taskPreviewModelList := models.NewModelTaskPreviewList(tasksList)
	for _, taskPreviewModel := range taskPreviewModelList {
		mediasCnt := task2medias[*taskPreviewModel.ID]
		taskPreviewModel.WithMedias = mediasCnt > 0
	}

	return tasks.NewPostTaskMyArchiveListOK().WithPayload(&models.AMyTaskList{
		Medias:       models.NewModelMediaList(mediasList),
		TaskPreviews: taskPreviewModelList,
		Users:        models.NewModelUserList(usersList),
	})
}
