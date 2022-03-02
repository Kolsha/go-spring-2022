//go:build !solution
// +build !solution

package varfmt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Sprintf(format string, args ...interface{}) string {
	argsStrings := make([]string, 0, len(args))
	argsSize := 0
	for _, arg := range args {
		s := fmt.Sprint(arg)
		argsStrings = append(argsStrings, s)
		argsSize += len(s)
	}

	var sb strings.Builder
	sb.Grow(len(format) + argsSize)

	pos := 0
	number := 0
	numberLen := -1
	for len(format) > 0 {
		r, size := utf8.DecodeRuneInString(format)
		format = format[size:]
		switch true {
		case r == '{':
			numberLen = 0
		case r == '}':
			write := number
			if numberLen <= 0 {
				write = pos
			}
			sb.WriteString(argsStrings[write])
			number = 0
			numberLen = -1
			pos++
		case numberLen >= 0:
			number = number*10 + (int(r) - '0')
			numberLen++
		default:
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
