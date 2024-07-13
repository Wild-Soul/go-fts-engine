package ds

import (
	"sync"
)

// TrieNode represents a node in the Trie
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	mu       sync.RWMutex
}

// NewTrieNode creates a new TrieNode
func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		isEnd:    false,
	}
}

// Init a child node
func (tn *TrieNode) InitChild(char rune) {
	tn.children[char] = NewTrieNode()
}

// Getter for children nodes
func (tn *TrieNode) GetChild(char rune) (*TrieNode, bool) {
	trieNode, exists := tn.children[char]
	return trieNode, exists
}

// Set isEnd property of TrieNode
func (tn *TrieNode) SetIsEnd(val bool) {
	tn.isEnd = val
}

// Trie represents the entire Trie data structure
type Trie struct {
	root *TrieNode
}

// NewTrie creates and returns a new Trie
func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

// Getter for root
func (t *Trie) GetRoot() *TrieNode {
	return t.root
}

// Insert adds a word to the Trie
func (t *Trie) Insert(word string) {
	node := t.root
	for _, char := range word {
		node.mu.Lock()
		if _, exists := node.children[char]; !exists {
			node.children[char] = &TrieNode{
				children: make(map[rune]*TrieNode),
				isEnd:    false,
			}
		}
		child := node.children[char]
		node.mu.Unlock()
		node = child
	}
	node.mu.Lock()
	defer node.mu.Unlock()
	node.isEnd = true
}

// Search checks if a word exists in the Trie
func (t *Trie) Search(word string) bool {
	node := t.root
	for _, char := range word {
		node.mu.RLock()
		child, exists := node.children[char]
		node.mu.RUnlock()
		if !exists {
			return false
		}
		node = child
	}
	node.mu.RLock()
	defer node.mu.RUnlock()
	return node.isEnd
}

// StartsWith checks if there is any word in the Trie that starts with the given prefix
func (t *Trie) StartsWith(prefix string) bool {
	node := t.root
	for _, char := range prefix {
		node.mu.RLock()
		child, exists := node.children[char]
		node.mu.RUnlock()
		if !exists {
			return false
		}
		node = child
	}
	return true
}

// FuzzySearch finds words in the Trie that are within the given maxDistance of the query
func (t *Trie) FuzzySearch(query string, maxDistance int) []string {
	results := []string{}
	var dfs func(node *TrieNode, current []rune, depth int, prevRow []int)

	dfs = func(node *TrieNode, current []rune, depth int, prevRow []int) {
		node.mu.RLock()
		defer node.mu.RUnlock()

		columns := len(query) + 1
		currentRow := make([]int, columns)
		currentRow[0] = depth

		for i := 1; i < columns; i++ {
			insertCost := currentRow[i-1] + 1
			deleteCost := prevRow[i] + 1
			replaceCost := prevRow[i-1]
			if rune(query[i-1]) != rune(current[depth-1]) {
				replaceCost++
			}
			currentRow[i] = min(insertCost, min(deleteCost, replaceCost))
		}

		if currentRow[len(query)] <= maxDistance && node.isEnd {
			results = append(results, string(current))
		}

		if minValue(currentRow) <= maxDistance {
			for char, child := range node.children {
				dfs(child, append(current, char), depth+1, currentRow)
			}
		}
	}

	dfs(t.root, []rune{}, 0, make([]int, len(query)+1))
	return results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minValue(row []int) int {
	m := row[0]
	for _, v := range row[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
