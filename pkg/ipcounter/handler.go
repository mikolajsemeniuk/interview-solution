package ipcounter

import (
	"io"
	"io/fs"
	"os"
)

// FileHandler wraps an fs.FS for reading and the OS filesystem for writing.
type FileHandler struct {
	FS fs.FS
}

// OpenRead returns an io.ReadCloser so the caller can defer Close().
func (h *FileHandler) OpenRead(name string) (io.ReadCloser, error) {
	f, err := h.FS.Open(name)
	if err != nil {
		return nil, ErrReadInputFile
	}

	return f, nil
}

// CreateWrite returns an io.WriteCloser so the caller can defer Close().
func (h *FileHandler) CreateWrite(filename, data string) (io.WriteCloser, error) {
	out, err := os.Create(filename)
	if err != nil {
		return nil, ErrCreateFile
	}

	if _, err := out.WriteString(data); err != nil {
		return nil, ErrWriteFile
	}

	return out, nil
}
