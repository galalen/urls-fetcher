package utils

import (
	"encoding/json"
	"fmt"
	"unicode"
)

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
