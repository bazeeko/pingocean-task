package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var ErrKeyAlreadyExists = errors.New("key alraedy exists in map")
var ErrKeyDoesNotExist = errors.New("key doesn't exist in map")

type Storage struct {
	Map map[string]int
	sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{Map: make(map[string]int)}
}

func (s *Storage) Set(url string, count int) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Map[url]; ok {
		return fmt.Errorf("Set: %w", ErrKeyAlreadyExists)
	}

	s.Map[url] = count

	return nil
}

func (s *Storage) Get(url string) (int, error) {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.Map[url]; !ok {
		return 0, fmt.Errorf("Get: %w", ErrKeyDoesNotExist)
	}

	return s.Map[url], nil
}

func (s *Storage) GetAll() map[string]int {
	s.RLock()
	defer s.RUnlock()

	return s.Map
}

func countStringInURL(str string, url string, ctx context.Context) (int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("countStringInURL: %w", err)
	}

	ctxReq, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req = req.WithContext(ctxReq)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 0, fmt.Errorf("countStringInURL: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("countStringInURL: %w", err)
	}

	count := strings.Count(string(data), str)

	return count, nil
}

func Search(str string, urls []string, ctx context.Context) map[string]int {
	storage := NewStorage()

	// wg := &sync.WaitGroup{}
	limiter := NewLimiter(MAX_GOROUTINES)

LOOP:
	for _, url := range urls {
		select {
		case <-ctx.Done():
			break LOOP
		default:
			limiter.Add()

			go func(url string) {
				defer limiter.Done()

				count, err := countStringInURL(str, url, ctx)
				if err != nil {
					log.Printf("Search: %s\n", err)
					storage.Set(url, -1)

					return
				}

				storage.Set(url, count)
			}(url)
		}
	}

	limiter.Wait()

	return storage.GetAll()
}
