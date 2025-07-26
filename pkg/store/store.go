package store

import "errors"

type Store struct {
	m map[string][]byte
}

type IStore interface {
	Get(key string) (string, error)
	Set(key, val string)
}

func NewStore() *Store {
	m := make(map[string][]byte)
	return &Store{m}
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

func (s *Store) Set(key string, val []byte) {
	s.m[key] = val
}
