package reader

import (
	"fmt"
	"os"
)

type Reader struct {
	path string
}

func NewReader(path string) *Reader {
	return &Reader{
		path,
	}
}

func (r *Reader) ReadFile() (string, error) {
	file, err := os.Open(r.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	size := stat.Size()
	if size == 0 {
		return "", nil
	}

	if stat.IsDir() {
		return "", fmt.Errorf("'%s' is a directory", r.path)
	}

	content, err := os.ReadFile(r.path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
