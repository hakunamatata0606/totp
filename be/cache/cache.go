package cache

import (
	"errors"
	"sync"
	"time"
)

type CacheIf interface {
	Set(string, interface{})
	Get(string) (interface{}, error)
}

type inMemRecord struct {
	exp  time.Time
	data interface{}
}

type inMemCacheImpl struct {
	m       sync.Mutex
	store   map[string]*inMemRecord
	timeout time.Duration
}

func NewInMemCache(timeout time.Duration) CacheIf {
	return &inMemCacheImpl{
		m:       sync.Mutex{},
		store:   make(map[string]*inMemRecord),
		timeout: timeout,
	}
}

func (c *inMemCacheImpl) Set(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()
	c.store[key] = &inMemRecord{
		exp:  time.Now().Add(c.timeout),
		data: value,
	}
}

func (c *inMemCacheImpl) Get(key string) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()
	record, ok := c.store[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	t := time.Now()

	if t.After(record.exp) {
		delete(c.store, key)
		return nil, errors.New("key expired")
	}
	return record.data, nil
}
