package compressor

import (
	"fmt"
	"io"
)

type FrequencyTable struct {
	Table map[byte]int
}

func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{
		Table: make(map[byte]int),
	}
}

func (f *FrequencyTable) Create(reader io.Reader) error {
	buf := make([]byte, 8192)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		// Process each byte directly
		for i := 0; i < n; i++ {
			f.Add(buf[i])
		}
	}

	return nil
}

func (f *FrequencyTable) Add(b byte) {
	f.Table[b]++
}

func (f *FrequencyTable) Get() map[byte]int {
	return f.Table
}

func (f *FrequencyTable) GetHumanReadable() string {
	var result string
	for k, v := range f.Table {
		result += fmt.Sprintf("%02x %d ", k, v) + "\n"
	}
	return result
}
