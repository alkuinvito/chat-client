package store

import "errors"

type Store struct {
	m map[string][]byte
}

type IStore interface {
	Clear()
	Delete(key string)
	Get(key string) (string, error)
	GetString(key string) (string, error)
	Keys() []string
	Set(key, val string)
}

func NewStore() *Store {
	m := make(map[string][]byte)
	return &Store{m}
}

func (s *Store) Clear() {
	keys := s.Keys()
	for _, k := range keys {
		s.Delete(k)
	}
}

func (s *Store) Delete(key string) {
	val, ok := s.m[key]
	if !ok {
		return
	}

	// zero-ing unused value
	for i := range val {
		val[i] = 0
	}

	// delete key
	delete(s.m, key)
}

func (s *Store) Get(key string) ([]byte, error) {
	val, ok := s.m[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	return val, nil
}

func (s *Store) GetString(key string) (string, error) {
	val, ok := s.m[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return string(val), nil
}

func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}

	return keys
}

func (s *Store) Set(key string, val []byte) {
	s.m[key] = val
}
