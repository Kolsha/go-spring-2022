//go:build !solution
// +build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

//func CollapseSpaces(input string) string {
//	runes := []rune(input)
//	prevSpace := false
//	l := 0
//	for r := 0; r < len(runes); r++ {
//		isSpace := unicode.IsSpace(runes[r])
//		if !isSpace {
//			runes[l] = runes[r]
//			l++
//		} else if !prevSpace {
//			runes[l] = ' '
//			l++
//		}
//		prevSpace = isSpace
//	}
//	return string(runes[:l])
//}

func CollapseSpaces(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	prevSpace := false
	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)
		input = input[size:]

		isSpace := unicode.IsSpace(r)
		if !isSpace {
			sb.WriteRune(r)
		} else if !prevSpace {
			sb.WriteRune(' ')
		}
		prevSpace = isSpace

	}

	return sb.String()
}
