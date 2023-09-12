package dto

//go:generate dbgen -type CollectionGroupRef

// CollectionGroupRef TBD
type CollectionGroupRef struct {
	CollectionID int64 `db:"collection_id"`
	GroupID      int64 `db:"group_id"`
}

// CollectionGroupRefList TBD
type CollectionGroupRefList []*CollectionGroupRef

// CollectionIDToGroupsIDs TBD
func (o CollectionGroupRefList) CollectionIDToGroupsIDs() map[int64][]int64 {
	o2a := make(map[int64][]int64, len(o))
	for _, ref := range o {
		o2a[ref.CollectionID] = append(o2a[ref.CollectionID], ref.GroupID)
	}
	return o2a
}

// GroupIDToCollectionsIDs TBD
func (o CollectionGroupRefList) GroupIDToCollectionsIDs() map[int64][]int64 {
	o2a := make(map[int64][]int64, len(o))
	for _, ref := range o {
		o2a[ref.GroupID] = append(o2a[ref.GroupID], ref.CollectionID)
	}
	return o2a
}

// GroupsIDs TBD
func (o CollectionGroupRefList) GroupsIDs() []int64 {
	groupsIDsMap := make(map[int64]struct{}, len(o))
	for _, groupRef := range o {
		groupsIDsMap[groupRef.GroupID] = struct{}{}
	}

	var groupsIDsList []int64
	for groupID := range groupsIDsMap {
		groupsIDsList = append(groupsIDsList, groupID)
	}

	return groupsIDsList
}
