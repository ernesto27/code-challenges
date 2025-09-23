package main

import (
	"container/heap"
	"fmt"
)

type Node struct {
	char        rune
	freq        int
	left, right *Node
}

func buildFreqTable(s string) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}
	return freq
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].freq < pq[j].freq }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(*Node)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func buildHuffmanTree(freq map[rune]int) *Node {
	pq := &PriorityQueue{}
	heap.Init(pq)

	for ch, f := range freq {
		heap.Push(pq, &Node{char: ch, freq: f})
	}

	for pq.Len() > 1 {
		left := heap.Pop(pq).(*Node)
		right := heap.Pop(pq).(*Node)
		merged := &Node{
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		}
		heap.Push(pq, merged)
	}

	return heap.Pop(pq).(*Node)
}

func buildCodes(node *Node, prefix string, codes map[rune]string) {
	if node == nil {
		return
	}

	if node.left == nil && node.right == nil {
		codes[node.char] = prefix
	}

	buildCodes(node.left, prefix+"0", codes)
	buildCodes(node.right, prefix+"1", codes)
}

func encode(s string, codes map[rune]string) string {
	result := ""
	for _, ch := range s {
		val := string(ch)
		fmt.Println(val)
		result += codes[ch]
	}
	return result
}

func decode(encoded string, root *Node) string {
	result := ""
	current := root

	for _, bit := range encoded {
		if bit == '0' {
			current = current.left
		} else {
			current = current.right
		}

		if current.left == nil && current.right == nil {
			result += string(current.char)
			current = root
		}
	}
	return result
}

func main() {
	input := "MISSISSIPPI"

	freq := buildFreqTable(input)

	root := buildHuffmanTree(freq)

	codes := make(map[rune]string)
	buildCodes(root, "", codes)
	fmt.Println("Huffman Codes:", codes)

	// Step 4: Encode
	encoded := encode(input, codes)
	fmt.Println("Encoded:", encoded)

	// Step 5: Decode
	decoded := decode(encoded, root)
	fmt.Println("Decoded:", decoded)
}
