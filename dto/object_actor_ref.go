package dto

import "sort"

//go:generate dbgen -type ObjectActorRef

// ObjectActorRef TBD
type ObjectActorRef struct {
	ID       int64 `db:"id"`
	ObjectID int64 `db:"object_id"`
	ActorID  int64 `db:"actor_id"`
}

// ObjectActorRefList TBD
type ObjectActorRefList []*ObjectActorRef

func (o ObjectActorRefList) Len() int           { return len(o) }
func (o ObjectActorRefList) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o ObjectActorRefList) Less(i, j int) bool { return o[i].ID < o[j].ID }

// ObjectIDToActorsIDs TBD
func (o ObjectActorRefList) ObjectIDToActorsIDs() map[int64][]int64 {
	o2r := make(map[int64]ObjectActorRefList, len(o))
	for _, ref := range o {
		o2r[ref.ObjectID] = append(o2r[ref.ObjectID], ref)
	}
	o2a := make(map[int64][]int64, len(o))
	for objectID, refs := range o2r {
		sort.Sort(refs)
		for _, ref := range refs {
			o2a[objectID] = append(o2a[objectID], ref.ActorID)
		}
	}
	return o2a
}

// ActorsIDs TBD
func (o ObjectActorRefList) ActorsIDs() []int64 {
	actorsIDsMap := make(map[int64]struct{}, len(o))
	for _, actorRef := range o {
		actorsIDsMap[actorRef.ActorID] = struct{}{}
	}

	var actorsIDsList []int64
	for actorID := range actorsIDsMap {
		actorsIDsList = append(actorsIDsList, actorID)
	}

	return actorsIDsList
}
