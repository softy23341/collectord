package resource

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/inconshreveable/log15.v2"

	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/errs"
	"git.softndit.com/collector/backend/util"
)

// Object TBD
type Object struct {
	Context Context
}

func (o *Object) validateCurrencies(currenciesIDs []int64) error {
	if len(currenciesIDs) == 0 {
		return nil
	}
	currenciesIDs = util.UniqInt64(currenciesIDs)
	currencies, err := o.Context.DBM.GetCurrenciesByIDs(currenciesIDs)
	if err != nil {
		o.Context.Log.Error("GetCurrenciesByIDs error", "error", err)
		return err
	}
	if len(currencies) != len(currenciesIDs) {
		err = fmt.Errorf("can't find currency ids: %+v ", currenciesIDs)
		o.Context.Log.Error("GetCurrenciesByIDs error", "error", err)
		return err
	}
	return nil
}

// material
type materialManager struct {
	DBM    dal.TxManager
	log    log15.Logger
	rootID int64
}

func (m *materialManager) createMaterialsByNormalNames(names []string) (dto.MaterialList, error) {
	var outMaterials dto.MaterialList
	for _, materialName := range names {
		if strings.TrimSpace(materialName) != "" {
			material, err := m.DBM.GetOrCreateMaterialByNormalName(&dto.Material{
				Name:       materialName,
				NormalName: util.NormalizeString(materialName),
				RootID:     &m.rootID,
			})
			if err != nil {
				m.log.Error("create or find object materials", "error", err)
				return nil, err
			}
			outMaterials = append(outMaterials, material)
		}
	}
	return outMaterials, nil
}

// badge
type badgeManager struct {
	DBM    dal.TxManager
	log    log15.Logger
	rootID int64
}

func (m *badgeManager) createBadgesByNormalNames(badges dto.InputBadgeList) (dto.BadgeList, error) {
	var outBadges dto.BadgeList
	for _, badge := range badges {
		if strings.TrimSpace(badge.Name) != "" {
			badge, err := m.DBM.GetOrCreateBadgeByNormalNameAndColor(&dto.Badge{
				Name:       badge.Name,
				NormalName: util.NormalizeString(badge.Name),
				Color:      badge.Color,
				RootID:     &m.rootID,
			})
			if err != nil {
				m.log.Error("create or find object badges", "error", err)
				return nil, err
			}
			outBadges = append(outBadges, badge)
		}
	}
	return outBadges, nil
}

// actor
type actorManager struct {
	DBM    dal.TxManager
	log    log15.Logger
	rootID int64
}

func (m *actorManager) createActorsByNormalNames(names []string) (dto.ActorList, error) {
	var outActors dto.ActorList
	for _, actorName := range names {
		if strings.TrimSpace(actorName) != "" {
			actor, err := m.DBM.GetOrCreateActorByNormalName(&dto.Actor{
				Name:       actorName,
				NormalName: util.NormalizeString(actorName),
				RootID:     &m.rootID,
			})
			if err != nil {
				m.log.Error("create or find object actors", "error", err)
				return nil, err
			}

			outActors = append(outActors, actor)
		}
	}
	return outActors, nil
}

type originLocationManager struct {
	DBM    dal.TxManager
	log    log15.Logger
	rootID int64
}

func (m *originLocationManager) createOriginLocationsByNormalNames(names []string) (dto.OriginLocationList, error) {
	var outOriginLocations dto.OriginLocationList
	for _, originLocationName := range names {
		if strings.TrimSpace(originLocationName) != "" {
			originLocation, err := m.DBM.GetOrCreateOriginLocationByNormalName(&dto.OriginLocation{
				Name:       originLocationName,
				NormalName: util.NormalizeString(originLocationName),
				RootID:     &m.rootID,
			})
			if err != nil {
				m.log.Error("create or find object originLocations", "error", err)
				return nil, err
			}
			outOriginLocations = append(outOriginLocations, originLocation)
		}
	}
	return outOriginLocations, nil
}

type mediaManager struct {
	DBM dal.Manager
	log log15.Logger
}

func (m *mediaManager) getMediasByIDs(mediaIDs []int64) (dto.MediaList, error) {
	medias, err := m.DBM.GetMediasByIDs(mediaIDs)
	if err != nil {
		return nil, errs.Internal.Wrap(err)
	}

	// TODO: check objectmediarefs cnt (objects can't share media)
	if len(medias) != len(mediaIDs) {
		return nil, errors.New("Can't find medias")
	}
	return medias, nil
}

func (m *mediaManager) MediasByObjectEntities(objectsIDs []int64) (dto.MediaList, map[int64][]int64, error) {
	mediaRefs, err := m.DBM.GetObjectsMediaRefs(objectsIDs)
	if err != nil {
		return nil, nil, err
	}

	medias, err := m.DBM.GetMediasByIDs(mediaRefs.UniqMediasIDs())
	if err != nil {
		return nil, nil, err
	}

	return medias, mediaRefs.ObjectIDToMediasIDs(), nil
}

// ObjectPreviewExtractorOpts TBD
type ObjectPreviewExtractorOpts struct {
	mediaExtractorOpts *objectMediasExtractorOpts
}

// ObjectPreviewExtractor TBD
type ObjectPreviewExtractor struct {
	DBM dal.TrManager
	log log15.Logger

	objectIDs []int64

	mediaExtractor          *objectMediasExtractor
	actorExtractor          *objectActorsExtractor
	originLocationExtractor *objectOriginLocationsExtractor
	badgeExtractor          *objectBadgesExtractor
	valuationsExtractor     *objectValuationsExtractor
	valuations              dto.ValuationList

	objectStatuses dto.ObjectStatusRefList

	err  error
	opts *ObjectPreviewExtractorOpts
}

// ObjectPreviewExtractorResult TBD
type ObjectPreviewExtractorResult struct {
	ObjectIDs []int64

	MediaExtractor          *objectMediasExtractor
	ActorExtractor          *objectActorsExtractor
	OriginLocationExtractor *objectOriginLocationsExtractor
	BadgeExtractor          *objectBadgesExtractor
	Valuations              dto.ValuationList

	ObjectStatuses dto.ObjectStatusRefList
}

// NewObjectPreviewExtractor TBD
func NewObjectPreviewExtractor(DBM dal.TrManager, log log15.Logger, opts *ObjectPreviewExtractorOpts) *ObjectPreviewExtractor {
	return &ObjectPreviewExtractor{
		DBM: DBM,
		log: log,

		opts: opts,
	}
}

// SetObjectsIDs TBD
func (o *ObjectPreviewExtractor) SetObjectsIDs(objectIDs []int64) *ObjectPreviewExtractor {
	o.objectIDs = objectIDs
	return o
}

// FetchMedia TBD
func (o *ObjectPreviewExtractor) FetchMedia() *ObjectPreviewExtractor {
	if o.Failed() {
		return o
	}

	// objects medias
	mediaExtractor, err := newObjectMediasExtractor(o.mhContext(), o.objectIDs, o.opts.mediaExtractorOpts)
	if err != nil {
		o.log.Error("newObjectMediasExtractor", "err", err)
		o.err = err
		return o
	}

	o.mediaExtractor = mediaExtractor
	return o
}

// FetchActors TBD
func (o *ObjectPreviewExtractor) FetchActors() *ObjectPreviewExtractor {
	if o.Failed() {
		return o
	}

	actorExtractor, err := newObjectActorsExtractor(o.mhContext(), o.objectIDs)
	if err != nil {
		o.log.Error("GetObjectsActorRefs", "err", err)
		o.err = err
	}

	o.actorExtractor = actorExtractor

	return o
}

// FetchOriginLocations TBD
func (o *ObjectPreviewExtractor) FetchOriginLocations() *ObjectPreviewExtractor {
	if o.Failed() {
		return o
	}

	// originLocations
	originLocationExtractor, err := newObjectOriginLocationsExtractor(o.mhContext(), o.objectIDs)
	if err != nil {
		o.log.Error("GetObjectsOriginLocationRefs", "err", err)
		o.err = err
	}

	o.originLocationExtractor = originLocationExtractor

	return o
}

// FetchBadges TBD
func (o *ObjectPreviewExtractor) FetchBadges() *ObjectPreviewExtractor {
	if o.Failed() {
		return o
	}

	// badges
	badgeExtractor, err := newObjectBadgesExtractor(o.mhContext(), o.objectIDs)
	if err != nil {
		o.log.Error("GetObjectsBadgeRefs", "err", err)
		o.err = err
	}

	o.badgeExtractor = badgeExtractor

	return o
}

// FetchStatuses TBD
func (o *ObjectPreviewExtractor) FetchStatuses() *ObjectPreviewExtractor {
	if o.Failed() {
		return o
	}

	// statuses
	objectStatuses, err := o.DBM.GetCurrentObjectsStatusesRefs(o.objectIDs)
	if err != nil {
		o.log.Error("GetCurrentObjectsStatusesRefs", "err", err)
		o.err = err
	}

	o.objectStatuses = objectStatuses

	return o
}

// FetchValuations TBD
func (o *ObjectPreviewExtractor) FetchValuations() *ObjectPreviewExtractor {
	if o.Failed() {
		return o
	}

	// objects valuations
	valuationsExtractor, err := newObjectValuationsExtractor(o.mhContext(), o.objectIDs)
	if err != nil {
		o.log.Error("newObjectValuationsExtractor", "err", err)
		o.err = err
		return o
	}

	o.valuations = valuationsExtractor.GetObjectValuations()
	return o
}

type valuationsManager struct {
	DBM dal.Manager
	log log15.Logger
}

func (m *valuationsManager) createValuations(objectID int64, valuations dto.ValuationList) (dto.ValuationList, error) {
	var outValuations dto.ValuationList

	for _, val := range valuations {
		valuation, err := m.DBM.CreateValuation(&dto.Valuation{
			ObjectID:   objectID,
			Name:       util.NormalizeString(val.Name),
			Comment:    util.NormalizeString(val.Comment),
			Date:       val.Date,
			Price:      val.Price,
			CurrencyID: val.CurrencyID,
		})
		if err != nil {
			m.log.Error("create valuation", "error", err)
			return nil, err
		}
		outValuations = append(outValuations, valuation)
	}
	return outValuations, nil

}

func (m *valuationsManager) updateValuations(objectID int64, valuations dto.ValuationList) (dto.ValuationList, error) {
	var outValuations dto.ValuationList
	err := m.DBM.DeleteValuationsByObjectID(objectID)
	if err != nil {
		m.log.Error("cleaning valuations by object ID", "error", err)
		return nil, err
	}

	if valuations == nil {
		return nil, nil
	}

	for _, val := range valuations {
		valuation, err := m.DBM.CreateValuation(&dto.Valuation{
			ObjectID:   objectID,
			Name:       util.NormalizeString(val.Name),
			Comment:    util.NormalizeString(val.Comment),
			Date:       val.Date,
			Price:      val.Price,
			CurrencyID: val.CurrencyID,
		})
		if err != nil {
			m.log.Error("create valuation", "error", err)
			return nil, err
		}
		outValuations = append(outValuations, valuation)
	}
	return outValuations, nil

}

// FetchAll TBD
func (o *ObjectPreviewExtractor) FetchAll() *ObjectPreviewExtractor {
	return o.
		FetchMedia().
		FetchActors().
		FetchOriginLocations().
		FetchBadges().
		FetchStatuses().
		FetchValuations()
}

// Failed TBD
func (o *ObjectPreviewExtractor) Failed() bool {
	return o.Err() != nil
}

// Err TBD
func (o *ObjectPreviewExtractor) Err() error {
	return o.err
}

// Result TBD
func (o *ObjectPreviewExtractor) Result() (*ObjectPreviewExtractorResult, error) {
	return &ObjectPreviewExtractorResult{
		ObjectIDs: o.objectIDs,

		MediaExtractor:          o.mediaExtractor,
		ActorExtractor:          o.actorExtractor,
		OriginLocationExtractor: o.originLocationExtractor,
		BadgeExtractor:          o.badgeExtractor,
		ObjectStatuses:          o.objectStatuses,
		Valuations:              o.valuations,
	}, o.Err()
}

func (o *ObjectPreviewExtractor) mhContext() *mhContext {
	return &mhContext{
		DBM: o.DBM,
		log: o.log,
	}
}
