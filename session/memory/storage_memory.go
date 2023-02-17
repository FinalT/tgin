package memory

import (
	"context"
	"errors"
	"github.com/FinalT/tgin/session"
	cache "github.com/patrickmn/go-cache"
	"time"
)

type Store struct {
	// 利用一个内存缓存来帮助我们管理过期时间
	c          *cache.Cache
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
// 实际上，这里也可以考虑使用 Option 设计模式，允许用户控制过期检查的间隔
func NewStore(expiration time.Duration) *Store {
	return &Store{
		c: cache.New(expiration, time.Second),
	}
}

func (m *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}
	m.c.Set(sess.ID(), sess, m.expiration)
	return sess, nil
}

func (m *Store) Refresh(ctx context.Context, id string) error {
	sess, err := m.Get(ctx, id)
	if err != nil {
		return nil
	}
	m.c.Set(sess.ID(), sess, m.expiration)
	return nil
}

func (m *Store) Remove(ctx context.Context, id string) error {
	m.c.Delete(id)
	return nil
}

func (m *Store) Get(ctx context.Context, id string) (session.Session, error) {
	sess, ok := m.c.Get(id)
	if !ok {
		return nil, errors.New("session not found")
	}
	return sess.(*memorySession), nil
}

type memorySession struct {
	id         string
	data       map[string]string
	expiration time.Duration
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("找不到这个 key")
	}
	return val, nil
}

func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	m.data[key] = val
	return nil
}

func (m *memorySession) ID() string {
	return m.id
}
