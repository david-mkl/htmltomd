package util

import "testing"

func TestStringPointer(t *testing.T) {
	str := "test"

	if str != *String(str) {
		t.Errorf("Expected %s. Got %s", str, *String(str))
	}
}

func TestBoolPointer(t *testing.T) {
	b := true

	if b != *Bool(b) {
		t.Errorf("Expected %t. Got %t", b, *Bool(b))
	}
}
