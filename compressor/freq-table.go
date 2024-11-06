package compressor

import "fmt"

type FrequencyTable struct {
	Table map[rune]int
}

func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{
		Table: make(map[rune]int),
	}
}

func (f *FrequencyTable) Create(text string) {
	for _, char := range text {
		f.Add(char)
	}
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
