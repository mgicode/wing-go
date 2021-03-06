package tcp

import (
	"testing"
	"errors"
)

// go test -v -test.run Test_isClosedConnError
func Test_isClosedConnError(t *testing.T) {
	err := errors.New("hello")
	if false != isClosedConnError(err) {
		t.Errorf("isClosedConnError fail")
		return
	}
	if false != isClosedConnError(nil) {
		t.Errorf("isClosedConnError fail")
		return
	}
	err = errors.New("hello use of closed network connection")
	if true != isClosedConnError(err) {
		t.Errorf("isClosedConnError fail")
		return
	}
}

// go test -v -test.run TestRandString
func TestRandString(t *testing.T) {
	if 10 != len(RandString(10)) {
		t.Errorf("RandString error")
	}
}