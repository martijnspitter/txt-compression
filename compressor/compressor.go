package compressor

type Compressor struct {
	BinaryTree  *BinaryTree
	FreqTable   *FrequencyTable
	orignalText string
}

func NewCompressor(orignalText string) *Compressor {
	newFreqTable := NewFrequencyTable()
	newFreqTable.Create(orignalText)
	freqTable := newFreqTable.Get()
	binaryTree := NewBinaryTree(&freqTable)
	binaryTree.GetPrefixCodeTable()

	return &Compressor{
		FreqTable:   newFreqTable,
		BinaryTree:  binaryTree,
		orignalText: orignalText,
	}
}

func (c *Compressor) GetHeader() string {
	return c.BinaryTree.GetCodeTableAsString()
}

func (c *Compressor) GetCompressedText() string {
	return c.BinaryTree.GetCompressedText(c.orignalText)
}
