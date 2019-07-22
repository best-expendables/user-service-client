package userclient

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"gopkg.in/redis.v5"
)

var redisAddr string

func init() {
	flag.StringVar(&redisAddr, "redisAddr", "", "")
}

func TestRedisPubSub(t *testing.T) {
	if redisAddr == "" {
		t.Skip("skipped because redisAddr isn't provided")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	pubsub, err := SubscribeRedis(redisClient)
	if err != nil {
		t.Fatalf("subscribe failed with err: %s", err)
	}

	// speedup tests
	pubsub.(*redisPubSub).receiveTimeout = 50 * time.Millisecond

	go func() {
		redisClient.Publish(channelName, "user:updated:123")
	}()

	msg, err := pubsub.ReceiveMsg()
	if err != nil {
		t.Fatalf("ReceiveMsg failed with err: %s", err)
	}
	if !reflect.DeepEqual(msg, Message{
		Source:  "user",
		Event:   "updated",
		Payload: "123",
	}) {
		t.Errorf("incorrect message: %+v", msg)
	}

	go func() {
		redisClient.Publish(channelName, "incorrect")
	}()

	_, err = pubsub.ReceiveMsg()
	if err == nil || err.Error() != "incorrect message" {
		t.Errorf("expected incorrect message error, but error is %s", err)
	}

	// check that ping works
	go func() {
		time.Sleep(150 * time.Millisecond)
		redisClient.Publish(channelName, "user:updated:123")
	}()

	_, err = pubsub.ReceiveMsg()
	if err != nil {
		t.Errorf("ReceiveMsg failed with err: %s", err)
	}
}
