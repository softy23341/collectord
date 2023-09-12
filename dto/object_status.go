package dto

//go:generate dbgen -type ObjectStatus

// ObjectStatus TBD
type ObjectStatus struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Description  string `db:"description"`
	ImageMediaID *int64 `db:"image_media_id"`
}

// ObjectStatusList TBD
type ObjectStatusList []*ObjectStatus

// GetImageMediaIDs TBD
func (c ObjectStatusList) GetImageMediaIDs() []int64 {
	var ids []int64
	for _, status := range c {
		if id := status.ImageMediaID; id != nil {
			ids = append(ids, *id)
		}
	}
	return ids
}
