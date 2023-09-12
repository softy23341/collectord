package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

const userSearchByNameMaxPopTags = 15

// SearchByName TBD
func (u *User) SearchByName(params users.PostUserSearchByNameParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = u.Context.DBM

		errorResponse = func(code int) middleware.Responder {
			return users.NewPostUserSearchByNameDefault(code)
		}

		sParams         = params.RSearchUsersByName
		withPopularTags = sParams.WithPopularTags

		page, perPage = int16(0), int16(30)
	)

	logger.Debug("search by name")

	if paginator := sParams.Paginator; paginator != nil {
		if p := paginator.Page; p != nil {
			page = *p
		}
		if pp := paginator.Cnt; pp != nil {
			perPage = *pp
		}
	}

	var popularTags []string
	if withPopularTags {
		popTg, err := DBM.GetPopularUserTags(userSearchByNameMaxPopTags)
		if err != nil {
			logger.Error("GetPopularUserTags", "err", err)
			return errorResponse(500)
		}

		popularTags = popTg
	}
	ucnt, usersList, err := DBM.SearchUsersByName(userContext.User.ID, sParams.Name, page, perPage)
	if err != nil {
		logger.Error("cant SearchUsersByName", "err", err)
		return errorResponse(500)
	}

	mediasIDs := usersList.GetAvatarsMediaIDs()

	var medias dto.MediaList
	if len(mediasIDs) > 0 {
		var err error
		medias, err = DBM.GetMediasByIDs(mediasIDs)
		if err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500)
		}
	}

	return users.NewPostUserSearchByNameOK().WithPayload(&models.AUsersByName{
		Medias:          models.NewModelMediaList(medias),
		Users:           models.NewModelUserList(usersList),
		PopularTags:     popularTags,
		UsersTotalCount: ucnt,
	})
}
