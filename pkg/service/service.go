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
	"strings"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

type Service interface {
	Get(ctx context.Context, code string) (redirect *Redirect, err error)
	Post(ctx context.Context, domain string) (redirect *Redirect, err error)
}

type service struct {
	repository Repository
	logger     log.Logger
	shortUri   string
}

func New(middleware []Middleware, logger log.Logger, repository Repository, shortUri string) Service {
	var svc Service = NewService(logger, repository, shortUri)
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}

func NewService(logger log.Logger, repository Repository, shortUri string) Service {
	return &service{repository: repository, shortUri: shortUri, logger: logger}
}

func (s *service) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	return s.repository.Find(code)
}

func (s *service) Post(ctx context.Context, domain string) (redirect *Redirect, err error) {
	now := time.Now()
	local, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		now = now.In(local)
	} else {
		_ = level.Warn(s.logger).Log("time", "LoadLocation", "err", err)
	}

	code := shortid.MustGenerate()

	redirect = &Redirect{
		Code:      code,
		URL:       domain,
		CreatedAt: now,
	}

	if err = s.repository.Store(redirect); err != nil {
		return
	}

	redirect.URL = strings.TrimRight(s.shortUri, "/") + "/" + code
	return
}
