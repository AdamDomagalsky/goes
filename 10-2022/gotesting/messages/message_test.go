package messages

//package messages_test // black box testing - all small caps not available

import "testing"

func TestGreet(t *testing.T) {
	got := Greet("Gopher")
	expect := "Hello, Gopher!\n"

	if got != expect {
		t.Errorf("Did not expected result. Wanted %q, got %q\n", got, expect)
	}
}

func TestDepart(t *testing.T) {
	got := depart("Gopher")
	expect := "Goodbye, Gopher!\n"

	if got != expect {
		t.Errorf("Did not expected result. Wanted %q, got %q\n", got, expect)
	}
}

func TestFailureTest(t *testing.T) {
	t.Error("not-immediate error")
	t.Fatal("immediate error")
	t.Errorf("wont see it")
}

func TestGreetTableDriven(t *testing.T) {
	scenarios := []struct {
		input  string
		expect string
	}{
		{input: "Gopher", expect: "Hello, Gopher!\n"},
		{input: "Ewa", expect: "Hello, Ewa!\n"},
	}
	for _, s := range scenarios {
		got := Greet(s.input)
		if got != s.expect {
			t.Errorf("\nInput:%q \nGot: %q\nExpected:%q\n", s.input, got, s.expect)
		}
	}
}
