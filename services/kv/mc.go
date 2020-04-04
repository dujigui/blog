package kv

import (
	"errors"
	"sync"
	"time"
)

var (
	kv = MemoryKV{
		values: make(map[string]interface{}),
		timers: make(map[string]*time.Timer),
	}
)

func KeyValue() KV {
	return &kv
}

type KV interface {
	Set(k, v string, expired time.Duration) error
	Get(k string) (string, error)
}

// todo 临时使用，如果日后项目用到键值对的地方增多，可以考虑引入 redis
type MemoryKV struct {
	values map[string]interface{}
	timers map[string]*time.Timer
	lock   sync.Mutex
}

func (s *MemoryKV) Set(k, v string, expired time.Duration) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	tv, ok := s.timers[k]
	if ok {
		delete(s.timers, k)
		tv.Stop()
	}
	if expired != 0 {
		t := time.NewTimer(expired)
		s.timers[k] = t
		go func() {
			_, ok := <-t.C
			if ok {
				s.lock.Lock()
				delete(s.values, k)
				delete(s.timers, k)
				s.lock.Unlock()
			}
		}()
	}
	s.values[k] = v
	return nil
}

func (s *MemoryKV) Get(k string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	v, ok := s.values[k]
	if !ok {
		return "", nil
	}
	if sv, ok := v.(string); ok {
		return sv, nil
	}
	return "", errors.New("非 string 类型")
}
