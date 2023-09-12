package main

import (
	"os"

	"git.softndit.com/collector/backend/config"
	"git.softndit.com/collector/backend/dal"
	dalpg "git.softndit.com/collector/backend/dal/pg"
	"git.softndit.com/collector/backend/dto"
	resource "git.softndit.com/collector/backend/resources"
	"git.softndit.com/collector/backend/services"
	"github.com/urfave/cli"
)

var cliFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config",
		Value: "collectord.toml",
	},
	cli.BoolFlag{
		Name: "dryRun",
	},
}

type context struct {
	DBM          dal.TrManager
	SearchClient services.SearchClient
}

func main() {
	app := cli.NewApp()
	app.Flags = cliFlags

	app.Action = func(c *cli.Context) {
		configPath := c.String("config")

		mainConfig := &config.Config{}
		if err := mainConfig.ReadConfig(configPath); err != nil {
			panic(err)
		}

		instanceHolder, err := config.NewInstanceHolder(mainConfig)
		if err != nil {
			panic(err)
		}

		// DBM
		serverConfig := mainConfig.ServerConfig
		pgConfig, err := instanceHolder.GetPgConfig(serverConfig.DBM.PgRef)
		if err != nil {
			panic(err)
		}
		DBM, err := dalpg.NewManager(&dalpg.ManagerContext{
			PoolConfig: pgConfig,
		})
		if err != nil {
			panic(err)
		}

		// search // elastic
		elasticConfig, err := instanceHolder.GetElasticSearchConfig(serverConfig.SearchClient.ElasticSearchRef)
		if err != nil {
			panic(err)
		}

		searchClient, err := services.NewElasticSearchClient(
			&services.ElasticSearchClientContext{
				URL:         elasticConfig.URL,
				ObjectIndex: serverConfig.SearchClient.ObjectIndex,
			},
		)
		if err != nil {
			panic(err)
		}

		perPage := 50

		for paginator := dto.NewPagePaginator(0, int16(perPage)); ; paginator.NextPage() {
			objects, err := DBM.GetObjects(paginator)
			if err != nil {
				panic(err)
			}

			if err := reindex(&context{DBM: DBM, SearchClient: searchClient}, objects); err != nil {
				panic(err)
			}

			if mlen := len(objects); mlen == 0 || mlen < perPage {
				break
			}
		}

	}
	app.Run(os.Args)
}

func reindex(context *context, objectsList dto.ObjectList) error {
	searchObjects, err := resource.PrepareToReindex(&resource.ReindexContext{DBM: context.DBM}, objectsList)
	if err != nil {
		return err
	}

	if len(searchObjects) > 0 {
		return context.SearchClient.BulkObjectIndex(searchObjects)
	}

	return nil
}
