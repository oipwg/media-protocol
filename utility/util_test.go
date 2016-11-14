package utility

import (
	"os"
	"testing"
)

func TestCheckAddress(t *testing.T) {
	// Don't need heavy testing of true address validity
	// The heavy lifting is done by btc suite
	cases := []struct {
		in  string
		out bool
	}{
		{`FRLJFGyzEiudhjvePdyk8Gn4bkpBcoXzGv`, true},
		{`nope`, false},
	}

	for _, c := range cases {
		got := CheckAddress(c.in)
		if got != c.out {
			t.Errorf("CheckAddress(%q) == %q, want %q", c.in, got, c.out)
		}
	}
}

func TestCheckSignature(t *testing.T) {
	if os.Getenv("F_USER") == "" {
		t.Skip("skipping test; $F_TOKEN not set")
	}
	if os.Getenv("F_TOKEN") == "" {
		t.Skip("skipping test; $F_TOKEN not set")
	}

	// Don't need heavy testing of true signature validity
	// The heavy lifting is done by the FlorinCoin daemon
	cases := []struct {
		mes, sig, addr string
		out            bool
	}{
		{`a`, `IODgNXJXYTb3XcKJswL1zPzzou50wn6oDzImBvzLRuVzOeGyruuJHvEO9C1p8+gErM8xNb3ZXGjhjoznitACG2k=`, `F8gFhCVvcBv18fQNf5U3RZ6Zotgcjy8JnF`, true},
		{`b`, `IODgNXJXYTb3XcKJswL1zPzzou50wn6oDzImBvzLRuVzOeGyruuJHvEO9C1p8+gErM8xNb3ZXGjhjoznitACG2k=`, `F8gFhCVvcBv18fQNf5U3RZ6Zotgcjy8JnF`, false},
	}

	for _, c := range cases {
		got := CheckSignature(c.addr, c.sig, c.mes)
		if got != c.out {
			t.Errorf("CheckSignature(%q, <...>) == %q, want %q", c.mes, c.out, c.out)
		}
	}
}

func TestIsJSON(t *testing.T) {
	// Don't need heavy testing of true JSON syntax
	// The actual JSON parsing is in a standard library
	cases := []struct {
		in  string
		out bool
	}{
		{`{"a": "b"}`, true},
		{`nope`, false},
	}

	for _, c := range cases {
		got := IsJSON(c.in)
		if got != c.out {
			t.Errorf("IsJSON(%q) == %q, want %q", c.in, got, c.out)
		}
	}
}
