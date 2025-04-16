package pubsub

import "time"

type Config struct {
	Timeout time.Duration `default:"1s"`
}
