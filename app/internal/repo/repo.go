package repo

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"urlshort/app/internal/repo/models"
	"urlshort/app/pkg"

	"github.com/google/uuid"
)

type UrlStore struct {
	urls map[string]string
	mu   sync.RWMutex
	file *os.File
}

func NewUrlStore(filename string) *UrlStore {
	store := &UrlStore{
		urls: make(map[string]string),
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("UrlStore", err)
	}

	store.file = f

	if err := store.load(); err != nil {
		log.Println("Error loading data in UrlStore", err)
	}

	return store
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
			if err := s.save(key, longUrl); err != nil {
				log.Println("Error saving to UrlStore", err)
			}
			return key
		}
	}
}

func (s *UrlStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.urls)
}

func (s *UrlStore) save(key, url string) error {
	e := gob.NewEncoder(s.file)
	urlRecord := models.UrlRecord{Key: key, Url: url, BaseModel: models.BaseModel{
		Identifier: uuid.UUID().String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}}
	return e.Encode(urlRecord)
}

func (s *UrlStore) load() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	d := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var r models.UrlRecord

		if err = d.Decode(&r); err != nil {
			s.Set(r.Key, r.Url)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}
