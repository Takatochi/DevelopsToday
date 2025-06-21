package services

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

// breedAPIURL is the URL for the cat breed API.
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
	mutex      sync.RWMutex
	apiURL     string // Make URL configurable
}

func NewBreed() Validator {
	return &serviceBreed{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		ttl:        10 * time.Minute,
		apiURL:     breedAPIURL, // Use global as default
	}
}

// NewBreedWithURL creates a breed validator with custom URL (for testing)
func NewBreedWithURL(apiURL string) Validator {
	return &serviceBreed{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		ttl:        10 * time.Minute,
		apiURL:     apiURL,
	}
}

func (s *serviceBreed) IsValid(breedName string) bool {
	s.mutex.RLock()

	// if the cache is fresh, we return the result quickly
	if time.Since(s.lastFetch) <= s.ttl && len(s.cache) > 0 {
		defer s.mutex.RUnlock()
		return s.searchInCache(breedName)
	}

	// if the cache is outdated, but there is data,
	//we return it from the old cache and update it asynchronously
	hasOldCache := len(s.cache) > 0
	result := false
	if hasOldCache {
		result = s.searchInCache(breedName)
	}
	s.mutex.RUnlock()

	// asynchronously update the cache without blocking the request
	s.mutex.RLock()
	needsUpdate := time.Since(s.lastFetch) > s.ttl
	s.mutex.RUnlock()

	if needsUpdate {
		go func() {
			s.mutex.Lock()
			defer s.mutex.Unlock()
			if time.Since(s.lastFetch) > s.ttl {
				_ = s.fetchBreeds() // ignore the error in background
			}
		}()
	}

	// if there is no old cache, we make a synchronous request
	if !hasOldCache {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if err := s.fetchBreeds(); err != nil {
			return false
		}
		return s.searchInCache(breedName)
	}

	return result
}

func (s *serviceBreed) searchInCache(breedName string) bool {
	for _, b := range s.cache {
		if strings.EqualFold(b.Name, breedName) {
			return true
		}
	}
	return false
}

func (s *serviceBreed) fetchBreeds() error {
	resp, err := s.httpClient.Get(s.apiURL) // Use instance URL instead of global
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var breeds []Breed
	if err = json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return err
	}

	s.cache = breeds
	s.lastFetch = time.Now()
	return nil
}
