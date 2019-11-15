/**
 * @Time : 2019-11-15 15:32
 * @Author : solacowa@gmail.com
 * @File : repository
 * @Software: GoLand
 */

package redis

import (
	"fmt"
	"github.com/icowan/shorter/src/pkg/shortener"
	"github.com/pkg/errors"
	"strconv"
)

type redisRepository struct {
	client RedisInterface
}

func NewRedisRepository(drive RedisDrive, hosts, password, prefix string, database int) (shortener.Repository, error) {
	rdsClient := NewRedisClient(drive, hosts, password, prefix, database)

	return &redisRepository{client: rdsClient}, nil
}

func (m *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (m *redisRepository) Find(code string) (redirect *shortener.Redirect, err error) {
	key := m.generateKey(code)
	data, err := m.client.HGet("redirect", code)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}

}

func (m *redisRepository) Store(redirect *shortener.Redirect) error {

	return
}
