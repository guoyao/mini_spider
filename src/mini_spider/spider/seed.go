package spider

import (
	"encoding/json"
	"io/ioutil"
)

type Seed struct {
	url string
}

func NewSeed(url string) *Seed {
	return &Seed{url: url}
}

func LoadSeedsFromFile(filePath string) ([]*Seed, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var urls []string
	if err := json.Unmarshal(content, &urls); err != nil {
		return nil, err
	}

	seeds := make([]*Seed, len(urls))
	for index, url := range urls {
		seeds[index] = NewSeed(url)
	}

	return seeds, nil
}
