package main

import (
	"errors"
	"fmt"
	"path"
	"reflect"

	pgx "github.com/jackc/pgx"

	"git.softndit.com/collector/backend/cleaver"
	cleaverclient "git.softndit.com/collector/backend/cleaver/client"
	"git.softndit.com/collector/backend/dal"
	dalpg "git.softndit.com/collector/backend/dal/pg"
	"git.softndit.com/collector/backend/dto"
	mediavariant "git.softndit.com/collector/backend/resources/media"
	"git.softndit.com/collector/backend/services"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

type mediaReprocessor struct {
	DBM           dal.Manager
	Storage       services.FileStorage
	Log           log15.Logger
	CleaverClient cleaverclient.ConnectClient
}

func main() {
	// logger
	logger := log15.New("app", "reporocess images")

	// DBM
	DBM, err := dalpg.NewManager(&dalpg.ManagerContext{
		Log: logger,
		PoolConfig: &pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Database: "collector_development",
				User:     "collector_app",
				Password: "collectordevpassstrong",
			},
			MaxConnections: 10,
		},
	})
	if err != nil {
		panic(err)
	}

	// storage
	storage, err := services.NewFSFiler(&services.FSStorageContext{
		BasePath:       "/collector-uploads/files/",
		BaseURL:        "/files/",
		BaseNamePrefix: "dev",
		Log:            logger,
	})
	if err != nil {
		panic(err)
	}

	// cleaver
	cleaverClient := cleaverclient.NewRMQClient("amqp://guest:guest@localhost:5672/")
	if err := cleaverClient.Connect(); err != nil {
		panic(err)
	}

	mr := &mediaReprocessor{
		DBM:           DBM,
		Storage:       storage,
		Log:           logger,
		CleaverClient: cleaverClient,
	}

	mediasForReprocessing := dto.MediaTypeSet{
		dto.MediaTypePhotoObject,
		dto.MediaTypeAvatar,
	}

	mr.processMedias(mediasForReprocessing)

}

func (m *mediaReprocessor) processMedias(types dto.MediaTypeSet) {
	perPage := 20

	for paginator := dto.NewPagePaginator(0, int16(perPage)); ; paginator.NextPage() {
		medias, err := m.DBM.GetMediaByPage(types, paginator)
		if err != nil {
			panic(err)
		}

		for _, media := range medias {
			if err := m.processMedia(media); err != nil {
				m.Log.Error("processMedia err", "err", err)
			}

		}
		if mlen := len(medias); mlen == 0 || mlen < perPage {
			break
		}

	}
}

func (m *mediaReprocessor) processMedia(media *dto.Media) error {
	if media.Type.IsPhoto() {
		var mediaLocation *services.FileLocation

		variants := media.Photo.Variants
		var originalVariant *dto.PhotoVariant
		oldMediaVariants := make(map[string]*dto.PhotoVariant, len(variants))

		for _, variant := range variants {
			// lets check variant exist
			variantMediaLocation, err := m.Storage.GetFileLocation(variant.URI)
			if err != nil || !m.Storage.IsExist(variantMediaLocation) {
				// TODO actualize variants
				m.Log.Warn("variant does not exist", "URL", variant.URI, "media_id", media.ID)
				continue
			}

			oldMediaVariants[variantMediaLocation.FullRelativePath()] = variant

			// skip original
			if variant.Original() {
				originalVariant = variant
				mediaLocation = variantMediaLocation
			}
		}

		if originalVariant == nil {
			m.Log.Error("cant find original for media", "id", media.ID)
			return errors.New("cant find original for media")
		}

		// get new variants
		typedTransforms, err := mediavariant.TransformsByMediaType(media.Type)
		if err != nil {
			m.Log.Debug("can't get transforms by type", "typo", media.Type)
			return err
		}

		var photoVariants []*dto.PhotoVariant

		if len(typedTransforms) > 0 {
			var transforms []*cleaver.Transform

			targetToPhotoVariant := make(map[string]dto.PhotoVariantType)

			for _, tTransform := range typedTransforms {
				transform := tTransform.Transform
				task := *transform

				taskTarget := mediaLocation.FullRelativePathWithPrefix(transform.Target)
				if variant, find := oldMediaVariants[taskTarget]; find {
					m.Log.Debug("variant already present", "variant target", taskTarget)
					photoVariants = append(photoVariants, variant)
					continue
				}

				task.Target = taskTarget
				transforms = append(transforms, &task)
				targetToPhotoVariant[task.Target] = tTransform.Type
			}

			resizeResults, err := m.CleaverClient.Resize(&cleaver.ResizeTask{
				Source:     mediaLocation.FullRelativePath(),
				Transforms: transforms,
			})
			if err != nil {
				m.Log.Error("cant process media", "err", err)
				return err
			}
			m.Log.Debug("media processing", "r", fmt.Sprintf("%+v", resizeResults))

			if len(resizeResults) > 0 {
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
			}
		}

		photoVariants = append(photoVariants, originalVariant)

		media.Photo.Variants = photoVariants

		// it's ok
		if !reflect.DeepEqual(photoVariants, variants) {
			m.Log.Debug("save media", "id", media.ID)
			err = m.DBM.UpdateMedia(media)
			if err != nil {
				return err
			}
		}

		// delete unused variants
		for _, oldVariant := range oldMediaVariants {
			find := false
			for _, newVariant := range photoVariants {
				if oldVariant.URI == newVariant.URI {
					find = true
					break
				}
			}

			if !find {
				m.Log.Debug("remove", "uri", oldVariant.URI)
				variantMediaLocation, _ := m.Storage.GetFileLocation(oldVariant.URI)
				_ = m.Storage.Remove(variantMediaLocation)
			}
		}
	}
	return nil
}
