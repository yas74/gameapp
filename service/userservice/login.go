package userservice

import (
	"fmt"
	"gocasts/gameapp/dto"
	"gocasts/gameapp/pkg/richerror"
)

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"
	// TODO - it would be better to use two separate methods for existance check and GetUserByPhoneNumber

	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("username or password is not correct")
	}

	if user.Password != getMD5Hash(req.Password) {
		return dto.LoginResponse{}, fmt.Errorf("username or password is not correct")
	}

	// jwt generate
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return dto.LoginResponse{
		User: dto.UserInfo{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
		},
		Tokens: dto.Tokens{AccessToken: accessToken,
			RefreshToken: refreshToken},
	}, nil
}
