//go:build !solution
// +build !solution

package testequal

import (
	"fmt"
	"reflect"
)

func equal(expected, actual interface{}) bool {
	if reflect.ValueOf(expected).Kind() == reflect.Struct {
		return false
	}
	if reflect.ValueOf(actual).Kind() == reflect.Struct {
		return false
	}
	return reflect.DeepEqual(expected, actual)
}

func errorf(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`
	msg := ``
	l := len(msgAndArgs)
	switch l {
	case 0:
		break
	case 1:
		msg = msgAndArgs[0].(string)
	default:
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	t.Errorf(format, expected, actual, msg)

}

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if equal(expected, actual) {
		return true
	}
	errorf(t, expected, actual, msgAndArgs...)

	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	if !equal(expected, actual) {
		return true
	}
	errorf(t, expected, actual, msgAndArgs...)

	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if equal(expected, actual) {
		return
	}
	errorf(t, expected, actual, msgAndArgs...)
	t.FailNow()
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !equal(expected, actual) {
		return
	}
	errorf(t, expected, actual, msgAndArgs...)
	t.FailNow()
}
