package main

import (
	"fmt"
	"net/url"
	"os"

	"git.softndit.com/collector/backend/config"
	"git.softndit.com/collector/backend/dal"
	dalpg "git.softndit.com/collector/backend/dal/pg"
	"git.softndit.com/collector/backend/dto"
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
	DBM    dal.TrManager
	dryRun bool
}

func migrateURL(oldURL url.URL) (*url.URL, error) {
	oldURL.Scheme = "http"
	oldURL.Host = "asd"
	return &oldURL, nil
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

		// dryrun
		dryRun := c.Bool("dryRun")

		perPage := 50

		context := &context{DBM: DBM, dryRun: dryRun}
		for paginator := dto.NewPagePaginator(0, int16(perPage)); ; paginator.NextPage() {
			medias, err := DBM.GetMediaByPage(dto.MediaTypeList, paginator)
			if err != nil {
				panic(err)
			}

			// TODO: bulk update
			for _, media := range medias {
				if err := migrateMediaURL(context, media); err != nil {
					panic(err)
				}
				if !dryRun {
					if err := DBM.UpdateMedia(media); err != nil {
						panic(err)
					}
				}
			}

			if mlen := len(medias); mlen == 0 || mlen < perPage {
				break
			}
		}

	}
	app.Run(os.Args)
}

func migrateMediaURL(context *context, media *dto.Media) error {
	if photo := media.Photo; photo != nil {
		for _, variant := range photo.Variants {
			URLString := variant.URI
			URL, err := url.Parse(URLString)
			if err != nil {
				return err
			}

			newURL, err := migrateURL(*URL)
			if err != nil {
				return err
			}

			newURLString := newURL.String()
			variant.URI = newURLString

			if context.dryRun {
				fmt.Printf("from: %s, to: %s\n", URLString, newURLString)
			}
		}
	}

	if doc := media.Document; doc != nil {
		URLString := doc.URI
		URL, err := url.Parse(URLString)
		if err != nil {
			return err
		}

		newURL, err := migrateURL(*URL)
		if err != nil {
			return err
		}

		newURLString := newURL.String()
		doc.URI = newURLString

		if context.dryRun {
			fmt.Printf("from: %s, to: %s\n", URLString, newURLString)
		}
	}

	return nil
}
