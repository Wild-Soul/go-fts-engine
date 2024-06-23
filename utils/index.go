package utils

type Index map[string][]int

func (idx Index) Add(docs []document) {
	for _, doc := range docs {
		tokens := analyze(doc.Text)
		for _, token := range tokens {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.Id {
				continue
			}
			idx[token] = append(ids, doc.Id)
		}
	}
}

func Interection(a, b []int) []int {
	// Find the common values between two arrays.
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	res := make([]int, 0, n)
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

func (idx Index) Search(text string) []int {
	// TODO:: Make res a Set, and have methods like Intersection and join.
	var res []int
	// search query should also follow the same trasformations that was done during indexing.
	for _, token := range analyze(text) {
		// check if token is present in Index
		if ids, ok := idx[token]; ok {
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
