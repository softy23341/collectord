package resource

import (
	cleaver "git.softndit.com/collector/backend/cleaver/client"
	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	npusherclient "git.softndit.com/collector/backend/npusher/client"
	"git.softndit.com/collector/backend/services"
	"gopkg.in/inconshreveable/log15.v2"
)

// Context TBD
type Context struct {
	DBM           dal.TrManager
	Log           log15.Logger
	CleaverClient cleaver.ConnectClient
	SearchClient  services.SearchClient
	FileStorage   services.FileStorage

	EventSender     *services.EventSender
	MessengerClient *services.MessengerClient
	PushClient      npusherclient.Client
	MailClient      services.MailSender

	JobPool *delayedjob.Pool
}

// ReindexContext TBD
type ReindexContext struct {
	DBM dal.TrManager
}

// PrepareToReindex TBD
func PrepareToReindex(c *ReindexContext, objectsList dto.ObjectList) (services.ObjectDocsForIndex, error) {
	collectionsListIDs := objectsList.GetCollectionsIDs()
	objectsIDs := objectsList.GetIDs()

	reindexExternalVersion, err := c.DBM.GetCurrentTransaction()
	if err != nil {
		return nil, err
	}

	// collections
	collectionsList, err := c.DBM.GetCollectionsByIDs(collectionsListIDs)
	if err != nil {
		return nil, err
	}
	cid2collection := collectionsList.IDToCollection()

	// materials
	materialRefList, err := c.DBM.GetMaterialRefsByObjectsIDs(objectsIDs)
	if err != nil {
		return nil, err
	}
	objectID2MaterialsIDs := materialRefList.ObjectIDToMaterialListIDsMap()

	// actors refs
	actorsRefs, err := c.DBM.GetObjectsActorRefs(objectsIDs)
	if err != nil {
		return nil, err
	}
	objectID2ActorsIDs := actorsRefs.ObjectIDToActorsIDs()

	// actors
	actors, err := c.DBM.GetActorsByIDs(actorsRefs.ActorsIDs())
	if err != nil {
		return nil, err
	}
	ID2Actor := actors.IDToActor()

	// badges
	badgesRefs, err := c.DBM.GetObjectsBadgeRefs(objectsIDs)
	if err != nil {
		return nil, nil
	}
	objectID2BadgesIDs := badgesRefs.ObjectIDToBadgesIDs()

	// originLocations
	originLocationsRefs, err := c.DBM.GetObjectsOriginLocationRefs(objectsIDs)
	if err != nil {
		return nil, err
	}
	objectID2OriginLocationIDs := originLocationsRefs.ObjectIDToOriginLocationsIDs()

	// statuses
	statusesRefs, err := c.DBM.GetCurrentObjectsStatusesRefs(objectsIDs)
	if err != nil {
		return nil, err
	}
	objectID2StatusesIDs := statusesRefs.ObjectToOneStatusMap()

	// production date interval
	productionIntervalsIDs := objectsList.ProductionDateIntervalIDs()
	intervalsList, err := c.DBM.GetNamedDateIntervalsByIDs(productionIntervalsIDs)
	if err != nil {
		return nil, err
	}
	intervalIDToInterval := intervalsList.IDToNamedDateInterval()

	searchObjects := make(services.ObjectDocsForIndex, 0, len(objectsList))
	for _, object := range objectsList {

		objectActorsIDs := objectID2ActorsIDs[object.ID]
		forIndex := object.SearchObject().ForceUpdate(true).
			WithRootID(cid2collection[object.CollectionID].RootID).
			WithMaterialsIDs(objectID2MaterialsIDs[object.ID]).
			WithActorsIDs(objectActorsIDs).
			WithBadgesIDs(objectID2BadgesIDs[object.ID]).
			WithOriginLocationsIDs(objectID2OriginLocationIDs[object.ID])

		if status := objectID2StatusesIDs[object.ID]; status != nil {
			forIndex.WithStatusesIDs([]int64{status.ObjectStatusID})
		}

		if len(objectActorsIDs) > 0 {
			actorNames := make([]string, 0, len(objectActorsIDs))
			for _, id := range objectActorsIDs {
				actorNames = append(actorNames, ID2Actor[id].NormalName)
			}
			forIndex.WithActorsNames(actorNames)
		}

		if productionDateID := object.ProductionDateIntervalID; productionDateID != nil {
			interval := intervalIDToInterval[*productionDateID]

			forIndex.
				WithProductionIntervalFrom(&interval.ProductionDateIntervalFrom).
				WithProductionIntervalTo(&interval.ProductionDateIntervalTo)
		}

		forIndex.SetExternalVersion(reindexExternalVersion)

		searchObjects = append(searchObjects, forIndex)
	}

	return searchObjects, nil
}

// ReindexObjects TBD
func (c *Context) ReindexObjects(objects dto.SearchObjectList) error {
	objectsList, err := c.DBM.GetObjectsByIDs(objects.IDs())
	if err != nil {
		return err
	}

	searchObjects, err := PrepareToReindex(&ReindexContext{DBM: c.DBM}, objectsList)
	if err != nil {
		return err
	}

	reindexExternalVersion, err := c.DBM.GetCurrentTransaction()
	if err != nil {
		return err
	}

	objectsToDelete := make(services.ObjectDocsForIndex, 0)
	if len(objects) > len(objectsList) {
		objectsMap := objectsList.ToMap()
		for _, foundObject := range objects {
			if _, foundInDB := objectsMap[foundObject.ID]; !foundInDB {
				foundObject.SetExternalVersion(reindexExternalVersion)
				objectsToDelete = append(objectsToDelete, foundObject)
			}
		}
	}

	if len(objectsToDelete) > 0 {
		if err := c.SearchClient.BulkObjectDelete(objectsToDelete); err != nil {
			return err
		}
	}

	if len(searchObjects) > 0 {
		return c.SearchClient.BulkObjectIndex(searchObjects)
	}
	return nil
}
