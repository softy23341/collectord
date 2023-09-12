package dto

//go:generate dbgen -type Material

// Material TBD
type Material struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	NormalName string `db:"normal_name"`
	RootID     *int64 `db:"root_id"`
}

// MaterialList TBD
type MaterialList []*Material

// IDs is id slice
func (ml MaterialList) IDs() []int64 {
	ids := make([]int64, len(ml))
	for i := range ml {
		ids[i] = ml[i].ID
	}
	return ids
}

// NormalNames is name slice
func (ml MaterialList) NormalNames() []string {
	ids := make([]string, len(ml))
	for i := range ml {
		ids[i] = ml[i].NormalName
	}
	return ids
}
