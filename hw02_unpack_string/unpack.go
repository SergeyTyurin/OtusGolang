package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func isEscaped(str string, pos int) bool {
	result := false
	for pos != 0 {
		r, size := utf8.DecodeLastRuneInString(str[:pos])
		if r != '\\' {
			break
		}
		result = !result
		pos -= size
	}
	return result
}

func isValid(str string, pos int) bool {
	substr := str[pos:]
	ch, size := utf8.DecodeRuneInString(substr)
	_, err := strconv.Atoi(string(ch))

	checkPrevNumber := func(str string, pos int) bool {
		ch, size := utf8.DecodeLastRuneInString(str[:pos])
		_, err := strconv.Atoi(string(ch))
		if err == nil && !isEscaped(str, pos-size) {
			return false
		}
		return true
	}

	// check number
	if err == nil {
		if pos == 0 {
			return false
		}
		return checkPrevNumber(str, pos)
	}

	// check escaped character
	if isEscaped(str, pos) && ch != '\\' {
		return false
	}

	// check last character
	if ch == '\\' && (pos+size) == len(str) {
		return false
	}

	return true
}

func Unpack(str string) (string, error) {
	var strBuilder strings.Builder
	var buf string
	for pos, ch := range str {
		if !isValid(str, pos) {
			return "", ErrInvalidString
		}

		if ch == '\\' && !isEscaped(str, pos) {
			continue
		}

		number, err := strconv.Atoi(string(ch))
		if err == nil && !isEscaped(str, pos) {
			strBuilder.WriteString(strings.Repeat(buf, number))
			buf = ""
		} else {
			strBuilder.WriteString(buf)
			buf = string(ch)
		}
	}
	strBuilder.WriteString(buf)
	return strBuilder.String(), nil
}
