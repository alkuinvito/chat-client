package store

import "errors"

type Store struct {
	m map[string]string
}

func NewStore() *Store {
	m := make(map[string]string)
	return &Store{m}
}

func (s *Store) Get(key string) (string, error) {
	val, ok := s.m[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return val, nil
}

func (s *Store) Set(key, val string) {
	s.m[key] = val
}
