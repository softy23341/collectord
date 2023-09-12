package client

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"git.softndit.com/collector/backend/cleaver"
	"github.com/inconshreveable/log15"

	"github.com/streadway/amqp"
)

// ErrRequestTimeout TBD
var ErrRequestTimeout = errors.New("cleaver rmqclient request timeout")

// RMQClientConfig TBD
type RMQClientConfig struct {
	Log               log15.Logger
	Servers           []string
	ResponseTTL       time.Duration
	RequestTimeout    time.Duration
	ReconnectInterval time.Duration
}

// DefaultRMQClientConfig TBD
var DefaultRMQClientConfig = RMQClientConfig{
	Log:               log15.Root(),
	ResponseTTL:       2 * time.Minute,
	RequestTimeout:    5 * time.Minute,
	ReconnectInterval: 3 * time.Second,
}

// NewRMQClient TBD
func NewRMQClient(servers ...string) ConnectClient {
	cfg := DefaultRMQClientConfig
	cfg.Servers = servers
	return NewRMQClientWithConfig(cfg)
}

// NewRMQClientWithConfig TBD
func NewRMQClientWithConfig(config RMQClientConfig) ConnectClient {
	client := &RMQClient{config: config, tasks: make(map[string]chan []byte)}
	if _, err := rand.Read(client.requestIDBase[:]); err != nil {
		panic(err)
	}
	return client
}

// RMQClient TBD
type RMQClient struct {
	config RMQClientConfig

	ch            *amqp.Channel
	cbQueueName   string
	requestIDBase [6]byte
	requestNum    uint64
	mu            sync.Mutex

	tasks   map[string]chan []byte
	tasksMu sync.RWMutex
}

// not thread-safe, use mutex lock before call
func (c *RMQClient) nextRequestID() string {
	var buf [14]byte // 6 for requestBase + 8 for requestNum
	copy(buf[:6], c.requestIDBase[:])
	binary.LittleEndian.PutUint64(buf[6:], c.requestNum)
	c.requestNum++
	return base64.StdEncoding.EncodeToString(buf[:])
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

	args := make(amqp.Table)
	if c.config.ResponseTTL > 0 {
		args["x-message-ttl"] = int32(c.config.ResponseTTL / time.Millisecond)
	}
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		args,  // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.ch, c.cbQueueName = ch, q.Name
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

	go func() {
		for m := range msgs {
			c.tasksMu.RLock()
			resCh, found := c.tasks[m.CorrelationId]
			c.tasksMu.RUnlock()
			if !found {
				continue
			}
			select {
			case resCh <- m.Body:
			default:
			}
		}
	}()

	return nil
}

func (c *RMQClient) processTask(queueName string, task, result interface{}) error {
	reqData, err := json.Marshal(task)
	if err != nil {
		return err
	}

	var (
		respCh = make(chan []byte)
		id     string
	)
	c.tasksMu.Lock()
	id = c.nextRequestID()
	c.tasks[id] = respCh
	c.tasksMu.Unlock()

	defer func() {
		c.tasksMu.Lock()
		delete(c.tasks, id)
		c.tasksMu.Unlock()
	}()

	c.mu.Lock()
	var (
		ch          = c.ch
		cbQueueName = c.cbQueueName
	)
	c.mu.Unlock()

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: id,
			ReplyTo:       cbQueueName,
			Body:          reqData,
		})
	if err != nil {
		return err
	}

	timeoutCh := time.Tick(c.config.RequestTimeout)
	select {
	case resData := <-respCh:
		if err = json.Unmarshal(resData, &result); err != nil {
			return err
		}
		return nil
	case <-timeoutCh:
		return ErrRequestTimeout
	}

}

// Resize TBD
func (c *RMQClient) Resize(task *cleaver.ResizeTask) ([]*cleaver.TransformResult, error) {
	var result cleaver.StatusedResizeResult
	if err := c.processTask(cleaver.RMQQueueResizeJSON, task, &result); err != nil {
		return nil, err
	}
	if result.Status.Code != cleaver.StatusOK {
		return nil, result.Status.ToErr()
	}
	return result.Transforms, nil
}

// Copy TBD
func (c *RMQClient) Copy(task *cleaver.CopyTask) (*cleaver.CopyResult, error) {
	var result cleaver.StatusedCopyResult
	if err := c.processTask(cleaver.RMQQueueCopyJSON, task, &result); err != nil {
		return nil, err
	}
	if result.Status.Code != cleaver.StatusOK {
		return nil, result.Status.ToErr()
	}
	return result.Result, nil
}
