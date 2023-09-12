package main

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/urfave/cli"

	"os/signal"
	"syscall"

	"git.softndit.com/collector/backend/loggy"
	"git.softndit.com/collector/backend/npusher/internal"
)

var (
	version = "unknown"

	log = log15.New("app", "npusherdd")

	cliFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "npusherd-conf.toml",
			EnvVar: "NPUSHERD_CONFIG",
		},
		cli.StringFlag{
			Name:   "loglvl",
			Value:  "info",
			EnvVar: "NPUSHERD_LOG_LVL",
		},
	}
)

func init() {
	log.SetHandler(loggy.StdoutLogHandler)
}

func die(msg string, err error) {
	if err != nil {
		log.Crit(msg, "err", err)
	} else {
		log.Crit(msg)
	}
	os.Exit(1)
}

func configureLogger(c *cli.Context) {
	llStr := c.String("loglvl")
	if llStr == "" {
		return
	}
	ll, err := log15.LvlFromString(llStr)
	if err != nil {
		die("Can't parse loglvl flag", err)
	}
	log.SetHandler(log15.LvlFilterHandler(ll, loggy.StdoutLogHandler))
}

func main() {
	app := cli.NewApp()

	app.Name = "npusherd"
	app.Usage = "Push notification service"
	app.Flags = cliFlags
	app.Version = version
	app.Copyright = "Cube Innovations Inc."
	app.Action = func(c *cli.Context) {
		configureLogger(c)
		log.Info("Starting", "version", version)

		srvConfig := &internal.ServerConfig{
			Log: log,
		}
		if err := srvConfig.ReadFromToml(c.String("config")); err != nil {
			die("Error while reading config", err)
		}

		srv, err := internal.NewServer(srvConfig)
		if err != nil {
			die("Can't configure server", err)
		}
		if err := srv.Run(); err != nil {
			die("Can't run server", err)
		}

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		log.Info("Signal catched, exiting...", "signal", <-ch)
	}

	app.Run(os.Args)
}
