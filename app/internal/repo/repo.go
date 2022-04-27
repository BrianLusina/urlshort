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
	"urlshort/app/pkg/errors"

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
	store := &UrlStore{urls: make(map[string]string)}

	if filename != "" {
		saveChan := make(chan models.UrlRecord, saveQueueLength)
		store.saveChan = saveChan

		if err := store.load(filename); err != nil {
			log.Println("Error loading data in UrlStore", err)
		}

		go store.saveLoop(filename)
	}

	return store
}

func (s *UrlStore) Get(key, url *string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if u, ok := s.urls[*key]; ok {
		*url = u
		return nil
	}

	return errors.ErrUrlNotFound
}

func (s *UrlStore) Set(key, url *string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.urls[*key]; ok {
		return errors.ErrUrlExists
	}

	s.urls[*key] = *url
	return nil
}

// Put generates a short key from a long url and returns it
func (s *UrlStore) Put(longUrl, key *string) error {
	for {
		*key = pkg.GenerateKey(s.Count())
		if err := s.Set(key, longUrl); err == nil {
			break
		}
	}

	if s.saveChan != nil {
		s.saveChan <- models.UrlRecord{Key: *key, Url: *longUrl, BaseModel: models.BaseModel{
			Identifier: uuid.New().String(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}}
	}
	return nil
}

func (s *UrlStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

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
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening UrlStore: ", err)
	}

	e := json.NewEncoder(f)

	for {
		record := <-s.saveChan
		if err := e.Encode(record); err != nil {
			log.Println("Error saving record: ", err)
		}
	}
}

func (s *UrlStore) load(filename string) error {
	f, err := os.Open(filename)

	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	for err == nil {
		var r models.UrlRecord

		if err = d.Decode(&r); err != nil {
			s.Set(&r.Key, &r.Url)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}
