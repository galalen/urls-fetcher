package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
	"unicode"
)

func GetRandomUserAgent() string {
	rand.Seed(time.Now().Unix())

	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	}

	return agents[rand.Intn(len(agents))]
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsValidWord(word string) bool {
	return len(word) >= 3 && IsLetter(word)
}

func IsValidFromBank(word string, bank []string) bool {
	for _, txt := range bank {
		if txt == word {
			return true
		}
	}
	return false
}

func PrettyPrint(words []string) {
	b, err := json.MarshalIndent(words, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func GetTopNWords(m *sync.Map, n int) []string {
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
