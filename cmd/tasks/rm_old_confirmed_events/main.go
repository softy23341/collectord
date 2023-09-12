package main

import (
	"fmt"
	"os"
	"time"

	"git.softndit.com/collector/backend/config"
	dalpg "git.softndit.com/collector/backend/dal/pg"
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

		ago := 30 * (24 * time.Hour)
		if dryRun := c.Bool("dryRun"); dryRun {
			cnt, err := DBM.GetOldEventsCnt(ago)
			if err != nil {
				panic(err)
			}
			fmt.Printf("events cnt %d\n", cnt)
		} else {
			if err := DBM.DeleteOldEvents(ago); err != nil {
				panic(err)
			}
		}

	}
	app.Run(os.Args)
}
