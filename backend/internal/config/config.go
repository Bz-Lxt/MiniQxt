package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env            string
	HTTPAddr       string
	MySQLDSN       string
	JWTSecret      string
	GraderWorkers  int
	GraderQueue    int
	TraceBatch     int
	TraceFlushMS   int
	MediaDir       string
}

func Load() Config {
	return Config{
		Env:           getenv("APP_ENV", "development"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		MySQLDSN:      getenv("MYSQL_DSN", "miniqxt:miniqxt_pass@tcp(127.0.0.1:28391)/miniqxt?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"),
		JWTSecret:     getenv("JWT_SECRET", "miniqxt-dev-secret-change-me"),
		GraderWorkers: getenvInt("GRADER_WORKERS", 8),
		GraderQueue:   getenvInt("GRADER_QUEUE_SIZE", 4096),
		TraceBatch:    getenvInt("TRACE_BATCH_SIZE", 50),
		TraceFlushMS:  getenvInt("TRACE_FLUSH_MS", 200),
		MediaDir:      getenv("MEDIA_DIR", "/app/media"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}
