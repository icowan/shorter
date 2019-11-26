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
	now = now.In(time.Local)
	code := shortid.MustGenerate()

	// todo 考虑如何处理垃圾数据的问题 得复的url 不同的code

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
