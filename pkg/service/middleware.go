/**
 * @Time : 19/11/2019 10:19 AM
 * @Author : solacowa@gmail.com
 * @File : middleware
 * @Software: GoLand
 */

package service

import (
	"context"
	"github.com/go-kit/kit/log"
)

type Middleware func(svc Service) Service

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{logger, next}
	}
}

func (l loggingMiddleware) Get(ctx context.Context, code string) (rs *Redirect, err error) {
	defer func() {
		_ = l.logger.Log("method", "Get", "s", code, "err", err)
	}()
	return l.next.Get(ctx, code)
}

func (l loggingMiddleware) Post(ctx context.Context, uri string) (err error) {
	defer func() {
		_ = l.logger.Log("method", "Post", "uri", uri, "err", err)
	}()
	return l.next.Post(ctx, uri)
}
