package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func toDigit(r rune) int {
	return int(r - '0')
}

func Unpack(str string) (string, error) {
	runes := []rune(str)
	length := len(runes)
	builder := strings.Builder{}
	i := 0

	var curChar rune
	var isEscaped bool

	for i < length {
		char := runes[i]

		if isEscaped { // case after escape char
			if char == '\\' {
				curChar = '\\'
			} else if unicode.IsLetter(char) { // letter after escape forbidden
				return "", ErrInvalidString
			} else if unicode.IsDigit(char) {
				curChar = char
			} else {
				return "", ErrInvalidString
			}

			isEscaped = false
		} else { // normal case
			if char == '\\' {
				if curChar != 0 { // write prev char
					builder.WriteRune(curChar)
				}

				isEscaped = true
			} else if unicode.IsLetter(char) {
				if curChar != 0 { // write prev char
					builder.WriteRune(curChar)
				}

				curChar = char
			} else if unicode.IsDigit(char) {
				if curChar == 0 { // only letter can be repeated
					return "", ErrInvalidString
				}

				for j := 0; j < toDigit(char); j++ {
					builder.WriteRune(curChar)
				}
				curChar = 0
			} else {
				return "", ErrInvalidString
			}
		}

		i += 1
	}

	if curChar != 0 { // write last char if exists
		builder.WriteRune(curChar)
	}

	return builder.String(), nil
}
