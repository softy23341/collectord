package dto

import "sort"

//go:generate dbgen -type ObjectMediaRef

// ObjectMediaRef TBD
type ObjectMediaRef struct {
	ObjectID      int64 `db:"object_id"`
	MediaID       int64 `db:"media_id"`
	MediaPosition int32 `db:"media_position"`
}

// ObjectMediaRefList TBD
type ObjectMediaRefList []*ObjectMediaRef

func (o ObjectMediaRefList) Len() int      { return len(o) }
func (o ObjectMediaRefList) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o ObjectMediaRefList) Less(i, j int) bool {
	// TODO make it cool
	if o[i].MediaPosition != o[j].MediaPosition {
		return o[i].MediaPosition < o[j].MediaPosition
	}

	return o[i].ObjectID < o[j].ObjectID
}

// ObjectIDToMediasIDs TBD
func (o ObjectMediaRefList) ObjectIDToMediasIDs() map[int64][]int64 {
	o2m := make(map[int64]ObjectMediaRefList, len(o))
	for _, ref := range o {
		o2m[ref.ObjectID] = append(o2m[ref.ObjectID], ref)
	}

	o2ids := make(map[int64][]int64, len(o))
	for objectID, refList := range o2m {
		sort.Sort(refList)
		o2ids[objectID] = refList.OrderedMediasIDs()
	}

	return o2ids
}

// UniqMediasIDs TBD
func (o ObjectMediaRefList) UniqMediasIDs() []int64 {
	mediasIDsMap := make(map[int64]struct{}, len(o))
	for _, mediaRef := range o {
		mediasIDsMap[mediaRef.MediaID] = struct{}{}
	}

	var mediasIDsList []int64
	for mediaID := range mediasIDsMap {
		mediasIDsList = append(mediasIDsList, mediaID)
	}

	return mediasIDsList
}

// OrderedMediasIDs TBD
func (o ObjectMediaRefList) OrderedMediasIDs() []int64 {
	mediasIDsList := make([]int64, 0, len(o))
	for _, mediaRef := range o {
		mediasIDsList = append(mediasIDsList, mediaRef.MediaID)
	}

	return mediasIDsList
}
