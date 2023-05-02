package main

import (
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/alitto/pond"
	"github.com/galalen/urls-fetcher/fileops"
	"github.com/galalen/urls-fetcher/utils"
)

const (
	urlsFile  = "data/endg-urls"
	wordsFile = "data/words.txt"
)

var wordBank = fileops.GetFilteredWordBank(wordsFile)

func processDoc(r io.ReadCloser, m *sync.Map) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Println(err)
		return
	}

	article := doc.Find(".article-text").Text()
	for _, txt := range strings.Fields(strings.ToLower(article)) {
		if !(utils.IsValidWord(txt) && utils.IsValidFromBank(txt, wordBank)) {
			continue
		}

		if v, ok := m.Load(txt); ok {
			m.Store(txt, v.(int)+1)
		} else {
			m.Store(txt, 1)
		}

	}

}

func fetchDataAndProcess(url string, m *sync.Map) {
	log.Printf("Fetching: %s ...", url)
	// client := &http.Client{Timeout: 10 * time.Second}
	// res, err := client.Get(url)

	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching %s:%s", url, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	log.Printf("Processing: %s ...", url)
	processDoc(res.Body, m)
}

func getTopNWords(m *sync.Map, n int) []string {
	var pairs []struct {
		key   string
		value int
	}

	m.Range(func(k, v interface{}) bool {
		pairs = append(pairs, struct {
			key   string
			value int
		}{k.(string), v.(int)})

		return true
	})

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value > pairs[j].value
	})

	var keys []string
	for _, pair := range pairs {
		keys = append(keys, pair.key)
	}

	if n > len(pairs) {
		return keys
	}

	return keys[:n]
}

func main() {
	urls := fileops.ReadUrls(urlsFile)

	m := &sync.Map{}

	pool := pond.New(1000, 10000)
	for _, url := range urls {
		pool.Submit(func() {
			fetchDataAndProcess(url, m)
		})
	}
	pool.StopAndWait()

	words := getTopNWords(m, 10)
	utils.PrettyPrint(words)
}
