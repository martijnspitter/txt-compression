package compressor

import "fmt"

type Compressor struct {
	text      string
	freqTable *map[rune]int
}

func NewCompressor(text string) *Compressor {
	return &Compressor{
		text:      text,
		freqTable: &map[rune]int{},
	}
}

func (c *Compressor) GenerateFreqTable() {
	freqTable := make(map[rune]int)
	for _, char := range c.text {
		freqTable[char]++
	}
	c.freqTable = &freqTable
}

func (c *Compressor) GetFreqTable() *map[rune]int {
	return c.freqTable
}

func (c *Compressor) GetHumanReadableFreqTable() string {
	var result string
	for k, v := range *c.freqTable {
		result += fmt.Sprintf("%c %d ", k, v) + "\n"
	}
	return result
}
