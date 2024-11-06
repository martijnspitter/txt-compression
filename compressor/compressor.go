package compressor

type Compressor struct {
	root *Node
}

func NewCompressor(root *Node) *Compressor {
	return &Compressor{
		root: root,
	}
}
