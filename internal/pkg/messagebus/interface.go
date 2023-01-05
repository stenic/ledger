package messagebus

import (
	"github.com/go-redis/redis/v8"
)

type MessageBus interface {
	Publish(channel, msg string) error
	Consume(channel string) <-chan *redis.Message
	Close()
}
