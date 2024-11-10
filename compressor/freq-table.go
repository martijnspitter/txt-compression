package compressor

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type FrequencyTable struct {
	Table map[rune]int
}

func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{
		Table: make(map[rune]int),
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

		// Process the buffer as UTF-8 text
		str := string(buf[:n])
		for len(str) > 0 {
			r, size := utf8.DecodeRuneInString(str)
			if r == utf8.RuneError {
				return fmt.Errorf("invalid UTF-8 sequence")
			}
			f.Add(r)
			str = str[size:]
		}
	}
	return nil
}

func (f *FrequencyTable) Add(char rune) {
	f.Table[char]++
}

func (f *FrequencyTable) Get() map[rune]int {
	return f.Table
}

func (f *FrequencyTable) GetHumanReadable() string {
	var result string
	for k, v := range f.Table {
		result += fmt.Sprintf("%c %d ", k, v) + "\n"
	}
	return result
}
