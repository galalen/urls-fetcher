package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
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
		return
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

func worker(id int, jobs <-chan string, results chan<- struct{}, m *sync.Map) {
	for j := range jobs {
		fetchDataAndProcess(j, m)
		results <- struct{}{}
	}
}

func main() {

	urls := fileops.ReadUrls(urlsFile)
	m := &sync.Map{}

	numJobs := len(urls)
	jobs := make(chan string, numJobs)
	results := make(chan struct{}, numJobs)

	numWorkers := 50
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, m)
	}

	for i, url := range urls {
		if i%100 == 0 {
			fmt.Println("Sleeping for 5 sec")
			time.Sleep(5 * time.Second)
		}
		jobs <- url
	}

	for a := 1; a <= numJobs; a++ {
		<-results
	}

	topWord := 10
	fmt.Printf("Viewing top %d words:\n", topWord)
	words := getTopNWords(m, topWord)
	utils.PrettyPrint(words)
}
