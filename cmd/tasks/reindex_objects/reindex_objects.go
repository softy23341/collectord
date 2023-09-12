package main

import (
	dalpg "git.softndit.com/collector/backend/dal/pg"
	resource "git.softndit.com/collector/backend/resources"
	"git.softndit.com/collector/backend/services"
	pgx "github.com/jackc/pgx"
	log15 "gopkg.in/inconshreveable/log15.v2"
)

type objectReindexer struct {
	Context      *resource.Context
	Log          log15.Logger
	SearchClient services.SearchClient
}

func main() {
	// logger
	logger := log15.New("app", "object reindexer")

	// DBM
	DBM, err := dalpg.NewManager(&dalpg.ManagerContext{
		Log: logger,
		PoolConfig: &pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "127.0.0.1",
				Database: "collector_development",
				User:     "collector_app",
				Password: "collectordevpassstrong",
			},
			MaxConnections: 10,
		},
	})
	if err != nil {
		panic(err)
	}

	// search // elastic

	searchClient, err := services.NewElasticSearchClient(
		&services.ElasticSearchClientContext{
			URL:         "http://127.0.0.1:9200",
			ObjectIndex: "object_dev",
		},
	)
	if err != nil {
		panic(err)
	}

	reindexer := &objectReindexer{
		Context: &resource.Context{
			DBM:          DBM,
			SearchClient: searchClient,
		},
		SearchClient: searchClient,
	}
	if err := reindexer.Run(); err != nil {
		panic(err)
	}
}

func (r *objectReindexer) Run() error {
	query := &services.ScrollSearchQuery{}
	return r.SearchClient.ScrollThrought(query, r.Context.ReindexObjects)
}
