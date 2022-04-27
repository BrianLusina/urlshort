package repo

import (
	"sync"
	"urlshort/app/pkg"
)

type UrlStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewUrlStore() *UrlStore {
	return &UrlStore{
		urls: make(map[string]string),
	}
}

func (s *UrlStore) Get(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.urls[key]
}

func (s *UrlStore) Set(key, url string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.urls[key]
	if ok {
		return false
	}

	s.urls[key] = url
	return true
}

// Put generates a short key from a long url and returns it
func (s *UrlStore) Put(longUrl string) string {
	for {
		key := pkg.GenerateKey(s.Count())
		if s.Set(key, longUrl) {
			return key
		}
	}
}

func (s *UrlStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.urls)
}
