package dto

import (
	"fmt"
)

// SearchObject TBD
type SearchObject struct {
	ID              int64   `json:"id,omitempty"`
	CollectionID    int64   `json:"collection_id,omitempty"`
	RootID          int64   `json:"root_id,omitempty"`
	Name            string  `json:"name,omitempty"`
	Description     string  `json:"description,omitempty"`
	Provenance      string  `json:"provenance,omitempty"`
	RootIDNumber    string  `json:"root_id_number,omitempty"`
	Materials       []int64 `json:"materials,omitempty"`
	Actors          []int64 `json:"actors,omitempty"`
	FirstActorID    int64   `json:"first_actor_id,omitempty"`
	FirstActorName  string  `json:"first_actor_name,omitempty"`
	Badges          []int64 `json:"badges,omitempty"`
	OriginLocations []int64 `json:"origin_locations,omitempty"`
	Statuses        []int64 `json:"statuses,omitempty"`

	ProductionDateIntervalID   *int64 `json:"production_date_interval_id,omitempty"`
	ProductionDateIntervalFrom *int64 `json:"production_date_interval_from,omitempty"`
	ProductionDateIntervalTo   *int64 `json:"production_date_interval_to,omitempty"`

	UpdateTime string `json:"update_time,omitempty"`

	forceUpdate     bool
	ExternalVersion int64
}

// SearchObjectForceRepresentation TBD
type SearchObjectForceRepresentation struct {
	ID              int64   `json:"id"`
	CollectionID    int64   `json:"collection_id"`
	RootID          int64   `json:"root_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Provenance      string  `json:"provenance"`
	RootIDNumber    string  `json:"root_id_number"`
	Materials       []int64 `json:"materials"`
	Actors          []int64 `json:"actors"`
	FirstActorID    int64   `json:"first_actor_id"`
	FirstActorName  string  `json:"first_actor_name"`
	Badges          []int64 `json:"badges"`
	OriginLocations []int64 `json:"origin_locations"`
	Statuses        []int64 `json:"statuses"`

	ProductionDateIntervalID   *int64 `json:"production_date_interval_id"`
	ProductionDateIntervalFrom *int64 `json:"production_date_interval_from"`
	ProductionDateIntervalTo   *int64 `json:"production_date_interval_to"`

	UpdateTime string `json:"update_time"`
}

// ToForceRepresentation TBD
func (s *SearchObject) ToForceRepresentation() *SearchObjectForceRepresentation {
	return &SearchObjectForceRepresentation{
		ID:              s.ID,
		CollectionID:    s.CollectionID,
		RootID:          s.RootID,
		Name:            s.Name,
		Description:     s.Description,
		Provenance:      s.Provenance,
		RootIDNumber:    s.RootIDNumber,
		Materials:       s.Materials,
		Actors:          s.Actors,
		FirstActorID:    s.FirstActorID,
		FirstActorName:  s.FirstActorName,
		Badges:          s.Badges,
		OriginLocations: s.OriginLocations,
		Statuses:        s.Statuses,

		ProductionDateIntervalID:   s.ProductionDateIntervalID,
		ProductionDateIntervalFrom: s.ProductionDateIntervalFrom,
		ProductionDateIntervalTo:   s.ProductionDateIntervalTo,

		UpdateTime: s.UpdateTime,
	}
}

// SearchObjects TBD
func (oo ObjectList) SearchObjects() SearchObjectList {
	rl := make(SearchObjectList, len(oo))
	for i := range oo {
		rl[i] = oo[i].SearchObject()
	}
	return rl
}

// SearchObject TBD
func (o *Object) SearchObject() *SearchObject {
	return &SearchObject{
		ID:                         o.ID,
		CollectionID:               o.CollectionID,
		Name:                       o.Name,
		Description:                o.Description,
		Provenance:                 o.Provenance,
		RootIDNumber:               o.RootIDNumber,
		ProductionDateIntervalID:   o.ProductionDateIntervalID,
		ProductionDateIntervalFrom: o.ProductionDateIntervalFrom,
		ProductionDateIntervalTo:   o.ProductionDateIntervalTo,
		// yyyyMMdd'T'HHmmssZ
		UpdateTime: o.UpdateTime.Format("20060102T150405-0700"),
	}
}

// ObjectSearchFilters TBD
type ObjectSearchFilters struct {
	CollectionsGroups [][]int64
	Collections       []int64
	Actors            []int64
	Badges            []int64
	Materials         []int64
	OriginLocations   []int64
	Statuses          []int64

	CollectionsToExclude []int64

	ProductionDateIntervalID   *int64
	ProductionDateIntervalFrom *int64
	ProductionDateIntervalTo   *int64
}

// ObjectSearchFiltersResults TBD
type ObjectSearchFiltersResults struct {
	Actors                   ObjectSearchFiltersResultList
	Badges                   ObjectSearchFiltersResultList
	Materials                ObjectSearchFiltersResultList
	OriginLocations          ObjectSearchFiltersResultList
	Statuses                 ObjectSearchFiltersResultList
	ProdutcionNamedIntervals ObjectSearchFiltersResultList
	Collections              ObjectSearchFiltersResultList
}

// IDs TBD
func (sr ObjectSearchFiltersResultList) IDs() []int64 {
	res := make([]int64, len(sr))
	for i := range sr {
		res[i] = sr[i].PropertyID
	}
	return res
}

// ObjectSearchFiltersResult TBD
type ObjectSearchFiltersResult struct {
	PropertyID int64
	Cnt        int64
}

// ObjectSearchFiltersResultList TBD
type ObjectSearchFiltersResultList []*ObjectSearchFiltersResult

// NewDefaultObjectOrders TBD
func NewDefaultObjectOrders() *ObjectOrders {
	return &ObjectOrders{
		CreationTime: -1,
	}
}

// ObjectOrders TBD
type ObjectOrders struct {
	CreationTime int16
	UpdateTime   int16
	ActorName    int16
	Name         int16
}

// ForceUpdate TBD
func (s *SearchObject) ForceUpdate(b bool) *SearchObject {
	s.forceUpdate = b
	return s
}

// SetExternalVersion TBD
func (s *SearchObject) SetExternalVersion(version int64) *SearchObject {
	s.ExternalVersion = version
	return s
}

// IDKey TBD
func (s *SearchObject) IDKey() string {
	return fmt.Sprintf("%d", s.ID)
}

// Doc TBD
func (s *SearchObject) Doc() interface{} {
	if s.forceUpdate {
		return s.ToForceRepresentation()
	}
	return s
}

// Type TBD
func (s *SearchObject) Type() string {
	return "object"
}

// Version TBD
func (s *SearchObject) Version() int64 {
	return s.ExternalVersion
}

// WithRootID TBD
func (s *SearchObject) WithRootID(rootID int64) *SearchObject {
	s.RootID = rootID
	return s
}

// WithMaterialsIDs TBD
func (s *SearchObject) WithMaterialsIDs(m []int64) *SearchObject {
	s.Materials = m
	return s
}

// WithActorsIDs TBD
func (s *SearchObject) WithActorsIDs(a []int64) *SearchObject {
	if len(a) != 0 {
		s.FirstActorID = a[0]
	}
	s.Actors = a
	return s
}

// WithActorsNames TBD
func (s *SearchObject) WithActorsNames(n []string) *SearchObject {
	if len(n) < 1 {
		return s
	}
	s.FirstActorName = n[0]
	return s
}

// WithBadgesIDs TBD
func (s *SearchObject) WithBadgesIDs(b []int64) *SearchObject {
	s.Badges = b
	return s
}

// WithOriginLocationsIDs TBD
func (s *SearchObject) WithOriginLocationsIDs(ol []int64) *SearchObject {
	s.OriginLocations = ol
	return s
}

// WithStatusesIDs TBD
func (s *SearchObject) WithStatusesIDs(st []int64) *SearchObject {
	s.Statuses = st
	return s
}

// WithProductionIntervalFrom TBD
func (s *SearchObject) WithProductionIntervalFrom(gdate *int64) *SearchObject {
	s.ProductionDateIntervalFrom = gdate
	return s
}

// WithProductionIntervalTo TBD
func (s *SearchObject) WithProductionIntervalTo(gdate *int64) *SearchObject {
	s.ProductionDateIntervalTo = gdate
	return s
}

// SearchObjectList TBD
type SearchObjectList []*SearchObject

// IDs TBD
func (s SearchObjectList) IDs() []int64 {
	ids := make([]int64, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ids
}
