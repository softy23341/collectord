package dto

import (
	"time"

	"git.softndit.com/collector/backend/util"
)

//go:generate dbgen -type Collection

// collection types
const (
	RegularCollectionTypo CollectionTypo = 0
	DraftCollectionTypo                  = 10
	TrashCollectionTypo                  = 20
)

// CollectionTypo - collection type
type CollectionTypo int16

// CollectionTypo - collection type
type HiddenField string

// MediaTypeSet TBD
type HiddenFieldSet []HiddenField

// NewCollection TBD
func NewCollection() *Collection {
	return &Collection{
		CreationTime: time.Now(),
	}
}

// Collection TBD
type Collection struct {
	ID           int64          `db:"id"`
	Name         string         `db:"name"`
	RootID       int64          `db:"root_id"`
	UserID       *int64         `db:"user_id"`
	UserUniqID   *int64         `db:"user_uniq_id"`
	CreationTime time.Time      `db:"creation_time"`
	Typo         CollectionTypo `db:"typo"`
	Description  string         `db:"description"`
	ImageMediaID *int64         `db:"image_media_id"`
	Public       bool           `db:"public"`
	IsAnonymous  bool           `db:"is_anonymous"`
}

// IsRegular TBD
func (c *Collection) IsRegular() bool {
	return c.Typo == RegularCollectionTypo
}

// IsTrash TBD
func (c *Collection) IsTrash() bool {
	return c.Typo == TrashCollectionTypo
}

// IsDraft TBD
func (c *Collection) IsDraft() bool {
	return c.Typo == DraftCollectionTypo
}

// CollectionList TBD
type CollectionList []*Collection

// RemoveByIDs TBD
func (c CollectionList) RemoveByIDs(ids []int64) CollectionList {
	filteredCollectionList := make(CollectionList, 0, len(c))
	toExclude := util.Int64SliceToSet(ids)
	for _, collection := range c {
		if toExclude.Contains(collection.ID) {
			continue
		}

		filteredCollectionList = append(filteredCollectionList, collection)
	}
	return filteredCollectionList
}

// IDToCollection TBD
func (c CollectionList) IDToCollection() map[int64]*Collection {
	id2col := make(map[int64]*Collection, 0)
	for _, collection := range c {
		id2col[collection.ID] = collection
	}
	return id2col
}

// GetMediasIDs TBD
func (c CollectionList) GetMediasIDs() []int64 {
	var ids []int64
	for _, collection := range c {
		if id := collection.ImageMediaID; id != nil {
			ids = append(ids, *id)
		}
	}
	return ids
}

// GetIDs TBD
func (c CollectionList) GetIDs() []int64 {
	ids := make([]int64, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return ids
}

// RootsIDs TBD
func (c CollectionList) RootsIDs() []int64 {
	rootsIDsMap := make(map[int64]struct{}, len(c))
	for _, collection := range c {
		rootsIDsMap[collection.RootID] = struct{}{}
	}

	rootsIDs := make([]int64, 0, len(rootsIDsMap))
	for rootID := range rootsIDsMap {
		rootsIDs = append(rootsIDs, rootID)
	}

	return rootsIDs
}

// NewRegularCollection TBD
func NewRegularCollection(rootID int64) *Collection {
	return &Collection{
		RootID: rootID,
		Typo:   RegularCollectionTypo,
	}
}

// NewDraftCollection TBD
func NewDraftCollection(rootID int64) *Collection {
	return &Collection{
		Name:   "Draft",
		RootID: rootID,
		Typo:   DraftCollectionTypo,
	}
}

// NewTrashCollection TBD
func NewTrashCollection(rootID int64) *Collection {
	return &Collection{
		Name:   "Trash",
		RootID: rootID,
		Typo:   TrashCollectionTypo,
	}
}
