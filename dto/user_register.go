package dto

//data transfer object

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string
	Password    string
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}
