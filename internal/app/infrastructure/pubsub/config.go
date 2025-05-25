package pubsub

import "time"

type Config struct {
	Timeout time.Duration `env:"TIMEOUT" envDefault:"1s"`
}
