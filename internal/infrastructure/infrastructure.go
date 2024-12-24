package infrastructure

import (
	"net/http"

	handler "github.com/isurucuma/rate-limiting/internal/server"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func NewRouter(httpHandler *handler.HTTPHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/get", httpHandler.HandleGet)
	return mux
}
