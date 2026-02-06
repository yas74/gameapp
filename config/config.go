package config

import (
	"gocasts/gameapp/repository/mysql"
	"gocasts/gameapp/service/authservice"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	Mysql      mysql.Config
}
