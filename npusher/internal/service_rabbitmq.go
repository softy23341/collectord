package internal

import (
	"encoding/json"
	"errors"
	"time"

	"git.softndit.com/collector/backend/npusher"

	"github.com/BurntSushi/toml"
	"github.com/streadway/amqp"
)

func init() {
	RegesterService("rabbitmq", newRabitmqService)
}

type rmqServiceConfig struct {
	QueueName string `toml:"queue_name"`
	Servers   []string
	Prefetch  int
	TaskTTL   int
}

const queueName = "np-send.json"

type rmqService struct {
	ctx    *ServiceCtx
	config rmqServiceConfig
}

func newRabitmqService(ctx *ServiceCtx) (Service, error) {
	s := &rmqService{ctx: ctx}
	if err := s.configure(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *rmqService) configure(ctx *ServiceCtx) error {
	if err := toml.PrimitiveDecode(*ctx.Config, &s.config); err != nil {
		return err
	}

	if len(s.config.Servers) == 0 {
		return errors.New("no servers specified")
	}

	if s.config.Prefetch < 0 {
		s.config.Prefetch = 0
	}

	if s.config.TaskTTL < 1 {
		s.config.TaskTTL = 60
	}

	if s.config.QueueName == "" {
		s.config.QueueName = queueName
	}

	return nil
}

func (s *rmqService) init() error {
	var (
		conn       *amqp.Connection
		connectErr error
	)
	for _, srv := range s.config.Servers {
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

	queueArgs := make(amqp.Table)
	if s.config.TaskTTL > 0 {
		queueArgs["x-message-ttl"] = int32(s.config.TaskTTL * 1000)
	}
	q, err := ch.QueueDeclare(
		s.config.QueueName, // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		queueArgs,          // arguments
	)
	if err != nil {
		return err
	}

	if err := ch.Qos(s.config.Prefetch, 0, false); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	// naive try to reconnect
	closeCh := conn.NotifyClose(make(chan *amqp.Error))
	go func() {
		closeErr := <-closeCh
		s.ctx.Log.Error("connection to RabbitMQ server closed", "err", closeErr)
		var n = 0
		for {
			n++
			s.ctx.Log.Info("try to reconnect", "try", n)
			err := s.init()
			if err == nil {
				s.ctx.Log.Info("reconnected")
				return
			}
			s.ctx.Log.Error("reconnect error", "try", n, "err", err)
			time.Sleep(5 * time.Second)
		}

	}()

	go func() {
		for m := range msgs {
			go func(d amqp.Delivery) {
				defer d.Ack(false)
				s.send(d.Body)
			}(m)
		}
	}()

	s.ctx.Log.Debug("Rabbitmq-service successful init")
	return nil
}

func (s *rmqService) send(body []byte) error {
	receiveTask := &npusher.NotificationTask{}
	if err := json.Unmarshal(body, receiveTask); err != nil {
		s.ctx.Log.Error("Can't unmarshal bytes",
			"str Bytes", string(body))

		return err
	}

	if err := s.ctx.Provider.Send(receiveTask); err != nil {
		s.ctx.Log.Error("can't send msg", "err", err)
		return err
	}
	return nil
}

func (s *rmqService) Run() error {
	return s.init()
}
