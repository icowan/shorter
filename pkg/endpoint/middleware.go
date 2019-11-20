/**
 * @Time : 19/11/2019 10:39 AM
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
	"time"
)

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				_ = logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

var ErrLimitExceed = errors.New("Rate limit exceed!")

func NewTokenBucketLimitter(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
