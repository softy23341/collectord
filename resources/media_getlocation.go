package resource

import (
	"net/http"
	"net/url"
	"path"
	"strconv"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/restapi/operations/medias"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

const (
	mediaIDParam       = "mediaID"
	serveLocationAlias = "/serve-files"
)

type (
	mediaAuthRespond struct {
		responder    middleware.Responder
		redirectPath string
	}
)

func newMediaAuthRespond(responder middleware.Responder, redirectMediaPath string) *mediaAuthRespond {
	return &mediaAuthRespond{
		responder:    responder,
		redirectPath: redirectMediaPath,
	}
}

// WriteResponse TBD
func (c *mediaAuthRespond) WriteResponse(rw http.ResponseWriter, p runtime.Producer) {
	rw.Header().Set("X-Accel-Redirect", c.redirectPath)
	rw.Header().Set("Content-Type", "")

	c.responder.WriteResponse(rw, p)
}

// GetMediaLocation TBD
func (m *Media) GetMediaLocation(params medias.GetMediasParams) middleware.Responder {
	var (
		authToken string
		DBM       = m.Context.DBM
		logger    = m.Context.Log.New("service", "mediachecker")

		HTTPRequest = params.HTTPRequest
		mediaPath   = params.Path
	)

	successResponse := func(mediaPath string) middleware.Responder {
		logger.Debug("redirect header", "x", path.Join(serveLocationAlias, mediaPath))
		return newMediaAuthRespond(medias.NewGetMediasNoContent(), path.Join(serveLocationAlias, mediaPath))
	}

	cookie, err := HTTPRequest.Cookie(AuthCookieKey)
	if err != nil && err != http.ErrNoCookie {
		logger.Error("cant get cookie", "err", err)
		return medias.NewGetMediasDefault(500)
	}

	if cookie != nil {
		authToken = cookie.Value
	}
	// try header
	if authToken == "" {
		authToken = HTTPRequest.Header.Get(AuthHeaderKey)
	}

	mediaURL, err := url.Parse(mediaPath)
	if err != nil {
		logger.Error("cant parse mediaPath to URL", "err", err)
		return medias.NewGetMediasDefault(500)
	}

	mediaIDstr := mediaURL.Query().Get(mediaIDParam)
	if mediaIDstr == "" {
		logger.Error("cant find mediaIDstr", "str", mediaURL.Query())
		return medias.NewGetMediasDefault(500)
	}

	mediaID, err := strconv.ParseInt(mediaIDstr, 10, 0)
	if err != nil {
		logger.Error("cant parse mediaIDstr", "str", mediaIDstr)
		return medias.NewGetMediasDefault(500)
	}

	media, err := DBM.GetMediaByIDAndVariantURI(mediaID, mediaURL.Path)
	if err != nil {
		logger.Error("GetMediaByIDAndVariantURI", "err", err)
		return medias.NewGetMediasDefault(404)
	}
	if media == nil {
		return medias.NewGetMediasDefault(404)
	}

	mediaTypeSet := dto.MediaTypeSet{
		dto.MediaTypeCollection,
		dto.MediaTypePhotoObject,
		dto.MediaTypeDocument,
		dto.MediaTypeVideo,
	}

	if !mediaTypeSet.Contain(media.Type) {
		return successResponse(mediaPath)
	}

	// check user access rights
	allow, err := DBM.CanUserGetMediaByAuthToken(authToken, media.ID)
	if err != nil {
		logger.Error("CanUserGetMediaByAuthToken", "err", err)
		return medias.NewGetMediasDefault(500)
	}
	if allow {
		return successResponse(mediaPath)
	}

	// check temporary access from chat message
	if params.MessageID != nil {
		allow, err := DBM.CheckUserTmpAccessToMedia(authToken, *params.MessageID, mediaID)
		if err != nil {
			logger.Error("CheckUserTmpAccessToMedia", "err", err)
			return medias.NewGetMediasDefault(500)
		}
		if allow {
			return successResponse(mediaPath)
		}
	}

	logger.Error("user is not allowed to see media", "err", err)
	return medias.NewGetMediasForbidden()
}
