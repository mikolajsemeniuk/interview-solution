package ipcounter_test

import (
	"context"
	"errors"
	"solution/pkg/ipcounter"
	"strings"
	"testing"
)

type storage struct {
	list      func(ctx context.Context, namespace, set string) ([]ipcounter.Record, error)
	increment func(ctx context.Context, namespace, set, key string, count int) error
}

func (s *storage) List(ctx context.Context, namespace, set string) ([]ipcounter.Record, error) {
	return s.list(ctx, namespace, set)
}

func (s *storage) Increment(ctx context.Context, namespace, set, key string, count int) error {
	return s.increment(ctx, namespace, set, key, count)
}

//nolint:funlen
func TestService_Import(t *testing.T) {
	t.Parallel()

	t.Run("success async", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			increment: func(_ context.Context, _, _, _ string, _ int) error {
				return nil
			},
		}

		data := "192.168.0.1\n192.168.0.2\n"
		reader := strings.NewReader(data)
		service := ipcounter.NewService(store)

		_, err := service.Import(context.Background(), reader, "namespace", "set", "async")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("success sync", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			increment: func(_ context.Context, _, _, _ string, _ int) error {
				return nil
			},
		}

		data := "10.0.0.1\n10.0.0.2\n"
		reader := strings.NewReader(data)
		service := ipcounter.NewService(store)

		_, err := service.Import(context.Background(), reader, "namespace", "set", "sync")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("failure sync", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			increment: func(_ context.Context, _, _, _ string, _ int) error {
				return ipcounter.ErrSetKey
			},
		}

		data := "127.0.0.1\n"
		reader := strings.NewReader(data)
		service := ipcounter.NewService(store)

		_, err := service.Import(context.Background(), reader, "namespace", "set", "sync")
		if !errors.Is(err, ipcounter.ErrSetKey) {
			t.Fatalf("expected error %v, got %v", ipcounter.ErrSetKey, err)
		}
	})

	t.Run("failure async", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			increment: func(_ context.Context, _, _, _ string, _ int) error {
				return ipcounter.ErrProcessRecords
			},
		}

		data := "127.0.0.1\n"
		reader := strings.NewReader(data)
		service := ipcounter.NewService(store)

		_, err := service.Import(context.Background(), reader, "namespace", "set", "async")
		if !errors.Is(err, ipcounter.ErrProcessRecords) {
			t.Fatalf("expected error %v, got %v", ipcounter.ErrProcessRecords, err)
		}
	})

	t.Run("failure reading input", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			increment: func(_ context.Context, _, _, _ string, _ int) error {
				return nil
			},
		}

		reader := &reader{
			read: func(_ []byte) (int, error) { return 0, ipcounter.ErrProcessInputFile },
		}

		service := ipcounter.NewService(store)

		_, err := service.Import(context.Background(), reader, "namespace", "set", "sync")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ipcounter.ErrProcessInputFile) {
			t.Fatalf("expected error %v, got %v", ipcounter.ErrProcessInputFile, err)
		}
	})
}

func TestService_Export(t *testing.T) {
	t.Parallel()

	records := []ipcounter.Record{
		{"192.168.1.1", 5},
		{"192.168.1.2", 3},
	}

	t.Run("success sync", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			list: func(_ context.Context, _, _ string) ([]ipcounter.Record, error) {
				return records, nil
			},
		}

		service := ipcounter.NewService(store)
		result, err := service.Export(context.Background(), "namespace", "set", "sync")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := "192.168.1.1, count=5\n192.168.1.2, count=3\n"
		if result != expected {
			t.Fatalf("expected %q, got %q", expected, result)
		}
	})

	t.Run("success async", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			list: func(_ context.Context, _, _ string) ([]ipcounter.Record, error) {
				return records, nil
			},
		}

		service := ipcounter.NewService(store)
		result, err := service.Export(context.Background(), "namespace", "set", "async")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := []string{
			"192.168.1.1, count=5\n",
			"192.168.1.2, count=3\n",
		}

		for _, line := range expected {
			if !strings.Contains(result, line) {
				t.Fatalf("expected result to contain %q, got %q", line, result)
			}
		}
	})

	t.Run("failure listing records", func(t *testing.T) {
		t.Parallel()

		store := &storage{
			list: func(_ context.Context, _, _ string) ([]ipcounter.Record, error) {
				return nil, ipcounter.ErrListRecords
			},
		}

		service := ipcounter.NewService(store)
		_, err := service.Export(context.Background(), "namespace", "set", "sync")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ipcounter.ErrListRecords) {
			t.Fatalf("expected error %v, got %v", ipcounter.ErrListRecords, err)
		}
	})
}
