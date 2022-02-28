//go:build !solution
// +build !solution

package spacecollapse

import "unicode"

func CollapseSpaces(input string) string {
	runes := []rune(input)
	prevSpace := false
	l := 0
	for r := 0; r < len(runes); r++ {
		isSpace := unicode.IsSpace(runes[r])
		if !isSpace {
			runes[l] = runes[r]
			l++
		} else if !prevSpace {
			runes[l] = ' '
			l++
		}
		prevSpace = isSpace
	}
	return string(runes[:l])
}
