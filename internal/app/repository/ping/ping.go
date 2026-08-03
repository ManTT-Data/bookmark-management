package ping

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Ping represents the interface of ping repository
//
//go:generate mockery --name Ping --filename ping.go
type Ping interface {
	Ping(ctx context.Context) error
}

type ping struct {
	c *redis.Client
}

// NewPing creates a new ping repository
func NewPing(c *redis.Client) Ping {
	return &ping{
		c: c,
	}
}

// Ping implements all pings functionalities of our application databases
func (p *ping) Ping(ctx context.Context) error {
	return p.c.Ping(ctx).Err()
}
