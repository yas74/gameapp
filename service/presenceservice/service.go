package presenceservice

import (
	"context"
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/pkg/richerror"
	"time"
)

type Config struct {
	Prefix         string        `koanf:"prefix"`
	ExpirationTime time.Duration `koanf:"expiration_time"`
}

type Repository interface {
	Upsert(ctx context.Context, key string, timeStamp int64, expTime time.Duration) error
}

type Service struct {
	config Config
	repo   Repository
}

func New(config Config, repo Repository) Service {
	return Service{config: config, repo: repo}
}

func (s Service) Upsert(ctx context.Context, req dto.UpsertPresenceRequest) (dto.UpsertPresenceResponse, error) {
	const op = richerror.Op("presenceservice.Upsert")

	key := fmt.Sprintf("%s:%d", s.config.Prefix, req.UserID)

	err := s.repo.Upsert(ctx, key, req.TimeStamp, s.config.ExpirationTime)
	if err != nil {
		return dto.UpsertPresenceResponse{},
			richerror.New(op).WithErr(err)
	}

	return dto.UpsertPresenceResponse{}, nil
}
