package compressor

import (
	"fmt"
	"sort"
)

type Node struct {
	Left  *Node
	Right *Node
	Value rune
	Count int
	ID    int
}

type BinaryTree struct {
	root      *Node
	queue     *Queue
	CodeTable CodeTable
}

type CodeTable map[rune]string

func NewBinaryTree(freqTable *map[rune]int) *BinaryTree {
	queue := createQueue(freqTable)
	root := createBinaryTree(queue)
	return &BinaryTree{
		queue: queue,
		root:  root,
	}
}

func createQueue(freqTable *map[rune]int) *Queue {
	queue := make(Queue, 0, len(*freqTable))
	id := 0
	for k, v := range *freqTable {
		queue = append(queue, getNewLeafNode(k, v, id))
		id++
	}
	sort.Sort(queue)
	return &queue
}

func createBinaryTree(queue *Queue) *Node {
	for queue.Len() > 1 {
		left := (*queue)[0]
		right := (*queue)[1]
		internal := getNewInternalNode(left, right)
		*queue = (*queue)[2:]
		*queue = append(*queue, internal)
		sort.Sort(*queue)
	}
	if queue.Len() == 0 {
		return nil
	}
	return (*queue)[0]
}

func (t *BinaryTree) GetPrefixCodeTable() {
	t.CodeTable = make(map[rune]string)

	t.root.Traverse("", &t.CodeTable)
}

func (t *BinaryTree) GetCodeTableAsString() string {
	var result string
	for k, v := range t.CodeTable {
		result += fmt.Sprintf("%c:%s,", k, v)
	}

	return result
}

func (t *BinaryTree) GetCompressedText(text string) string {
	var result string
	for _, char := range text {
		result += t.CodeTable[char]
	}

	return result
}

func (t *BinaryTree) GetCodeTable() *CodeTable {
	return &t.CodeTable
}

func (n *Node) Traverse(seed string, codeTable *CodeTable) {
	if n == nil {
		return
	}
	if n.Left != nil {
		n.Left.Traverse(seed+"0", codeTable)
	}
	if n.Right != nil {
		n.Right.Traverse(seed+"1", codeTable)
	}
	if n.Left == nil && n.Right == nil {
		if seed == "" {
			(*codeTable)[n.Value] = "0"
			return
		}
		(*codeTable)[n.Value] = seed
	}
}

type Queue []*Node

func (q Queue) Len() int {
	return len(q)
}
func (q Queue) Less(i, j int) bool {
	if len(q) == 0 {
		return false
	}
	if q[i].Count == q[j].Count {
		return q[i].ID < q[j].ID
	}
	return q[i].Count < q[j].Count
}
func (q Queue) Swap(i, j int) {
	if len(q) == 0 {
		return
	}
	q[i], q[j] = q[j], q[i]
}

func getNewLeafNode(value rune, count, id int) *Node {
	return &Node{
		Value: value,
		Count: count,
		ID:    id,
	}
}

func getNewInternalNode(left, right *Node) *Node {
	return &Node{
		Left:  left,
		Right: right,
		Count: left.Count + right.Count,
		ID:    left.ID + right.ID,
	}
}
