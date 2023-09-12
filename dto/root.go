package dto

//go:generate dbgen -type Root

type (
	// Root TBD
	Root struct {
		ID int64 `db:"id"`
	}

	// RootList TBD
	RootList []*Root
)

// IDs TBD
func (r RootList) IDs() []int64 {
	m := make(map[int64]struct{})
	for _, root := range r {
		m[root.ID] = struct{}{}
	}

	ids := make([]int64, 0, len(m))
	for k := range m {
		ids = append(ids, k)
	}
	return ids
}
