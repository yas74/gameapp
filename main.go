package main

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/delivery/httpserver"
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/service/authservice"
	"gocasts/gameapp/service/userservice"
	"time"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	fmt.Println("start echo server")

	authConfig := authservice.Config{
		SignKey:               JwtSignKey,
		AccessExpirationTime:  AccessTokenExpireDuration,
		RefreshExpirationTime: RefreshTokenExpireDuration,
		AccessSubject:         AccessTokenSubject,
		RefreshSubject:        RefreshTokenSubject,
	}

	mysqlCfg := mysql.Config{
		Username: "gameapp",
		Password: "gameappt0lk2o20",
		Port:     3308,
		Host:     "localhost",
		DBName:   "gameapp_db",
	}

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth:       authConfig,
		Mysql:      mysqlCfg,
	}

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)

	server.Serve()
}

// func userLoginHandler(writer http.ResponseWriter, req *http.Request) {
// 	if req.Method != http.MethodPost {
// 		fmt.Fprintf(writer, `{"error":"invalid method"}`)
// 	}

// 	data, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		writer.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))

// 		return
// 	}

// 	var lReq userservice.LoginRequest
// 	err = json.Unmarshal(data, &lReq)
// 	if err != nil {
// 		writer.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))

// 		return
// 	}

// 	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
// 		AccessTokenExpireDuration, RefreshTokenExpireDuration)

// 	mysqlRepo := mysql.New(
// 		mysql.Config{
// 			Username: "root",
// 			Password: "gameappRoo7t0lk2o20",
// 			Port:     3308,
// 			Host:     "127.0.0.1",
// 			DBName:   "gameapp_db",
// 		},
// 	)
// 	userSvc := userservice.New(authSvc, mysqlRepo)

// 	resp, err := userSvc.Login(lReq)
// 	if err != nil {
// 		writer.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))

// 		return
// 	}

// 	data, err = json.Marshal(resp)
// 	if err != nil {
// 		writer.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))

// 		return
// 	}
// 	writer.Write(data)

// }

// func userProfileHandler(writer http.ResponseWriter, req *http.Request) {
// 	if req.Method != http.MethodGet {
// 		fmt.Fprintf(writer, `{"error":"invalid method"}`)
// 	}

// 	authSvc := authservice.New(JwtSignKey, AccessTokenSubject, RefreshTokenSubject,
// 		AccessTokenExpireDuration, RefreshTokenExpireDuration)

// 	authToken := req.Header.Get("Authorization")

// 	claims, err := authSvc.ParseToken(authToken)
// 	if err != nil {
// 		fmt.Fprintf(writer, `{"error":"token is not valid"}`)

// 		return
// 	}

// 	mysqlRepo := mysql.New(
// 		mysql.Config{
// 			Username: "root",
// 			Password: "gameappRoo7t0lk2o20",
// 			Port:     3308,
// 			Host:     "127.0.0.1",
// 			DBName:   "gameapp_db",
// 		},
// 	)
// 	userSvc := userservice.New(authSvc, mysqlRepo)

// 	pReq := userservice.ProfileRequest{UserID: claims.UserID}
// 	resp, err := userSvc.Profile(pReq)
// 	if err != nil {
// 		writer.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))

// 		return
// 	}

// 	data, err := json.Marshal(resp)
// 	writer.Write(data)

// }

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.New(authSvc, mysqlRepo)

	return authSvc, userSvc

}
