package dto

import "time"

//go:generate dbgen -type Object
//go:generate dbgen -type ObjectPreview

// ObjectEntity TBD
type ObjectEntity interface {
	GetID() int64
}

// ObjectEntityList TBD
type ObjectEntityList interface {
	GetIDs() []int64
}

// NewObject TBD
func NewObject() *Object {
	return &Object{
		CreationTime: time.Now(),
		UpdateTime:   time.Now(),
	}
}

// Object TBD
type Object struct {
	ID                         int64      `db:"id"`
	CollectionID               int64      `db:"collection_id"`
	UserID                     int64      `db:"user_id"`
	UserUniqID                 int64      `db:"user_uniq_id"`
	Name                       string     `db:"name"`
	ProductionDateIntervalID   *int64     `db:"production_date_interval_id"`
	ProductionDateIntervalFrom *int64     `db:"production_date_interval_from"`
	ProductionDateIntervalTo   *int64     `db:"production_date_interval_to"`
	Description                string     `db:"description"`
	Provenance                 string     `db:"provenance"`
	PurchaseDate               *time.Time `db:"purchase_date"`
	PurchasePrice              *int64     `db:"purchase_price"`
	PurchaseCurrencyID         *int64     `db:"purchase_price_currency_id"`
	RootIDNumber               string     `db:"root_id_number"`
	CreationTime               time.Time  `db:"creation_time"`
	UpdateTime                 time.Time  `db:"update_time"`
}

// IsProductionDatePresent TBD
func (o *Object) IsProductionDatePresent() bool {
	return o.ProductionDateIntervalID != nil || o.IsProductionDateIntervalPresent()
}

// IsProductionDateIntervalPresent TBD
func (o *Object) IsProductionDateIntervalPresent() bool {
	return o.ProductionDateIntervalFrom != nil && o.ProductionDateIntervalTo != nil
}

// GetID TBD
func (o *Object) GetID() int64 {
	return o.ID
}

// Preview TBD
func (o *Object) Preview() *ObjectPreview {
	return &ObjectPreview{
		ID:           o.ID,
		CollectionID: o.CollectionID,
		Name:         o.Name,
	}
}

var _ ObjectEntity = (*Object)(nil)

// ObjectList TBD
type ObjectList []*Object

// Preview TBD
func (os ObjectList) Preview() ObjectPreviewList {
	pl := make(ObjectPreviewList, 0, len(os))
	for _, o := range os {
		pl = append(pl, o.Preview())
	}

	return pl
}

// GetIDs TBD
func (os ObjectList) GetIDs() []int64 {
	ids := make([]int64, 0, len(os))
	for _, o := range os {
		ids = append(ids, o.GetID())
	}
	return ids
}

var _ ObjectEntityList = (*ObjectList)(nil)

// ObjectPreview TBD
type ObjectPreview struct {
	ID           int64  `db:"id"`
	CollectionID int64  `db:"collection_id"`
	Name         string `db:"name"`
	RootID       int64  `db:"root_id"`
}

// GetID TBD
func (o *ObjectPreview) GetID() int64 {
	return o.ID
}

var _ ObjectEntity = (*ObjectPreview)(nil)

// ObjectPreviewList TBD
type ObjectPreviewList []*ObjectPreview

// GetIDs TBD
func (os ObjectPreviewList) GetIDs() []int64 {
	ids := make([]int64, len(os))
	for i := range os {
		ids[i] = os[i].GetID()
	}
	return ids
}

// GetCollectionsIDs TBD
func (os ObjectPreviewList) GetCollectionsIDs() []int64 {
	collectionsIDsMap := make(map[int64]struct{}, len(os))
	for _, collectionRef := range os {
		collectionsIDsMap[collectionRef.CollectionID] = struct{}{}
	}

	var collectionsIDsList []int64
	for collectionID := range collectionsIDsMap {
		collectionsIDsList = append(collectionsIDsList, collectionID)
	}

	return collectionsIDsList
}

// GetCollectionsIDs TBD
func (os ObjectList) GetCollectionsIDs() []int64 {
	collectionsIDsMap := make(map[int64]struct{}, len(os))
	for _, collectionRef := range os {
		collectionsIDsMap[collectionRef.CollectionID] = struct{}{}
	}

	var collectionsIDsList []int64
	for collectionID := range collectionsIDsMap {
		collectionsIDsList = append(collectionsIDsList, collectionID)
	}

	return collectionsIDsList
}

// ProductionDateIntervalIDs TBD
func (os ObjectList) ProductionDateIntervalIDs() []int64 {
	productionDateIntervalIDs := make(map[int64]struct{}, len(os))
	for _, object := range os {
		if intervalID := object.ProductionDateIntervalID; intervalID != nil {
			productionDateIntervalIDs[*intervalID] = struct{}{}
		}
	}

	var idsList []int64
	for intervalID := range productionDateIntervalIDs {
		idsList = append(idsList, intervalID)
	}

	return idsList
}

// ToMap TBD
func (os ObjectList) ToMap() map[int64]*Object {
	id2object := make(map[int64]*Object, 0)
	for _, object := range os {
		id2object[object.ID] = object
	}
	return id2object
}

var _ ObjectEntityList = (*ObjectPreviewList)(nil)
