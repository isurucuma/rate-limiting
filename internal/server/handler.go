package handler

import (
	"context"
	"net/http"

	"github.com/isurucuma/rate-limiting/internal/usecase"
)

type HTTPHandler struct {
	RateLimiterUsecase usecase.RateLimiterUsecase
}

func NewHTTPHandler(usecase usecase.RateLimiterUsecase) *HTTPHandler {
	return &HTTPHandler{RateLimiterUsecase: usecase}
}

func (h *HTTPHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if h.RateLimiterUsecase.HandleRequest(ctx) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request processed successfully."))
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many requests. Please try again later."))
	}
}
