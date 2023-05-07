package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	client := http.Client{}
	res, err := client.Do(req)

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

func worker(jobs <-chan string, m *sync.Map) {
	for j := range jobs {
		fetchDataAndProcess(j, m)
	}
}

func main() {
	urls := fileops.ReadUrls(urlsFile)
	m := &sync.Map{}

	numJobs := len(urls)
	jobs := make(chan string, numJobs)
	numWorkers := 10

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			defer wg.Done()
			worker(jobs, m)
		}(w)
	}

	// rand.Seed(time.Now().UnixNano())
	delay := 5
	for i, url := range urls {
		fmt.Printf("i => (%d)\n", i+1)
		if i > 0 && i%10 == 0 {
			// randomDelay := rand.Intn(delay) + 1
			fmt.Printf("Sleeping for %d sec\n", delay)
			// time.Sleep(time.Duration(randomDelay) * time.Second)
			time.Sleep(time.Duration(delay) * time.Second)
		}
		if i > 0 && i%1000 == 0 {
			delay += 5
		}
		jobs <- url
	}
	close(jobs)

	wg.Wait()

	topWord := 10
	fmt.Printf("Viewing top %d words:\n", topWord)
	words := utils.GetTopNWords(m, topWord)
	utils.PrettyPrint(words)
}
