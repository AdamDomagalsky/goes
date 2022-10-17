package main_test

import "testing"

func TestAddition(t *testing.T) {
	got := 2 + 2
	expect := 4
	if got != expect {
		t.Errorf("Got:%v Expected%v", got, expect)
	}
}
func TestSubs(t *testing.T) {
	got := 10 - 5
	expected := 4
	if got != expected {
		t.Errorf("Got %v Expected %v", got, expected)
	}
}
