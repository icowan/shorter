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
	//key := m.generateKey(code)
	//data, err := m.client.HGet("redirect", code)
	//if err != nil {
	//	return nil, errors.Wrap(err, "repository.Redirect.Find")
	//}
	//if len(data) == 0 {
	//	return nil, errors.Wrap(service.ErrRedirectNotFound, "repository.Redirect.Find")
	//}

	return
}

func (m *redisRepository) Store(redirect *service.Redirect) error {
	return nil
}
