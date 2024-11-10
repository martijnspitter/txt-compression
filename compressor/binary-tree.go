package compressor

import (
	"container/heap"
	"fmt"
	"strings"
)

type Node struct {
	Left  *Node
	Right *Node
	Value byte
	Count int
	ID    int
}

type BinaryTree struct {
	root      *Node
	queue     *Queue
	CodeTable CodeTable
}

type CodeTable map[byte]string

func NewBinaryTree(freqTable *map[byte]int) *BinaryTree {
	queue := createQueue(freqTable)
	root := createBinaryTree(queue)
	return &BinaryTree{
		queue: queue,
		root:  root,
	}
}

func createQueue(freqTable *map[byte]int) *Queue {
	queue := &Queue{}
	heap.Init(queue)
	id := 0
	for k, v := range *freqTable {
		heap.Push(queue, getNewLeafNode(k, v, id))
		id++
	}
	return queue
}

func createBinaryTree(queue *Queue) *Node {
	for queue.Len() > 1 {
		left := heap.Pop(queue).(*Node)
		right := heap.Pop(queue).(*Node)
		internal := getNewInternalNode(left, right)
		heap.Push(queue, internal)
	}
	if queue.Len() == 0 {
		return nil
	}
	return heap.Pop(queue).(*Node)
}

func (t *BinaryTree) GetPrefixCodeTable() {
	t.CodeTable = make(map[byte]string)

	t.root.Traverse("", &t.CodeTable)
}

func (t *BinaryTree) GetCodeTableAsString() string {
	var builder strings.Builder
	for k, v := range t.CodeTable {
		builder.WriteString(fmt.Sprintf("%c:%s,", k, v))
	}

	return builder.String()
}

func (t *BinaryTree) GetCompressedText(text []byte) string {
	// Pre-calculate the final size to avoid reallocations
	totalSize := 0
	for _, b := range text {
		totalSize += len(t.CodeTable[b])
	}

	// Use a single builder with pre-allocated capacity
	builder := strings.Builder{}
	builder.Grow(totalSize)

	for _, b := range text {
		builder.WriteString(t.CodeTable[b])
	}

	return builder.String()
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

func (q *Queue) Push(n interface{}) {
	*q = append(*q, n.(*Node))
}

func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	node := old[n-1]
	*q = old[0 : n-1]
	return node
}

func getNewLeafNode(value byte, count, id int) *Node {
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
