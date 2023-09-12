package dto

//go:generate dbgen -type ObjectOriginLocationRef

// ObjectOriginLocationRef TBD
type ObjectOriginLocationRef struct {
	ObjectID         int64 `db:"object_id"`
	OriginLocationID int64 `db:"origin_location_id"`
}

// ObjectOriginLocationRefList TBD
type ObjectOriginLocationRefList []*ObjectOriginLocationRef

// ObjectIDToOriginLocationsIDs TBD
func (o ObjectOriginLocationRefList) ObjectIDToOriginLocationsIDs() map[int64][]int64 {
	o2a := make(map[int64][]int64, len(o))
	for _, ref := range o {
		o2a[ref.ObjectID] = append(o2a[ref.ObjectID], ref.OriginLocationID)
	}
	return o2a
}

// OriginLocationsIDs TBD
func (o ObjectOriginLocationRefList) OriginLocationsIDs() []int64 {
	originLocationsIDsMap := make(map[int64]struct{}, len(o))
	for _, originLocationRef := range o {
		originLocationsIDsMap[originLocationRef.OriginLocationID] = struct{}{}
	}

	var originLocationsIDsList []int64
	for originLocationID := range originLocationsIDsMap {
		originLocationsIDsList = append(originLocationsIDsList, originLocationID)
	}

	return originLocationsIDsList
}
