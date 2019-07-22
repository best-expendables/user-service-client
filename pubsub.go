package userclient

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	redis "gopkg.in/redis.v5"
)

// PubSub allows to subscribe to user service messages
// It's NOT safe for concurrent use by multiple goroutines.
type PubSub interface {
	ReceiveMsg() (Message, error)
}

// Message pubsub message from user service
type Message struct {
	Source  string
	Event   string
	Payload string
}

// SubscribeRedis returns PubSub implemented on top of redis
func SubscribeRedis(c *redis.Client) (PubSub, error) {
	pubsub, err := c.Subscribe(channelName)
	if err != nil {
		return nil, err
	}
	return &redisPubSub{
		pubsub:         pubsub,
		receiveTimeout: time.Minute,
	}, nil
}

const channelName = "events-channel"

type redisPubSub struct {
	pubsub         *redis.PubSub
	receiveTimeout time.Duration
}

func (c *redisPubSub) ReceiveMsg() (Message, error) {
	for {
		msgi, err := c.pubsub.ReceiveTimeout(c.receiveTimeout)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				err = c.pubsub.Ping()
				if err != nil {
					err = fmt.Errorf("PubSub.Ping failed: %s", err)
				} else {
					continue
				}
			}
			return Message{}, err
		}

		switch msg := msgi.(type) {
		case *redis.Subscription:
			// Ignore.
		case *redis.Pong:
			// Ignore.
		case *redis.Message:
			return c.parseMsg(msg)
		default:
			return Message{}, fmt.Errorf("unknown message: %T", msgi)
		}
	}
}

func (c *redisPubSub) parseMsg(m *redis.Message) (Message, error) {
	// example: "user:updated:uuid"
	parts := strings.Split(m.Payload, ":")
	if len(parts) != 3 {
		return Message{}, errors.New("incorrect message")
	}
	return Message{
		Source:  parts[0],
		Event:   parts[1],
		Payload: parts[2],
	}, nil
}
