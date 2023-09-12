package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"

	"git.softndit.com/collector/backend/npusher"
	"git.softndit.com/collector/backend/npusher/client"
	"github.com/inconshreveable/log15"
)

var (
	version string

	cliFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Value: "amqp://guest:guest@localhost:5672/",
			Usage: "RabbitMQ server uri",
		},
		cli.StringFlag{
			Name:  "queue",
			Value: "np-send.json",
			Usage: "RabbitMQ queue name",
		},
		cli.StringFlag{
			Name:  "token",
			Usage: "Push notification target app/device token",
		},
		cli.StringFlag{
			Name:  "message",
			Value: "Default message from args",
			Usage: "Push notification message",
		},
		cli.BoolFlag{
			Name:  "sandbox",
			Usage: "Push notification sandbox flag",
		},
	}
)

func main() {
	app := cli.NewApp()

	app.Flags = cliFlags
	app.Name = "npusherd-cli"
	app.Usage = "Npusher RabbitMQ test client"
	app.Version = version
	app.Copyright = "Cube Innovations Inc."
	app.Action = func(c *cli.Context) {
		cfg := client.RMQClientConfig{
			Log:               log15.Root(),
			Servers:           []string{c.String("server")},
			ReconnectInterval: 3 * time.Second,
			QueueName:         c.String("queue"),
		}

		rabbitClient := client.NewRMQClientWithConfig(cfg)
		if err := rabbitClient.Connect(); err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}

		token := c.String("token")

		notify := npusher.APNSNotification{
			Alert: c.String("message"),
		}

		err := rabbitClient.SendPush(token, c.Bool("sandbox"), notify)

		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
		fmt.Println("Was send")
	}

	app.Run(os.Args)
}
