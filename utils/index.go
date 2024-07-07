package utils

import (
	"sync"

	"github.com/Wild-Soul/go-fts-engine/ds"
)

type Index struct {
	ds.SafeMap[string, []string]
}

func NewIndex() *Index {
	return &Index{
		SafeMap: *ds.NewSafeMap[string, []string](),
	}
}

func (idx *Index) processAdd(ch <-chan document) {
	doc := <-ch
	tokens := analyze(doc.Text)
	for _, token := range tokens {
		if ids, exists := idx.Get(token); exists {
			// Check if doc.Id already exists in the slice
			for _, id := range ids {
				if id == doc.Id {
					return // docID already exists, no need to add
				}
			}
			// Append docID if it doesn't exist
			idx.Set(token, append(ids, doc.Id))
		} else {
			// Create new slice with docID if token doesn't exist
			idx.Set(token, []string{doc.Id})
		}
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

func (idx *Index) Search(text string) []string {
	// TODO:: Make res a Set, and have methods like Intersection and join.
	var res []string
	// search query should also follow the same trasformations that was done during indexing.
	for _, token := range analyze(text) {
		// check if token is present in Index
		if ids, exists := idx.Get(token); exists {
			if res == nil {
				res = ids
			} else {
				res = Interection(res, ids) // find the document ids that match for each token
			}
		} else {
			// TODO:: Make it a score based response, return all/top n docs that matches at-least one work from query string.
			// For each doc that contains a distinct token, increase the score.
			return nil // token doesn't exists
		}
	}
	return res
}
