package dto

//go:generate dbgen -type OriginLocation

// OriginLocation TBD
type OriginLocation struct {
	ID         int64  `db:"id"`
	RootID     *int64 `db:"root_id"`
	Name       string `db:"name"`
	NormalName string `db:"normal_name"`
}

// OriginLocationList TBD
type OriginLocationList []*OriginLocation

// IDs ids slice
func (a OriginLocationList) IDs() []int64 {
	o := make([]int64, len(a))
	for i := range a {
		o[i] = a[i].ID
	}
	return o
}
