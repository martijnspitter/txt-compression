package compressor

import (
	"fmt"
	"sort"
)

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

type Node struct {
	Left  *Node
	Right *Node
	Value rune
	Count int
}

type Queue []*Node

func (q Queue) Len() int {
	return len(q)
}
func (q Queue) Less(i, j int) bool {
	return q[i].Count < q[j].Count
}
func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (c *Compressor) CreateBinaryTree() {
	queue := make(Queue, len(*c.freqTable))
	for k, v := range *c.freqTable {
		queue = append(queue, getNewNode(k, v))
	}
	sort.Sort(queue)

}

func getNewNode(value rune, count int) *Node {
	return &Node{
		Value: value,
		Count: count,
	}
}
