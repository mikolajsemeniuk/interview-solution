package ipcounter

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Storage provides methods to list and increment IP count records.
type Storage interface {
	List(ctx context.Context, namespace, set string) ([]Record, error)
	Increment(ctx context.Context, namespace, set, key string, count int) error
}

// Service performs import and export operations between I/O and storage.
type Service struct {
	storage Storage
}

// NewService returns a new Service using the given Storage.
func NewService(s Storage) *Service {
	return &Service{storage: s}
}

// Import reads IPs from r, increments their counts in storage, and returns elapsed seconds.
func (s *Service) Import(ctx context.Context, r io.Reader, namespace, set, mode string) (int, error) {
	start := time.Now()
	var eg errgroup.Group
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		ip := scanner.Text()
		if ip == "" {
			continue
		}

		if mode == "async" {
			eg.Go(func() error { return s.storage.Increment(ctx, namespace, set, ip, 1) })
			continue
		}

		if err := s.storage.Increment(ctx, namespace, set, ip, 1); err != nil {
			return 0, ErrSetKey
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, ErrProcessInputFile
	}

	if err := eg.Wait(); err != nil {
		return 0, ErrProcessRecords
	}

	return int(time.Since(start).Seconds()), nil
}

// Export retrieves IP counts from storage, formats them, and returns the result string.
func (s *Service) Export(ctx context.Context, namespace, set, mode string) (string, error) {
	items, err := s.storage.List(ctx, namespace, set)
	if err != nil {
		return "", ErrListRecords
	}

	var sb strings.Builder
	var mu sync.Mutex
	var eg errgroup.Group
	for _, item := range items {
		if mode == "sync" {
			line := fmt.Sprintf("%s, count=%d\n", item.Key, item.Count)
			if _, err := sb.WriteString(line); err != nil {
				return "", ErrWriteToStringBuilder
			}

			continue
		}

		eg.Go(func() error {
			mu.Lock()
			defer mu.Unlock()
			_, err := sb.WriteString(fmt.Sprintf("%s, count=%d\n", item.Key, item.Count))
			if err != nil {
				return errors.Join(ErrWriteToStringBuilder, err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return "", ErrProcessRecords
	}

	return sb.String(), nil
}
