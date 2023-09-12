package dto

//go:generate dbgen -type Badge

// Badge TBD
type Badge struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	NormalName string `db:"normal_name"`
	Color      string `db:"color"`
	RootID     *int64 `db:"root_id"`
}

// BadgeList TBD
type BadgeList []*Badge

// IDs is id slice
func (bl BadgeList) IDs() []int64 {
	ids := make([]int64, len(bl))
	for i := range bl {
		ids[i] = bl[i].ID
	}
	return ids
}

// NormalNames is name slice
func (bl BadgeList) NormalNames() []string {
	ids := make([]string, len(bl))
	for i := range bl {
		ids[i] = bl[i].NormalName
	}
	return ids
}

// InputBadge TBD
type InputBadge struct {
	Name  string `db:"name"`
	Color string `db:"color"`
}

// InputBadgeList TBD
type InputBadgeList []*InputBadge
