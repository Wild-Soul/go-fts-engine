package utils

import (
	"sync"

	"github.com/Wild-Soul/go-fts-engine/ds"
)

const EDIT_DISTANCE = 1

type Index struct {
	invertedIndex ds.SafeMap[string, []string]
	trie          *ds.Trie
	mu            sync.RWMutex
}

func NewIndex() *Index {
	return &Index{
		invertedIndex: *ds.NewSafeMap[string, []string](),
		trie:          ds.NewTrie(),
	}
}

func (idx *Index) Insert(word string, docID string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Update inverted index
	if ids, exists := idx.invertedIndex.Get(word); exists {
		idx.invertedIndex.Set(word, append(ids, docID))
	} else {
		idx.invertedIndex.Set(word, []string{docID})
	}

	// Update Trie
	node := idx.trie.GetRoot()
	for _, char := range word {
		if _, exists := node.GetChild(char); !exists {
			node.InitChild(char)
		}
		node, _ = node.GetChild(char)
	}
	node.SetIsEnd(true)
}

func (idx *Index) processAdd(ch <-chan document) {
	doc := <-ch
	tokens := analyze(doc.Text)
	for _, token := range tokens {
		// Update Trie and inverted index
		idx.Insert(token, doc.Id)
	}
}

func (idx *Index) Add(docs []document) {
	var wg sync.WaitGroup

	ch := make(chan document, 100) // 100 docs processed concurrently
	for _, doc := range docs {
		wg.Add(1)
		ch <- doc
		go func() {
			idx.processAdd(ch)
			wg.Done()
		}()
	}

	wg.Wait() // wait for all docs to be processed.
}

func Interection(a, b []string) []string {
	// Find the common values between two arrays.
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	res := make([]string, 0, n)
	for i, j := 0, 0; i < len(a) && j < len(b); {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			res = append(res, a[i])
			i++
			j++
		}
	}
	return res
}

func findInInvertedIdx(idx *Index, token string) []string {
	if ids, exists := idx.invertedIndex.Get(token); exists {
		return ids
	}
	return []string{}
}

func (idx *Index) ExactSearch(text string) []string {
	var res []string
	// search query should also follow the same trasformations that was done during indexing.
	for _, token := range analyze(text) {
		// check if token is present in Index
		ids := findInInvertedIdx(idx, token)
		res = Interection(res, ids) // only keep those documents that match for each token in query.
	}
	return res
}

// Searches through Trie datastructure using for similar tokens.
// Then gets the docIds from inverted index for these new tokens.
func (idx *Index) FuzzySearch(text string) []string {
	var res []string
	// search query should also follow the same trasformations that was done during indexing.
	for _, token := range analyze(text) {
		// check if token is present in Index
		similarTokens := idx.trie.FuzzySearch(token, EDIT_DISTANCE)
		for _, sT := range similarTokens {
			ids := findInInvertedIdx(idx, sT)
			res = Interection(res, ids) // only keep those documents that match for each token in query.
		}
	}
	return res
}
