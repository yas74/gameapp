package backofficeuserservice

import "gocasts/gameapp/entity"

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) ListAllUsers() ([]entity.User, error) {
	list := make([]entity.User, 0)

	list = append(list, entity.User{
		ID:          0,
		PhoneNumber: "09155110222",
		Name:        "fff",
		Password:    "fake",
		Role:        entity.AdminRole,
	})

	return list, nil
}
