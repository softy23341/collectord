package services

import "git.softndit.com/collector/backend/dto"

// SearchClient TBD
type SearchClient interface {
	IndexObject(ObjectDocForIndex) error
	BulkObjectDelete(ObjectDocsForIndex) error
	BulkObjectIndex(ObjectDocsForIndex) error
	SearchObjects(sq *SearchQuery) (*SearchObjectsResult, error)

	ScrollThrought(sq *ScrollSearchQuery, reindex func(objects dto.SearchObjectList) error) error
}

// ObjectDocForIndex TBD
type ObjectDocForIndex interface {
	IDKey() string
	Type() string
	Doc() interface{}
	Version() int64
}

// ObjectDocsForIndex TBD
type ObjectDocsForIndex []ObjectDocForIndex

// SearchQuery TBD
type SearchQuery struct {
	Query     string
	RootID    int64
	Filters   *dto.ObjectSearchFilters
	Orders    *dto.ObjectOrders
	Paginator *dto.PagePaginator
}

// ScrollSearchQuery TBD
type ScrollSearchQuery struct {
	ObjectIDs []int64
	RootID    *int64
	Filters   *dto.ObjectSearchFilters
}

// SearchObjectsResult TBD
type SearchObjectsResult struct {
	TotalHits int64
	Filters   *dto.ObjectSearchFiltersResults
	Objects   dto.SearchObjectList
}
