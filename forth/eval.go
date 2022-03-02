//go:build !solution
// +build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	stack []int
	words map[string]string
	curr  string
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{stack: []int{}, words: map[string]string{}, curr: ""}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	var opErr error
	ops := strings.Split(row, " ")
	for len(ops) > 0 {
		word := strings.ToLower(ops[0])
		ops = ops[1:]
		if opErr != nil {
			break
		}

		if len(word) == 0 {
			continue
		}

		if command, ok := e.words[word]; ok {
			_, opErr = e.Process(command)
			continue
		}

		if intWord, intErr := strconv.Atoi(word); intErr == nil {
			e.stack = append(e.stack, intWord)
			continue
		}

		if word == "+" || word == "-" || word == "/" || word == "*" {
			opErr = e.arithmetic(word)
			continue
		}

		stackLen := len(e.stack)
		switch word {
		case "dup":
			opErr = e.dup(stackLen - 1)
		case "over":
			opErr = e.dup(stackLen - 2)
		case "drop":
			if stackLen <= 0 {
				opErr = fmt.Errorf("invalid drop")
			} else {
				e.stack = e.stack[:stackLen-1]
			}
		case "swap":
			if stackLen <= 1 {
				opErr = fmt.Errorf("invalid swap")
			} else {
				e.stack[stackLen-1], e.stack[stackLen-2] = e.stack[stackLen-2], e.stack[stackLen-1]
			}
		case ":":
			ops, opErr = e.userDefined(ops)
		default:
			opErr = fmt.Errorf("invalid word")
		}

	}
	return e.stack, opErr
}

func (e *Evaluator) userDefined(ops []string) ([]string, error) {
	name := ""
	for len(name) == 0 && len(ops) > 0 {
		name = strings.ToLower(ops[0])
		ops = ops[1:]
	}

	if len(name) == 0 {
		return ops, fmt.Errorf("empty comand name")
	}

	if _, intErr := strconv.Atoi(name); intErr == nil {
		return ops, fmt.Errorf("redefine numbers")
	}

	command := ""
	for len(ops) > 0 && ops[0] != ";" {
		o := strings.ToLower(ops[0])
		if cmd, ok := e.words[o]; ok {
			o = cmd
		}
		command += o + " "
		ops = ops[1:]
	}
	ops = ops[1:]

	if len(command) == 0 {
		return ops, fmt.Errorf("empty comand")
	}
	e.words[name] = command
	return ops, nil
}

func (e *Evaluator) dup(i int) error {
	if i < 0 {
		return fmt.Errorf("invalid dup")
	}
	e.stack = append(e.stack, e.stack[i])
	return nil
}

func (e *Evaluator) arithmetic(o string) error {
	size := len(e.stack)
	if size < 2 {
		return fmt.Errorf("to few args at stack")
	}
	lhs, rhs := e.stack[size-2], e.stack[size-1]
	e.stack = e.stack[:size-1]
	switch o {
	case "+":
		e.stack[size-2] = lhs + rhs
	case "-":
		e.stack[size-2] = lhs - rhs
	case "/":
		if rhs == 0 {
			return fmt.Errorf("division by zero")
		}
		e.stack[size-2] = lhs / rhs
	case "*":
		e.stack[size-2] = lhs * rhs
	}

	return nil
}
