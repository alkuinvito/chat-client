package store

import (
	"context"
	"sync"
	"time"
)

type Store struct {
	sync.Mutex
	ctx         context.Context
	m           map[string][]byte
	expirations map[string]context.CancelFunc
}

type IStore interface {
	Clear()
	Delete(key string)
	Get(key string) []byte
	GetString(key string) string
	Keys() []string
	Set(key, val string)
	SetEx(key string, val []byte, expiredIn time.Duration)
	Startup(ctx context.Context)
}

func NewStore() *Store {
	m := make(map[string][]byte)
	expirations := make(map[string]context.CancelFunc)
	return &Store{m: m, expirations: expirations}
}

func (s *Store) Clear() {
	keys := s.Keys()
	for _, k := range keys {
		s.Delete(k)
	}
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()

	// cancel on-going timer and delete cancel func
	if cancel, ok := s.expirations[key]; ok {
		cancel()
		delete(s.expirations, key)
	}

	// zeroing unused value
	if val, ok := s.m[key]; ok {
		for i := range val {
			val[i] = 0
		}

		// delete key
		delete(s.m, key)
	}
}

func (s *Store) Get(key string) []byte {
	s.Lock()
	defer s.Unlock()

	val, ok := s.m[key]
	if !ok {
		return nil
	}

	return val
}

func (s *Store) GetString(key string) string {
	val, ok := s.m[key]
	if !ok {
		return ""
	}

	return string(val)
}

func (s *Store) Keys() []string {
	s.Lock()
	defer s.Unlock()

	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

func (s *Store) Set(key string, val []byte) {
	s.Lock()
	defer s.Unlock()

	if cancel, ok := s.expirations[key]; ok {
		cancel()
		delete(s.expirations, key)
	}

	s.m[key] = val
}

func (s *Store) SetEx(key string, val []byte, expiredIn time.Duration) {
	s.Lock()

	if cancel, ok := s.expirations[key]; ok {
		cancel()
		delete(s.expirations, key)
	}

	ctx, cancel := context.WithCancel(s.ctx)
	s.expirations[key] = cancel
	s.m[key] = append([]byte(nil), val...)

	s.Unlock()

	go func() {
		select {
		case <-time.After(expiredIn):
			s.Delete(key)
		case <-ctx.Done():
			// canceled
		}
	}()
}

func (s *Store) Startup(ctx context.Context) {
	s.ctx = ctx
}
