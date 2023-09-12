package services

import (
	"encoding/json"
	"fmt"
	"io"

	"git.softndit.com/collector/backend/dto"
	"gopkg.in/olivere/elastic.v3"
)

const elasticSearchExternalType = "external"

// ElasticSearchClientContext TBD
type ElasticSearchClientContext struct {
	URL         string
	ObjectIndex string
}

// NewElasticSearchClient TBD
func NewElasticSearchClient(c *ElasticSearchClientContext) (*ElasticSearchClient, error) {
	cl, err := elastic.NewClient(
		elastic.SetURL(c.URL),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{
		cl:          cl,
		ObjectIndex: c.ObjectIndex,
	}, nil
}

// ElasticSearchClient TBD
type ElasticSearchClient struct {
	cl          *elastic.Client
	ObjectIndex string
}

// IndexObject TBD
func (e *ElasticSearchClient) IndexObject(o ObjectDocForIndex) error {
	insertCommand := e.cl.Index().
		Index(e.ObjectIndex).
		Type(o.Type()).
		Id(o.IDKey()).
		BodyJson(o.Doc()).
		Refresh(true)

	if version := o.Version(); version != 0 {
		insertCommand.
			VersionType(elasticSearchExternalType).
			Version(version)
	}

	_, err := insertCommand.Do()

	return err
}

// BulkObjectIndex TBD
func (e *ElasticSearchClient) BulkObjectIndex(objects ObjectDocsForIndex) error {
	bulkService := e.cl.Bulk()
	for _, object := range objects {
		indexRequest := elastic.NewBulkIndexRequest().
			Index(e.ObjectIndex).
			Type(object.Type()).
			Doc(object.Doc()).
			Id(object.IDKey())

		if version := object.Version(); version != 0 {
			indexRequest.
				VersionType(elasticSearchExternalType).
				Version(version)
		}

		bulkService.Add(indexRequest)
	}

	_, err := bulkService.Do()
	return err
}

// BulkObjectDelete TBD
func (e *ElasticSearchClient) BulkObjectDelete(objects ObjectDocsForIndex) error {
	bulkService := e.cl.Bulk()
	for _, object := range objects {
		deleteRequest := elastic.NewBulkDeleteRequest().
			Index(e.ObjectIndex).
			Type(object.Type()).
			Id(object.IDKey())

		if version := object.Version(); version != 0 {
			deleteRequest.
				VersionType(elasticSearchExternalType).
				Version(version)
		}

		bulkService.Add(deleteRequest)
	}

	_, err := bulkService.Do()
	return err
}

// ScrollThrought TBD
func (e *ElasticSearchClient) ScrollThrought(sq *ScrollSearchQuery, reindex func(objects dto.SearchObjectList) error) error {
	scroll := e.cl.Scroll(e.ObjectIndex)
	{
		bq := elastic.NewBoolQuery()

		bq.Must(elastic.NewMatchAllQuery())
		if sq.RootID != nil {
			bq.Must(elastic.NewMatchQuery("root_id", sq.RootID))
		}

		if len(sq.ObjectIDs) != 0 {
			// TODO change to multiget
			values := make([]interface{}, len(sq.ObjectIDs))
			for i := range sq.ObjectIDs {
				values[i] = sq.ObjectIDs[i]
			}

			bq.Must(elastic.NewTermsQuery("id", values...))
		}

		if sq.Filters != nil { // filters
			for _, f := range sq.Filters.Collections {
				bq.Must(elastic.NewTermQuery("collection_id", f))
			}

			for _, collections := range sq.Filters.CollectionsGroups {
				values := make([]interface{}, len(collections))
				for i := range collections {
					values[i] = interface{}(collections[i])
				}
				bq.Must(elastic.NewTermsQuery("collection_id", values...))
			}

			for _, f := range sq.Filters.Actors {
				bq.Filter(elastic.NewTermQuery("actors", f))
			}

			for _, f := range sq.Filters.Badges {
				bq.Filter(elastic.NewTermQuery("badges", f))
			}

			for _, f := range sq.Filters.Materials {
				bq.Filter(elastic.NewTermQuery("materials", f))
			}

			for _, f := range sq.Filters.OriginLocations {
				bq.Filter(elastic.NewTermQuery("origin_locations", f))
			}

			for _, f := range sq.Filters.Statuses {
				bq.Filter(elastic.NewTermQuery("statuses", f))
			}

			// named date interval
			if id := sq.Filters.ProductionDateIntervalID; id != nil {
				bq.Must(elastic.NewTermQuery("production_date_interval_id", *id))
			}

		}

		scroll.Query(bq)
	}

	var (
		err          error
		searchResult *elastic.SearchResult
	)
	for searchResult, err = scroll.Do(); err == nil; searchResult, err = scroll.Do() {
		searchObjectsList := make(dto.SearchObjectList, len(searchResult.Hits.Hits))
		for i, h := range searchResult.Hits.Hits {
			object := &dto.SearchObject{}
			if err := json.Unmarshal(*h.Source, &object); err != nil {
				return err
			}
			searchObjectsList[i] = object
		}
		if err := reindex(searchObjectsList); err != nil {
			return err
		}
	}
	if err != io.EOF {
		return err
	}

	return nil
}

// SearchObjects TBD
// curl -XGET "$EL/object/_search?pretty" -d '
//
//	{
//		"from" : 0, "size" : 1,
//		"query": {
//			"bool": {
//				"must": {
//					"multi_match" : {
//						"query" : "Hope",
//						"fields" : ["name", "description"]
//					}
//				},
//				"filter": {
//					"term": { "materials": 11 }
//				},
//				"filter": {
//					"term": { "materials": 7 }
//				}
//			}
//		},
//		"aggs" : {
//			"actors" : {
//				"terms" : {
//					"field" : "actors"
//				}
//			},
//			"materials" : {
//				"terms" : {
//					"field" : "materials"
//				}
//			}
//		}
//	}
//
// '
func (e *ElasticSearchClient) SearchObjects(sq *SearchQuery) (*SearchObjectsResult, error) {
	search := e.cl.Search().
		Index(e.ObjectIndex).
		From(int(sq.Paginator.Offset())).Size(int(sq.Paginator.Limit()))

	{ // query
		bq := elastic.NewBoolQuery()

		{ // must
			if len(sq.Query) > 0 {
				mm := elastic.
					NewMultiMatchQuery(sq.Query, "name", "description", "provenance").
					Type("phrase_prefix")
				bq.Must(mm)
			} else {
				bq.Must(elastic.NewMatchAllQuery())
			}
			bq.Must(elastic.NewMatchQuery("root_id", sq.RootID))
		}

		if sq.Filters != nil { // filters
			for _, f := range sq.Filters.Collections {
				bq.Must(elastic.NewTermQuery("collection_id", f))
			}

			for _, collections := range sq.Filters.CollectionsGroups {
				values := make([]interface{}, len(collections))
				for i := range collections {
					values[i] = interface{}(collections[i])
				}
				bq.Must(elastic.NewTermsQuery("collection_id", values...))
			}

			for _, f := range sq.Filters.Actors {
				bq.Filter(elastic.NewTermQuery("actors", f))
			}

			for _, f := range sq.Filters.Badges {
				bq.Filter(elastic.NewTermQuery("badges", f))
			}

			for _, f := range sq.Filters.Materials {
				bq.Filter(elastic.NewTermQuery("materials", f))
			}

			for _, f := range sq.Filters.OriginLocations {
				bq.Filter(elastic.NewTermQuery("origin_locations", f))
			}

			for _, f := range sq.Filters.Statuses {
				bq.Filter(elastic.NewTermQuery("statuses", f))
			}

			// excludes
			for _, f := range sq.Filters.CollectionsToExclude {
				bq.MustNot(elastic.NewTermQuery("collection_id", f))
			}

			// named date interval
			if sq.Filters.ProductionDateIntervalTo != nil &&
				sq.Filters.ProductionDateIntervalFrom != nil {
				bq.Filter(elastic.NewRangeQuery("production_date_interval_to").
					Lte(sq.Filters.ProductionDateIntervalTo).
					Gte(sq.Filters.ProductionDateIntervalFrom))

				bq.Filter(elastic.NewRangeQuery("production_date_interval_from").
					Lte(sq.Filters.ProductionDateIntervalTo).
					Gte(sq.Filters.ProductionDateIntervalFrom))
			}
		}

		search.Query(bq)
	}

	{ // aggregation
		search.
			Aggregation("actors", elastic.NewTermsAggregation().Field("actors").Size(0)).
			Aggregation("badges", elastic.NewTermsAggregation().Field("badges").Size(0)).
			Aggregation("materials", elastic.NewTermsAggregation().Field("materials").Size(0)).
			Aggregation("origin_locations", elastic.NewTermsAggregation().Field("origin_locations").Size(0)).
			Aggregation("statuses", elastic.NewTermsAggregation().Field("statuses").Size(0)).
			Aggregation("production_date_interval_id", elastic.NewTermsAggregation().Field("production_date_interval_id").Size(0)).
			Aggregation("collection_id", elastic.NewTermsAggregation().Field("collection_id").Size(0))
	}

	// Sorting
	{
		if orders := sq.Orders; orders != nil {
			// creation time
			if t := orders.CreationTime; t != 0 {
				search.Sort("id", t > 0)
			}
			// actor name
			if t := orders.ActorName; t != 0 {
				search.SortBy(
					elastic.
						NewFieldSort("first_actor_name.raw").
						Order(t > 0).
						Missing("_last"),
				)
			}
			// object name
			if t := orders.Name; t != 0 {
				search.Sort("name.raw", t > 0)
			}

			// update time
			if t := orders.UpdateTime; t != 0 {
				search.Sort("update_time", t > 0)
			}
		} else {
			search.Sort("id", false)
		}
	}

	elasticSearchResult, err := search.Do()
	if err != nil {
		return nil, err
	}

	searchResult := &SearchObjectsResult{}
	searchResult.TotalHits = elasticSearchResult.Hits.TotalHits
	searchResult.Filters = &dto.ObjectSearchFiltersResults{}

	searchResult.Objects = make(dto.SearchObjectList, len(elasticSearchResult.Hits.Hits))
	for i, h := range elasticSearchResult.Hits.Hits {
		object := &dto.SearchObject{}
		if err := json.Unmarshal(*h.Source, &object); err != nil {
			return nil, err
		}
		searchResult.Objects[i] = object
	}

	// aggregations result
	// actors
	actors, err := rollElasticBacket(elasticSearchResult.Aggregations, "actors")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.Actors = actors

	// badges
	badges, err := rollElasticBacket(elasticSearchResult.Aggregations, "badges")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.Badges = badges

	// materials
	materials, err := rollElasticBacket(elasticSearchResult.Aggregations, "materials")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.Materials = materials

	// origin_locations
	originLocation, err := rollElasticBacket(elasticSearchResult.Aggregations, "origin_locations")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.OriginLocations = originLocation

	// statuses
	statuses, err := rollElasticBacket(elasticSearchResult.Aggregations, "statuses")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.Statuses = statuses

	// production_date_intervals
	prDateIntervals, err := rollElasticBacket(elasticSearchResult.Aggregations, "production_date_interval_id")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.ProdutcionNamedIntervals = prDateIntervals

	// production_date_intervals
	prCollections, err := rollElasticBacket(elasticSearchResult.Aggregations, "collection_id")
	if err != nil {
		return nil, err
	}
	searchResult.Filters.Collections = prCollections

	return searchResult, nil
}

func rollElasticBacket(aggr elastic.Aggregations, aggrName string) ([]*dto.ObjectSearchFiltersResult, error) {
	var keys []*dto.ObjectSearchFiltersResult
	terms, found := aggr.Terms(aggrName)
	if !found {
		return nil, fmt.Errorf("can't find aggr by name %s", aggrName)
	}

	for _, backet := range terms.Buckets {
		key, ok := (backet.Key).(float64)
		if !ok {
			return nil, fmt.Errorf("can't convert key to int64 %+v", key)
		}

		filterResult := &dto.ObjectSearchFiltersResult{
			PropertyID: int64(key),
			Cnt:        backet.DocCount,
		}

		keys = append(keys, filterResult)
	}

	return keys, nil
}

var _ SearchClient = (*ElasticSearchClient)(nil)
