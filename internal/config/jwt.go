package config

import "time"

type JWTConfig struct {
	Secret         string
	AccessTokenTTL time.Duration
}
