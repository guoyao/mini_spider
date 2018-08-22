package spider

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
)

type Seed struct {
	url string
}

func NewSeed(url string) *Seed {
	return &Seed{url: url}
}

func (seed *Seed) getHost() (string, error) {
	u, err := url.Parse(seed.url)
	if err != nil {
		return "", err
	}

	return u.Hostname(), nil
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

func getHosts(seeds []*Seed) []string {
	hosts := make([]string, 0, len(seeds))
	for _, seed := range seeds {
		host, err := seed.getHost()
		if err == nil {
			hosts = append(hosts, host)
		}
	}
	return hosts
}
