package services

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

var breedAPIURL = "https://api.thecatapi.com/v1/breeds"

type Breed struct {
	Name string `json:"name"`
}
type Validator interface {
	IsValid(breedName string) bool
}
type serviceBreed struct {
	lastFetch  time.Time
	httpClient *http.Client
	cache      []Breed
	ttl        time.Duration
	mutex      sync.Mutex
}

func NewBreed() Validator {
	return &serviceBreed{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		ttl:        10 * time.Minute,
	}
}

func (s *serviceBreed) IsValid(breedName string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if time.Since(s.lastFetch) > s.ttl || len(s.cache) == 0 {
		if err := s.fetchBreeds(); err != nil {
			return false
		}
	}

	for _, b := range s.cache {
		if strings.EqualFold(b.Name, breedName) {
			return true
		}
	}
	return false
}

func (s *serviceBreed) fetchBreeds() error {
	resp, err := s.httpClient.Get(breedAPIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var breeds []Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return err
	}

	s.cache = breeds
	s.lastFetch = time.Now()
	return nil
}
