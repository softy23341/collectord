package resource

import (
	"fmt"
	"io"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"time"

	cleaver "git.softndit.com/collector/backend/cleaver"
	cleaverClient "git.softndit.com/collector/backend/cleaver/client"
	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/go-openapi/runtime/middleware"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	mediavariant "git.softndit.com/collector/backend/resources/media"
	"git.softndit.com/collector/backend/restapi/operations/medias"
	"git.softndit.com/collector/backend/services"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Media TBD
type Media struct {
	Context Context
}

// CreateMedia creates media and media variants; must be sync with reprocess_images.go!
func (m *Media) CreateMedia(params medias.PostMediasParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	logger.Debug("Create media")

	errorResponse := func(code int) middleware.Responder {
		c := int32(code)
		return medias.NewPostMediasDefault(code).WithPayload(
			&models.Error{Code: &c},
		)
	}
	successResponse := func(media *dto.Media) middleware.Responder {
		return medias.NewPostMediasOK().WithPayload(&models.ANewMedia{
			Media: models.NewModelMedia(media),
		})
	}

	// check media for exists
	if findedMedia, err := m.Context.DBM.GetMediaByUserUniqID(userContext.User.ID, *params.ClientUniqID); err != nil {
		logger.Error("get media by user uniq id", "err", err)
		return errorResponse(500)
	} else if findedMedia != nil {
		logger.Debug("media already exists", "mediaID", findedMedia.ID)
		return successResponse(findedMedia)
	}

	if params.Typo == nil || !dto.MediaTypeList.Contain(dto.MediaType(*params.Typo)) {
		logger.Error("typo nil")
		return errorResponse(422)
	}

	mediaSaver := &MediaSaver{
		DBM:           m.Context.DBM,
		Logger:        logger,
		Storage:       m.Context.FileStorage,
		CleaverClient: m.Context.CleaverClient,
	}
	media, err := mediaSaver.SaveMedia(&InputFileData{
		OriginalFileName: params.File.Header.Filename,
		Typo:             dto.MediaType(*params.Typo),
		UserID:           &userContext.User.ID,
		UniqID:           *params.ClientUniqID,
		Content:          params.File.Data,
	})
	if err != nil {
		logger.Error("cant save media", "err", err)
		return errorResponse(500)
	}

	return successResponse(media)
}

// GetMediasByIDs TBD
func (m *Media) GetMediasByIDs(params medias.GetMediasByIdsParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	logger.Debug("get medias")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return medias.NewGetMediasByIdsDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	return medias.NewGetMediasByIdsForbidden()

	mediaList, err := m.Context.DBM.GetMediasByIDs(params.Ids)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	return medias.NewGetMediasByIdsOK().WithPayload(&models.AGetMediasByIds{
		Medias: models.NewModelMediaList(mediaList),
	})
}

// MediaSaver TBD
type MediaSaver struct {
	DBM           dal.Manager
	Logger        log15.Logger
	Storage       services.FileStorage
	CleaverClient cleaverClient.ConnectClient
}

// InputFileData TBD
type InputFileData struct {
	OriginalFileName string
	Typo             dto.MediaType
	UserID           *int64
	UniqID           int64
	Content          io.Reader
}

// SaveMedia TBD
func (m *MediaSaver) SaveMedia(file *InputFileData) (*dto.Media, error) {
	// file name
	ext := filepath.Ext(file.OriginalFileName)
	rand := strconv.FormatInt(rand.Int63(), 10)
	systemFileName := strconv.FormatInt(file.UniqID, 10) + "-" + rand + ext

	// save media to storage
	mediaLocation, err := m.Storage.Save(systemFileName, file.Content)
	if err != nil {
		m.Logger.Error("save file error", "err", err)
		return nil, err
	}

	// media type
	typo := file.Typo
	var media *dto.Media

	if typo.IsPhoto() {
		originalPrefix := "original_"
		typedTransforms, err := mediavariant.TransformsByMediaType(typo)
		if err != nil {
			m.Logger.Debug("can't get transforms by type", "typo", typo)
			return nil, err
		}

		originalTransform := &mediavariant.TypedTransform{
			Type: dto.PhotoVariantOriginal,
			Transform: &cleaver.Transform{
				Target:    originalPrefix,
				Geometry:  cleaver.Geometry{},
				CopyEqual: true,
			},
		}
		typedTransforms = append(typedTransforms, originalTransform)

		var photoVariants []*dto.PhotoVariant

		var transforms []*cleaver.Transform

		targetToPhotoVariant := make(map[string]dto.PhotoVariantType)

		for _, tTransform := range typedTransforms {
			transform := tTransform.Transform
			task := *transform

			prefix := ""
			if tTransform.Type != dto.PhotoVariantOriginal {
				prefix = originalPrefix
			}
			task.Target = mediaLocation.FullRelativePathWithPrefix(transform.Target + prefix)

			transforms = append(transforms, &task)
			targetToPhotoVariant[task.Target] = tTransform.Type
		}

		resizeResults, err := m.CleaverClient.Resize(&cleaver.ResizeTask{
			Source:     mediaLocation.FullRelativePath(),
			Transforms: transforms,
		})
		if err != nil {
			m.Logger.Error("cant process media", "err", err)
			return nil, err
		}
		m.Logger.Debug("media processing", "r", fmt.Sprintf("%+v", resizeResults))

		for _, r := range resizeResults {
			pVariant := &dto.PhotoVariant{
				URI:      path.Join(mediaLocation.BaseURL, r.Target),
				Type:     targetToPhotoVariant[r.Target],
				Width:    int32(r.NewGeometry.Width),
				Height:   int32(r.NewGeometry.Height),
				ByteSize: int64(r.NewByteSize),
			}
			photoVariants = append(photoVariants, pVariant)
		}

		if err := m.Storage.Remove(mediaLocation); err != nil {
			m.Logger.Error("cant remove original", "err", err)
			return nil, err
		}

		media = &dto.Media{
			UserID:     file.UserID,
			UserUniqID: &file.UniqID,
			Type:       file.Typo,
			MediaUnion: dto.MediaUnion{
				Photo: &dto.Photo{
					Name:     file.OriginalFileName,
					Variants: photoVariants,
				},
			},
		}
	}
	if typo.IsDocument() {
		media = &dto.Media{
			UserID:     file.UserID,
			UserUniqID: &file.UniqID,
			Type:       typo,
			MediaUnion: dto.MediaUnion{
				Document: &dto.Document{
					Name: file.OriginalFileName,
					URI:  mediaLocation.FullURL(),
				},
			},
		}
	}

	if err := m.DBM.CreateMedia(media); err != nil {
		m.Logger.Error("cant create media", "err", err)
		// TODO kill medias?
		return nil, err
	}
	return media, nil
}
