package config

var defaultConfig = map[string]interface{}{
	"auth.access_expiration_time":  AccessTokenExpireDuration,
	"auth.refresh_expiration_time": RefreshTokenExpireDuration,
	"auth.refresh_subject":         RefreshTokenSubject,
	"auth.access_subject":          AccessTokenSubject,
}
