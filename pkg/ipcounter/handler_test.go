package ipcounter_test

import (
	"errors"
	"io"
	"os"
	"solution/pkg/ipcounter"
	"testing"
	"testing/fstest"
)

type reader struct {
	read func(p []byte) (int, error)
}

func (r *reader) Read(p []byte) (int, error) {
	return r.read(p)
}

func TestFileHandler_OpenRead(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		data := "192.168.0.1\n192.168.0.2"
		file := "testfile.txt"
		f := fstest.MapFS{
			file: &fstest.MapFile{Data: []byte(data)},
		}

		h := ipcounter.FileHandler{FS: f}
		reader, err := h.OpenRead(file)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		content, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read from reader: %v", err)
		}

		if string(content) != data {
			t.Fatalf("expected content %q, got %q", data, string(content))
		}
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		f := fstest.MapFS{}
		h := ipcounter.FileHandler{FS: f}

		_, err := h.OpenRead("nonexistent.txt")
		if !errors.Is(err, ipcounter.ErrReadInputFile) {
			t.Fatalf("expected error of type *fs.PathError, got %v", err)
		}
	})
}

func TestFileHandler_CreateWrite(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		tmp := t.TempDir()
		filePath := tmp + "/output.txt"
		data := "192.168.1.1\n192.168.1.2"

		var h ipcounter.FileHandler

		_, err := h.CreateWrite(filePath, data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		if string(content) != data {
			t.Fatalf("want %q, got %q", data, string(content))
		}
	})

	t.Run("failure - invalid filename", func(t *testing.T) {
		t.Parallel()

		var h ipcounter.FileHandler

		_, err := h.CreateWrite("", "data")
		if !errors.Is(err, ipcounter.ErrCreateFile) {
			t.Fatal("expected an error when filename is invalid, got nil")
		}
	})
}
