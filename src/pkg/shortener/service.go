package shortener

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
	"time"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

//Service
type Service interface {
	Find(ctx context.Context, code string) (r *Redirect, err error)
	Store(ctx context.Context, redirect *Redirect) error
}

type service struct {
	repository Repository
	logger     log.Logger
}

func NewService(repository Repository, logger log.Logger) Service {
	return &service{repository: repository, logger: logger}
}

func (s *service) Find(ctx context.Context, code string) (r *Redirect, err error) {
	return s.repository.Find(code)
}

func (s *service) Store(ctx context.Context, redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errors.Wrap(ErrRedirectInvalid, err.Error())
	}
	now := time.Now()
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		_ = level.Warn(s.logger).Log("time", "LoadLocation", "err", err.Error())
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = now.In(local)
	return s.repository.Store(redirect)
}
