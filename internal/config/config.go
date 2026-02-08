package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env                 string
	HTTPAddr            string
	DBURL               string
	RedisAddr           string
	RedisPassword       string
	RedisDB             int
	ChallengeTTL        time.Duration
	EnrollmentTTL       time.Duration
	NonceTTL            time.Duration
	RateLimitPerMinute  int
	EnableTOTP          bool
	DefaultTenantPolicy string
}

func Load() Config {
	return Config{
		Env:                 getenv("TRUSTPIN_ENV", "dev"),
		HTTPAddr:            getenv("TRUSTPIN_HTTP_ADDR", ":8080"),
		DBURL:               getenv("TRUSTPIN_DB_URL", "postgres://trustpin:trustpin@postgres:5432/trustpin?sslmode=disable"),
		RedisAddr:           getenv("TRUSTPIN_REDIS_ADDR", "redis:6379"),
		RedisPassword:       getenv("TRUSTPIN_REDIS_PASSWORD", ""),
		RedisDB:             getenvInt("TRUSTPIN_REDIS_DB", 0),
		ChallengeTTL:        getenvDuration("TRUSTPIN_CHALLENGE_TTL", 2*time.Minute),
		EnrollmentTTL:       getenvDuration("TRUSTPIN_ENROLLMENT_TTL", 10*time.Minute),
		NonceTTL:            getenvDuration("TRUSTPIN_NONCE_TTL", 5*time.Minute),
		RateLimitPerMinute:  getenvInt("TRUSTPIN_RATE_LIMIT_PER_MINUTE", 120),
		EnableTOTP:          getenvBool("TRUSTPIN_ENABLE_TOTP", false),
		DefaultTenantPolicy: getenv("TRUSTPIN_DEFAULT_POLICY", "default"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getenvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
