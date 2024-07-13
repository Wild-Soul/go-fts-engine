package utils

import (
	"compress/gzip"
	"encoding/xml"
	"os"
)

type document struct {
	Title string `xml:"title"`
	Url   string `xml:"url"`
	Text  string `xml:"abstract"`
	Id    string
}

func LoadDocuments(path string) ([]document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(gz)
	dump := struct {
		Documents []document `xml:"doc"`
	}{}

	if err = dec.Decode(&dump); err != nil {
		return nil, err
	}

	docs := dump.Documents
	for i, doc := range docs {
		docs[i].Id = doc.Url // TODO:: this should be uuid
	}
	return docs, nil
}
