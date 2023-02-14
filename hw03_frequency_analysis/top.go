package hw03frequencyanalysis

import (
	"errors"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	errorValidString = errors.New("string contains more 1 words")
	N                = 10
)

func getFormattedWord(str string) (string, error) {
	if len(strings.Fields(str)) > 1 {
		return "", errorValidString
	}
	str = strings.ToLower(str)
	matcher := `[^\p{P}]+[\p{P}]*[^\p{P}]+`
	if utf8.RuneCountInString(str) == 1 {
		matcher = `[^\p{P}]`
	}
	validID := regexp.MustCompile(matcher)
	formatted := validID.FindString(str)
	return formatted, nil
}

func Top10(input string) []string {
	words := strings.Fields(input)
	counter := make(map[string]int)
	for _, word := range words {
		word, _ = getFormattedWord(word)
		if len(word) == 0 {
			continue
		}
		counter[word]++
	}

	keys := make([]string, 0, len(counter))
	for k := range counter {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if counter[keys[i]] == counter[keys[j]] {
			return keys[i] < keys[j]
		}
		return counter[keys[i]] > counter[keys[j]]
	})

	if len(keys) > N {
		keys = keys[:N]
	}

	return keys
}
