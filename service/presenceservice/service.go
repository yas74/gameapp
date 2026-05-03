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
	GetPresence(ctx context.Context, key string) (int64, error)
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

func (s Service) GetPresence(ctx context.Context, req dto.GetPresenceRequest) (dto.GetPresenceResponse, error) {
	const op = richerror.Op("presenceservice.GetPresence")

	presenceList := make([]dto.GetPresenceItem, 0)

	for _, i := range req.UserIDs {
		timeStamp, err := s.repo.GetPresence(ctx, s.GetPresenceKey(i))
		if err != nil {
			// TODO - log error
			// TODO - update metrics
			continue
		}

		presenceItem := dto.GetPresenceItem{
			UserID:    i,
			TimeStamp: timeStamp,
		}

		presenceList = append(presenceList, presenceItem)
	}

	return dto.GetPresenceResponse{Items: presenceList}, nil
}

func (s Service) GetPresenceKey(userID uint) string {
	return fmt.Sprintf("%s:%d", s.config.Prefix, userID)
}
