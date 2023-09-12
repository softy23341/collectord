package main

import (
	"os"
	"os/signal"
	"syscall"

	"git.softndit.com/collector/backend/cleaver/internal"
	"git.softndit.com/collector/backend/loggy"
	"github.com/inconshreveable/log15"

	"github.com/urfave/cli"
)

var (
	version = "unknown"

	log = log15.New("app", "cleaverdd")

	cliFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Value:  "cleaverd.toml",
			EnvVar: "CL_CONFIG",
			Usage:  "config file name",
		},
		cli.StringFlag{
			Name:   "loglvl",
			Value:  "info",
			EnvVar: "CL_LOGLVL",
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
	app.Name = "cleaverd"
	app.Usage = "Image resize server"
	app.Flags = cliFlags
	app.Version = version
	app.Copyright = "Cube Innovations Inc."
	app.Authors = []cli.Author{cli.Author{Name: "Dmitry Panin", Email: "dmitry.panin@yahoo.com"}}

	app.Action = func(c *cli.Context) {
		configureLogger(c)
		log.Info("Starting", "version", version)

		server, err := internal.NewServer(c.String("config"), log)
		if err != nil {
			die("Can't configure server", err)
		}

		if err := server.Run(); err != nil {
			die("Can't start server", err)
		}

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		log.Info("Signal catched, exiting...", "signal", <-ch)
	}

	app.Run(os.Args)
}
