package internal

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"git.softndit.com/collector/backend/cleaver"

	"github.com/BurntSushi/toml"
	logext "github.com/inconshreveable/log15/ext"
	"github.com/streadway/amqp"
)

func init() {
	RegisterService("rabbitmq", newRMQService)
}

func newRMQService(ctx *ServiceContext) (Service, error) {
	s := new(rmqService)
	s.ctx = ctx
	if err := s.configure(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

type rmqService struct {
	ctx    *ServiceContext
	config struct {
		Servers  []string
		Prefetch int
		TaskTTL  int
	}
}

func (s *rmqService) configure(ctx *ServiceContext) error {
	if err := toml.PrimitiveDecode(*ctx.config, &s.config); err != nil {
		return err
	}

	if len(s.config.Servers) == 0 {
		return errors.New("no servers specified")
	}

	if s.config.Prefetch < 0 {
		s.config.Prefetch = 0
	}

	if s.config.TaskTTL < 1 {
		s.config.TaskTTL = 120
	}

	return nil
}

func dump(b []byte) string {
	return strings.Replace(string(b), "\"", "'", -1)
}

func (s *rmqService) resize(taskData []byte) ([]byte, error) {
	var (
		lg   = s.ctx.log.New("id", logext.RandId(8), "type", "resize")
		task cleaver.ResizeTask
	)

	lg.Debug("new task", "task", dump(taskData))
	if err := json.Unmarshal(taskData, &task); err != nil {
		lg.Error("can't decode task", "err", err)
		return nil, err
	}

	transforms, err := s.ctx.executor.Resize(&task)
	if err != nil {
		lg.Error("task error", "err", err)
	}

	result := cleaver.StatusedResizeResult{
		Status:     cleaver.NewStatusFromErr(err),
		Transforms: transforms,
	}

	resultData, err := json.Marshal(&result)
	if err != nil {
		lg.Error("can't encode task result", "err", err)
		return nil, err
	}
	lg.Debug("task finished", "info", dump(resultData))

	return resultData, nil
}

func (s *rmqService) copy(taskData []byte) ([]byte, error) {
	var (
		lg   = s.ctx.log.New("id", logext.RandId(8), "type", "copy")
		task cleaver.CopyTask
	)

	lg.Debug("new task", "task", dump(taskData))
	if err := json.Unmarshal(taskData, &task); err != nil {
		lg.Error("can't decode task", "err", err)
		return nil, err
	}

	copyResult, err := s.ctx.executor.Copy(&task)
	if err != nil {
		lg.Error("task error", "err", err)
	}

	result := cleaver.StatusedCopyResult{
		Status: cleaver.NewStatusFromErr(err),
		Result: copyResult,
	}

	resultData, err := json.Marshal(&result)
	if err != nil {
		lg.Error("can't encode task result", "err", err)
		return nil, err
	}
	lg.Debug("task finished", "info", dump(resultData))

	return resultData, nil
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
	declareQueue := func(name string) (<-chan amqp.Delivery, error) {
		q, err := ch.QueueDeclare(
			name,      // name
			false,     // durable
			false,     // delete when usused
			false,     // exclusive
			false,     // no-wait
			queueArgs, // arguments
		)
		if err != nil {
			return nil, err
		}

		if err := ch.Qos(s.config.Prefetch, 0, false); err != nil {
			return nil, err
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
			return nil, err
		}

		return msgs, nil
	}

	resizeMsgs, err := declareQueue(cleaver.RMQQueueResizeJSON)
	if err != nil {
		return err
	}

	copyMsgs, err := declareQueue(cleaver.RMQQueueCopyJSON)
	if err != nil {
		return err
	}

	// naive try to reconnect
	closeCh := conn.NotifyClose(make(chan *amqp.Error))
	go func() {
		closeErr := <-closeCh
		s.ctx.log.Error("connection to RabbitMQ server closed", "err", closeErr)
		var n = 0
		for {
			n++
			s.ctx.log.Info("try to reconnect", "try", n)
			err := s.init()
			if err == nil {
				s.ctx.log.Info("reconnected")
				return
			}
			s.ctx.log.Error("reconnect error", "try", n, "err", err)
			time.Sleep(5 * time.Second)
		}

	}()

	consumeMsgs := func(msgs <-chan amqp.Delivery, fn func([]byte) ([]byte, error)) {
		for m := range msgs {
			go func(d amqp.Delivery) {
				resp, err := fn(d.Body)
				if err != nil {
					d.Ack(false)
					return
				}

				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "application/json",
						CorrelationId: d.CorrelationId,
						Body:          resp,
					})
				if err != nil {
					s.ctx.log.Error("failed to publish a message", "err", err)
				}

				d.Ack(false)
			}(m)
		}
	}

	go consumeMsgs(resizeMsgs, s.resize)
	go consumeMsgs(copyMsgs, s.copy)

	return nil
}

func (s *rmqService) Run() error {
	return s.init()
}
