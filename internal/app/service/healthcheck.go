package service

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/app/repository"
)

var (
	errRedisNotAvailable = errors.New("redis not available")
)

// Healthcheck represents the interface of healthcheck healthcheck
//
//go:generate mockery --name Healthcheck --filename healthcheck.go
type Healthcheck interface {
	HealthCheck(ctx context.Context) (*model.Message, error)
}

// healthcheck represents the structure of a healthcheck service
type healthcheck struct {
	serviceName string
	instanceID  string
	pingRepo    repository.Ping
}

// NewHealthcheck returns a new healthcheck service
func NewHealthcheck(serviceName, instanceID string, pingRepo repository.Ping) Healthcheck {
	s := &healthcheck{
		serviceName: serviceName,
		instanceID:  instanceID,
		pingRepo:    pingRepo,
	}

	return s
}

// HealthCheck returns a HealthCheck message
func (s *healthcheck) HealthCheck(ctx context.Context) (*model.Message, error) {
	err := s.pingRepo.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(errRedisNotAvailable, err.Error())
	}

	return &model.Message{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  s.instanceID,
	}, nil
}
