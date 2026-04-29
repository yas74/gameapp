package userservice

import (
	"context"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/pkg/richerror"
)

func (s Service) Profile(ctx context.Context, req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return dto.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return dto.ProfileResponse{Name: user.Name}, nil

}
