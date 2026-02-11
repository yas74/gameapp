package dto

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string
	Password    string
}

type UserInfo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}
