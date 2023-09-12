package services

import (
	"fmt"

	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

// NewEventSender TBD
func NewEventSender(eventCh chan *models.UserPayloadedEvent, sf *PayloadedEventFiller) *EventSender {
	return &EventSender{
		EventCh:              eventCh,
		PayloadedEventFiller: sf,
	}
}

// EventSender TBD
type EventSender struct {
	EventCh              chan *models.UserPayloadedEvent
	PayloadedEventFiller *PayloadedEventFiller
}

// Send TBD
func (c *EventSender) Send(events ...*dto.Event) {
	for _, event := range events {
		payloadedEvent, err := c.PayloadedEventFiller.MapToUserPayloadedEvent(event)
		if err != nil {
			continue
		}
		c.EventCh <- payloadedEvent
	}
}

// PayloadedEventFiller TBD
type PayloadedEventFiller struct {
	Log log15.Logger
	DBM dal.TrManager
}

// NewPayloadedEventFillerFromContext TBD
func NewPayloadedEventFillerFromContext(logger log15.Logger, DBM dal.TrManager) *PayloadedEventFiller {
	return &PayloadedEventFiller{
		Log: logger,
		DBM: DBM,
	}
}

// MapToUserPayloadedEvent TBD
func (sf *PayloadedEventFiller) MapToUserPayloadedEvent(e *dto.Event) (*models.UserPayloadedEvent, error) {
	payloadedEvent := models.NewModelPayloadedEvent(e)

	if e.NewEntity != nil || e.EditedEntity != nil {
		entityRef := e.NewEntity
		if e.EditedEntity != nil {
			entityRef = e.EditedEntity
		}

		if err := sf.fillPayloadedEvent(payloadedEvent, entityRef); err != nil {
			return nil, err
		}
	}

	userPayloadedEvent := &models.UserPayloadedEvent{
		UserID:         e.UserID,
		PayloadedEvent: payloadedEvent,
	}

	return userPayloadedEvent, nil
}

func (sf *PayloadedEventFiller) fillPayloadedEvent(se *models.PayloadedEvent, er *dto.EntityRef) error {
	switch typo := er.Typo; typo {
	case dto.CollectionEntityType:
		if err := sf.fillCollection(se, er.ID); err != nil {
			return err
		}
	case dto.GroupEntityType:
		if err := sf.fillGroup(se, er.ID); err != nil {
			return err
		}
	case dto.ActorEntityType:
		if err := sf.fillActor(se, er.ID); err != nil {
			return err
		}
	case dto.BadgeEntityType:
		if err := sf.fillBadge(se, er.ID); err != nil {
			return err
		}
	case dto.MaterialEntityType:
		if err := sf.fillMaterial(se, er.ID); err != nil {
			return err
		}
	case dto.OriginLocationEntityType:
		if err := sf.fillOriginLocation(se, er.ID); err != nil {
			return err
		}
	case dto.NamedIntervalEntityType:
		if err := sf.fillNamedDateInterval(se, er.ID); err != nil {
			return err
		}
	}

	return nil
}

func (sf *PayloadedEventFiller) fillCollection(se *models.PayloadedEvent, collectionID int64) error {
	// Get group
	var collection *dto.Collection
	if collections, err := sf.DBM.GetCollectionsByIDs([]int64{collectionID}); err != nil {
		return err
	} else if len(collections) != 1 {
		return fmt.Errorf("cant find collection: %d", collectionID)
	} else {
		collection = collections[0]
	}

	objectCnt, err := sf.DBM.GetObjectsCnt([]int64{collection.ID})
	if err != nil {
		return err
	}

	if collectionMediaID := collection.ImageMediaID; collectionMediaID != nil {
		medias, err := sf.DBM.GetMediasByIDs([]int64{*collectionMediaID})
		if err != nil {
			return err
		}

		se.Medias = models.NewModelMediaList(medias)
	}

	collectionGroupRefs, err := sf.DBM.GetCollectionsGroupRefs([]int64{collection.ID})
	if err != nil {
		return err
	}

	se.Collection = models.NewModelCollection(collection).
		WithObjectsCnt(objectCnt).
		WithGroupsIDs(collectionGroupRefs.GroupsIDs())

	return nil
}

func (sf *PayloadedEventFiller) fillGroup(se *models.PayloadedEvent, groupID int64) error {
	var group *dto.Group
	if groupsList, err := sf.DBM.GetGroupsByIDs([]int64{groupID}); err != nil {
		return err
	} else if len(groupsList) != 1 {
		return fmt.Errorf("group not find: %d", groupID)
	} else {
		group = groupsList[0]
	}

	// collections
	collections, err := sf.DBM.GetCollectionsByGroupIDs([]int64{groupID})
	if err != nil {
		return err
	}

	collectionsIDs := collections.GetIDs()
	collections2cnt, err := sf.DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		return err
	}

	var medias dto.MediaList
	if mediasIDs := collections.GetMediasIDs(); len(mediasIDs) > 0 {
		medias, err = sf.DBM.GetMediasByIDs(mediasIDs)
		if err != nil {
			return err
		}
	}

	// collections
	collectionModelList := models.NewModelCollectionList(collections)
	for _, collectionModel := range collectionModelList.List {
		collectionModel.WithObjectsCnt(collections2cnt[*collectionModel.ID])
	}

	se.Group = &models.FullGroup{
		Collections: collectionModelList,
		Group:       models.NewModelGroup(group),
	}
	se.Medias = models.NewModelMediaList(medias)

	return nil
}

func (sf *PayloadedEventFiller) fillActor(se *models.PayloadedEvent, actorID int64) error {
	var actor *dto.Actor
	if actorsList, err := sf.DBM.GetActorsByIDs([]int64{actorID}); err != nil {
		return err
	} else if len(actorsList) != 1 {
		err := fmt.Errorf("cant find actor with id: %d", actorID)
		return err
	} else {
		actor = actorsList[0]
	}

	se.Actor = models.NewModelActor(actor)

	return nil
}

func (sf *PayloadedEventFiller) fillBadge(se *models.PayloadedEvent, badgeID int64) error {
	var badge *dto.Badge
	if badgesList, err := sf.DBM.GetBadgesByIDs([]int64{badgeID}); err != nil {
		return err
	} else if len(badgesList) != 1 {
		err := fmt.Errorf("cant find badge with id: %d", badgeID)
		return err
	} else {
		badge = badgesList[0]
	}

	se.Badge = models.NewModelBadge(badge)

	return nil
}

func (sf *PayloadedEventFiller) fillMaterial(se *models.PayloadedEvent, materialID int64) error {
	var material *dto.Material
	if materialsList, err := sf.DBM.GetMaterialsByIDs([]int64{materialID}); err != nil {
		return err
	} else if len(materialsList) != 1 {
		err := fmt.Errorf("cant find material with id: %d", materialID)
		return err
	} else {
		material = materialsList[0]
	}

	se.Material = models.NewModelMaterial(material)

	return nil
}

func (sf *PayloadedEventFiller) fillOriginLocation(se *models.PayloadedEvent, originLocationID int64) error {
	var originLocation *dto.OriginLocation
	if originLocationsList, err := sf.DBM.GetOriginLocationsByIDs([]int64{originLocationID}); err != nil {
		return err
	} else if len(originLocationsList) != 1 {
		err := fmt.Errorf("cant find origin Location with id: %d", originLocationID)
		return err
	} else {
		originLocation = originLocationsList[0]
	}

	se.OriginLocation = models.NewModelOriginLocation(originLocation)

	return nil
}

func (sf *PayloadedEventFiller) fillNamedDateInterval(se *models.PayloadedEvent, namedDateIntervalID int64) error {
	var namedDateInterval *dto.NamedDateInterval
	if namedDateIntervalsList, err := sf.DBM.GetNamedDateIntervalsByIDs([]int64{namedDateIntervalID}); err != nil {
		return err
	} else if len(namedDateIntervalsList) != 1 {
		err := fmt.Errorf("cant find named Date Interval with id: %d", namedDateIntervalID)
		return err
	} else {
		namedDateInterval = namedDateIntervalsList[0]
	}

	se.NamedDateInterval = models.NewModelDateInterval(namedDateInterval)

	return nil
}
