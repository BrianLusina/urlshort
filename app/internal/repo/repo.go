package repo

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"urlshort/app/internal/repo/models"
	"urlshort/app/pkg"

	"github.com/google/uuid"
)

const saveQueueLength = 1000

type UrlStore struct {
	urls     map[string]string
	mu       sync.RWMutex
	file     *os.File
	saveChan chan models.UrlRecord
}

func NewUrlStore(filename string) *UrlStore {
	saveChan := make(chan models.UrlRecord, saveQueueLength)
	store := &UrlStore{
		urls:     make(map[string]string),
		saveChan: saveChan,
	}

	if err := store.load(); err != nil {
		log.Println("Error loading data in UrlStore", err)
	}

	go store.saveLoop(filename)

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
			s.saveChan <- models.UrlRecord{Key: key, Url: longUrl, BaseModel: models.BaseModel{
				Identifier: uuid.New().String(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}}
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
		Identifier: uuid.New().String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}}
	return e.Encode(urlRecord)
}

func (s *UrlStore) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("UrlStore", err)
	}
	defer f.Close()

	e := json.NewEncoder(f)

	for {
		record := <-s.saveChan
		if err := e.Encode(record); err != nil {
			log.Println("Error saving record", err)
		}
	}
}

func (s *UrlStore) load() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	d := json.NewDecoder(s.file)
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
