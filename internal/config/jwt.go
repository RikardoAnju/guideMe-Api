package config

import (
	"time"
)

var (
	JWTSecret       = []byte(GetEnv("JWT_SECRET", "secret"))
	JWTExpire       = time.Hour * 24 * 7 // 7 hari
)
