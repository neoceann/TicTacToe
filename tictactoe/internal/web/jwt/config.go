package _jwt

import "time"

type JWTConf struct {
	SecretAccessToken string
	AccessTTL time.Duration
	SecretRefreshToken string
	RefreshTTL time.Duration
}

func NewJWTConfig() *JWTConf {
	return &JWTConf{SecretAccessToken: "actest",
					AccessTTL: 1*time.Minute,
					SecretRefreshToken: "reftest",
					RefreshTTL: 1*time.Hour,
}
}