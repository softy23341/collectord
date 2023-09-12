package dto

import "git.softndit.com/collector/backend/util"

//go:generate dbgen -type Group

// Group TBD
type Group struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	RootID     int64  `db:"root_id"`
	UserID     *int64 `db:"user_id"`
	UserUniqID *int64 `db:"user_uniq_id"`
}

// GroupList TBD
type GroupList []*Group

// GetIDs TBD
func (grl GroupList) GetIDs() []int64 {
	o := make([]int64, len(grl))
	for i := range grl {
		o[i] = grl[i].ID
	}
	return o
}

// RemoveByIDs TBD
func (grl GroupList) RemoveByIDs(ids []int64) GroupList {
	filteredGroupList := make(GroupList, 0, len(grl))
	toExclude := util.Int64SliceToSet(ids)
	for _, group := range grl {
		if toExclude.Contains(group.ID) {
			continue
		}

		filteredGroupList = append(filteredGroupList, group)
	}
	return filteredGroupList
}

// RootsIDs TBD
func (grl GroupList) RootsIDs() []int64 {
	rootsIDsMap := make(map[int64]struct{}, len(grl))
	for _, collection := range grl {
		rootsIDsMap[collection.RootID] = struct{}{}
	}

	var rootsIDs []int64
	for rootID := range rootsIDsMap {
		rootsIDs = append(rootsIDs, rootID)
	}

	return rootsIDs
}
