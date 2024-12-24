package usecase

import (
	"context"

	"github.com/isurucuma/rate-limiting/internal/domain"
)

type RateLimiterUsecase interface {
	HandleRequest(ctx context.Context) bool
}

type rateLimiterUsecase struct {
	rateLimiter *domain.RateLimiter
}

func NewRateLimiterUsecase(rateLimiter *domain.RateLimiter) RateLimiterUsecase {
	return &rateLimiterUsecase{rateLimiter: rateLimiter}
}

func (u *rateLimiterUsecase) HandleRequest(ctx context.Context) bool {
	return u.rateLimiter.AllowRequest(ctx)
}
