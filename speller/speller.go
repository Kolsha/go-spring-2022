//go:build !solution

package speller

import "math"

func pow(p int) int64 {
	return int64(math.Pow(1000, float64(p)))
}

func Spell(n int64) string {
	if n < 0 {
		return "minus " + Spell(-n)
	}
	to19 := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve",
		"thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}

	tens := []string{"twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
	if n == 0 {
		return "zero"
	}
	if n < 20 {
		return to19[n-1]
	}

	spellNonZero := func(num int64, delimiter string) string {
		if num == 0 {
			return ""
		}
		return delimiter + Spell(num)
	}
	if n < 100 {
		return tens[n/10-2] + spellNonZero(n%10, "-")
	}
	if n < 1000 {
		return to19[n/100-1] + " hundred" + spellNonZero(n%100, " ")
	}

	for idx, w := range []string{"thousand", "million", "billion"} {
		p := idx + 1
		if n < pow(p+1) {
			return Spell(n/pow(p)) + " " + w + spellNonZero(n%pow(p), " ")
		}
	}

	panic("number is strange")
}
