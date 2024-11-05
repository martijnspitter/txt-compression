package compressor

type Compressor struct {
	text string
}

func NewCompressor(text string) *Compressor {
	return &Compressor{
		text,
	}
}

func (c *Compressor) GenerateFreqTable() map[rune]int {
	return make(map[rune]int, 0)
}
