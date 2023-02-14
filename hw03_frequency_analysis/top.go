package hw03frequencyanalysis

import (
	"errors"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	errorValidString = errors.New("string contains more 1 words")
	N                = 10
)

func split(str string) []string {
	f := func(c rune) bool {
		return unicode.IsMark(c) || unicode.IsSpace(c)
	}
	return strings.FieldsFunc(str, f)
}

func getFormattedWord(str string) (string, error) {
	if len(split(str)) > 1 {
		return "", errorValidString
	}
	str = strings.ToLower(str)
	matcher := `[^\p{P}]+[\p{P}\p{L}]*[^\p{P}]+`
	if utf8.RuneCountInString(str) <= 3 {
		matcher = `[^\p{P}]+`
	}
	validID := regexp.MustCompile(matcher)
	formatted := validID.FindString(str)
	return formatted, nil
}

func Top10(input string) []string {
	words := split(input)
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
