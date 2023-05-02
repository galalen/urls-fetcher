package fileops

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/galalen/urls-fetcher/utils"
)

func ReadUrls(path string) (urls []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error closing scanner: %s", err)
	}

	return
}

func GetFilteredWordBank(path string) (words []string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if !utils.IsValidWord(word) {
			continue
		}

		words = append(words, strings.ToLower(word))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error closing scanner: %s", err)
	}

	return
}
