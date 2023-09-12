package dto

//go:generate dbgen -type ObjectBadgeRef

// ObjectBadgeRef TBD
type ObjectBadgeRef struct {
	ObjectID int64 `db:"object_id"`
	BadgeID  int64 `db:"badge_id"`
}

// ObjectBadgeRefList TBD
type ObjectBadgeRefList []*ObjectBadgeRef

// ObjectIDToBadgesIDs TBD
func (o ObjectBadgeRefList) ObjectIDToBadgesIDs() map[int64][]int64 {
	o2a := make(map[int64][]int64, len(o))
	for _, ref := range o {
		o2a[ref.ObjectID] = append(o2a[ref.ObjectID], ref.BadgeID)
	}
	return o2a
}

// BadgesIDs TBD
func (o ObjectBadgeRefList) BadgesIDs() []int64 {
	badgesIDsMap := make(map[int64]struct{}, len(o))
	for _, badgeRef := range o {
		badgesIDsMap[badgeRef.BadgeID] = struct{}{}
	}

	var badgesIDsList []int64
	for badgeID := range badgesIDsMap {
		badgesIDsList = append(badgesIDsList, badgeID)
	}

	return badgesIDsList
}
