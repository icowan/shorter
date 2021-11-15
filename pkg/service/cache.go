/**
 * @Time : 2021/11/15 5:23 PM
 * @Author : solacowa@gmail.com
 * @File : cache
 * @Software: GoLand
 */

package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type lruCache struct {
	m          map[string]*LinkNode // 指向哈希表的指针
	cap        int                  // 长度
	head, tail *LinkNode            // 两个哨兵
}

type LinkNode struct {
	key       string
	val       []byte
	pre, next *LinkNode
}

func NewLruCache(capacity int) *lruCache {
	head := &LinkNode{"", nil, nil, nil}
	tail := &LinkNode{"", nil, nil, nil}
	head.next = tail
	tail.pre = head
	return &lruCache{make(map[string]*LinkNode), capacity, head, tail}
}

func (s *lruCache) Get(ctx context.Context, key string) (res []byte, err error) {
	if node, exist := s.m[key]; exist {
		s.moveToHead(node)
		return node.val, nil
	}
	return nil, errors.New("key notfound")
}

func (s *lruCache) moveToHead(node *LinkNode) {
	head := s.head
	// 从当前位置删除节点
	node.pre.next = node.next
	node.next.pre = node.pre
	// 将节点插入头部
	node.next = head.next
	head.next.pre = node
	node.pre = head
	head.next = node
}

func (s *lruCache) Put(ctx context.Context, key string, value []byte) {
	head := s.head
	tail := s.tail
	if node, exist := s.m[key]; exist { // 如果已经存在对应的节点，要将对应的元素交换到头部；并更新值（值可能变了）
		node.val = value
		s.moveToHead(node)
	} else { // 节点不存在
		// 创建节点
		node := &LinkNode{key, value, nil, nil}
		// 判断缓存容量是否已经满了
		if len(s.m) == s.cap {
			// 删除最后的元素
			delete(s.m, tail.pre.key)
			tail.pre.pre.next = tail
			tail.pre = tail.pre.pre
		}
		// 将节点放到头部
		node.next = head.next
		node.pre = head
		head.next.pre = node
		head.next = node
		// 参入哈希表中
		s.m[key] = node
	}
}

type cache struct {
	lruCache *lruCache
	traceId  string
	next     Service
	pkgName  string
	logger   log.Logger
}

func (s *cache) Get(ctx context.Context, code string) (redirect *Redirect, err error) {
	get, err := s.lruCache.Get(ctx, code)
	if err == nil {
		err = json.Unmarshal(get, &redirect)
		if err == nil {
			return redirect, nil
		}
		_ = level.Error(s.logger).Log("json", "Unmarshal", "err", err.Error())
	}
	defer func() {
		b, _ := json.Marshal(redirect)
		s.lruCache.Put(ctx, code, b)
	}()
	_ = level.Warn(s.logger).Log("lruCache", "Get", "err", err.Error())
	return s.next.Get(ctx, code)
}

func (s *cache) Post(ctx context.Context, domain string) (redirect *Redirect, err error) {
	return s.next.Post(ctx, domain)
}

func NewCache(logger log.Logger, traceId string, cacheCap int) Middleware {
	logger = log.With(logger, "service", "cache")
	return func(next Service) Service {
		return &cache{
			logger:   logger,
			next:     next,
			traceId:  traceId,
			pkgName:  "service",
			lruCache: NewLruCache(cacheCap),
		}
	}
}
