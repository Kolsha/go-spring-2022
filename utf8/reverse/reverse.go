//go:build !solution
// +build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

//func ReverseOld(input string) string {
//	runes := []rune(input)
//	for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
//		runes[l], runes[r] = runes[r], runes[l]
//	}
//	return string(runes)
//}

func Reverse(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	for len(input) > 0 {
		r, size := utf8.DecodeLastRuneInString(input)
		sb.WriteRune(r)
		input = input[:len(input)-size]
	}

	return sb.String()
}
