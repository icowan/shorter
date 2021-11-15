/**
 * @Time : 2021/11/15 4:55 PM
 * @Author : solacowa@gmail.com
 * @File : logging
 * @Software: GoLand
 */

package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type logging struct {
	traceId string
	logger  log.Logger
	next    Service
}

func (s *logging) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Get", "code", code,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Get(ctx, code)
}

func (s *logging) Post(ctx context.Context, domain string) (redirect *Redirect, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			s.traceId, ctx.Value(s.traceId),
			"method", "Post", "domain", domain,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Post(ctx, domain)
}

func NewLogging(logger log.Logger, traceId string) Middleware {
	return func(next Service) Service {
		return &logging{traceId, logger, next}
	}
}
