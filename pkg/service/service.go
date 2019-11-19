/**
 * @Time : 19/11/2019 10:12 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("not found")
)

type Service interface {
	Get(ctx context.Context, code string) (redirect *Redirect, err error)
	Post(ctx context.Context, domain string) error
}

type service struct {
	repository Repository
	logger     log.Logger
}

func New(middleware []Middleware, repository Repository) Service {
	var svc Service = NewService(repository)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	return s.repository.Find(code)
}

func (s *service) Post(ctx context.Context, domain string) error {
	now := time.Now()
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		_ = level.Warn(s.logger).Log("time", "LoadLocation", "err", err.Error())
	}
	return s.repository.Store(&Redirect{
		Code:      shortid.MustGenerate(),
		URL:       domain,
		CreatedAt: now.In(local),
	})
}
