/**
 * @Time : 2019-11-15 15:32
 * @Author : solacowa@gmail.com
 * @File : repository
 * @Software: GoLand
 */

package redis

import (
	"fmt"
	"github.com/icowan/shorter/pkg/service"
	"github.com/pkg/errors"
	"time"
)

type redisRepository struct {
	client RedisInterface
}

func NewRedisRepository(drive RedisDrive, hosts, password, prefix string, database int) (service.Repository, error) {
	rdsClient := NewRedisClient(drive, hosts, password, prefix, database)

	return &redisRepository{client: rdsClient}, nil
}

func (m *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (m *redisRepository) Find(code string) (redirect *service.Redirect, err error) {
	data, err := m.client.HGetAll(m.generateKey(code))
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(service.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	now, err := time.Parse("2006-01-02 15:04:05", data["created_at"])
	if err != nil {
		return
	}

	return &service.Redirect{
		Code:      data["code"],
		URL:       data["url"],
		CreatedAt: now.In(time.Local),
	}, nil
}

func (m *redisRepository) Store(redirect *service.Redirect) error {
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	err := m.client.HMSet(m.generateKey(redirect.Code), data)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
