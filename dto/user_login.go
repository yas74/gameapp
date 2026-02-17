package dto

//data transfer object

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string
}

type LoginResponse struct {
	User   UserInfo `json:"user"`
	Tokens Tokens   `json:"tokens"`
}
