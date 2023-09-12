package client

import (
	"encoding/json"
	"sync"
	"time"

	"git.softndit.com/collector/backend/npusher"

	"github.com/inconshreveable/log15"
	"github.com/streadway/amqp"
)

// RMQClientConfig TBD
type RMQClientConfig struct {
	Log               log15.Logger
	Servers           []string
	ReconnectInterval time.Duration
	QueueName         string
}

// DefaultRMQClientConfig TBD
var DefaultRMQClientConfig = RMQClientConfig{
	Log:               log15.Root(),
	ReconnectInterval: 3 * time.Second,
	QueueName:         "np-send.json",
}

// NewRMQClient TBD
func NewRMQClient(servers ...string) *RMQClient {
	cfg := DefaultRMQClientConfig
	cfg.Servers = servers
	return NewRMQClientWithConfig(cfg)
}

// NewRMQClientWithConfig TBD
func NewRMQClientWithConfig(config RMQClientConfig) *RMQClient {
	return &RMQClient{config: config}
}

// RMQClient TBD
type RMQClient struct {
	config RMQClientConfig

	ch *amqp.Channel
	mu sync.Mutex
}

// Connect TBD
func (c *RMQClient) Connect() error {
	var (
		conn       *amqp.Connection
		connectErr error
	)
	for _, srv := range c.config.Servers {
		conn, connectErr = amqp.Dial(srv)
		if connectErr == nil {
			break
		}
	}

	if connectErr != nil {
		return connectErr
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.ch = ch
	c.mu.Unlock()

	closeCh := conn.NotifyClose(make(chan *amqp.Error))
	go func() {
		closeErr := <-closeCh
		c.config.Log.Error("connection to RabbitMQ server closed", "err", closeErr)
		for n := 0; false; n++ {
			c.config.Log.Info("try to reconnect", "try", n)
			err := c.Connect()
			if err == nil {
				c.config.Log.Info("reconnected")
				return
			}
			c.config.Log.Error("reconnect error", "try", n, "err", err)
			time.Sleep(c.config.ReconnectInterval)
		}
	}()

	return nil
}

// SendPush TBD
func (c *RMQClient) SendPush(token string, sandbox bool, n npusher.Notification) error {
	notificationTask, err := c.createTask(token, sandbox, n)
	if err != nil {
		return err
	}

	taskBody, err := json.Marshal(notificationTask)
	if err != nil {
		return err
	}

	if err := c.publish(taskBody); err != nil {
		return err
	}

	return nil
}

// createTask TBD
func (c *RMQClient) createTask(token string, sandbox bool, n npusher.Notification) (*npusher.NotificationTask, error) {
	payload, err := n.PayloadJSON()
	if err != nil {
		return nil, err
	}
	return &npusher.NotificationTask{
		Token:          token,
		Type:           n.Type(),
		Sandbox:        sandbox,
		Payload:        payload,
		DefaultMessage: n.DefaultMessage(),
	}, nil
}

// Publish TBD
func (c *RMQClient) publish(body []byte) error {
	c.mu.Lock()
	var ch = c.ch
	c.mu.Unlock()

	err := ch.Publish(
		"",                 // exchange
		c.config.QueueName, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	return err
}
