package matchingservice

import (
	"gocasts/gameapp/dto"
	"gocasts/gameapp/entity"
	"gocasts/gameapp/pkg/richerror"
	"time"
)

type Config struct {
	WaitingTimeout time.Duration `koanf:"timeout"`
}

type Repository interface {
	AddToWaitingList(userID uint, category entity.Category) error
}

type Service struct {
	config Config
	repo   Repository
}

func New(config Config, repo Repository) Service {
	return Service{config: config, repo: repo}
}

func (s Service) AddToWaitingList(req dto.AddToWaitingListRequest) (dto.AddToWaitingListResponse, error) {
	op := richerror.Op("matchingservice.AddToWaitingList")

	// add the user to the waiting list for the given category if they already dont exist
	err := s.repo.AddToWaitingList(req.UserID, req.Category)
	if err != nil {
		return dto.AddToWaitingListResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return dto.AddToWaitingListResponse{Timeout: s.config.WaitingTimeout}, nil
}

func (s Service) MatchWaitedUsers(req dto.MatchWaitedUsersRequest) (dto.MatchWaitedUsersResponse, error) {
	return dto.MatchWaitedUsersResponse{}, nil
}
