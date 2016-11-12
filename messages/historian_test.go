package messages

import (
	"fmt"
	"os"
	"testing"
)

func TestVerifyHistorianMessage(t *testing.T) {
	s := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="

	if os.Getenv("F_USER") == "" {
		t.Skip("skipping test; $F_TOKEN not set")
	}
	if os.Getenv("F_TOKEN") == "" {
		t.Skip("skipping test; $F_TOKEN not set")
	}

	// Don't need heavy testing of true address validity
	// The heavy lifting is done by the FlorinCoin daemon
	cases := []struct {
		in  string
		err error
	}{
		{s, nil},
		{s[:len(s)-1] + "a", ErrHistorianMessageBadSignature},
	}

	for _, c := range cases {
		got, err := VerifyHistorianMessage([]byte(c.in))
		if err != c.err {
			t.Errorf("VerifyHistorianMessage(%q) | err == %q, want %q", c.in, err, c.err)
		}
		// ToDo: check the decoded result
		fmt.Printf("%v\n", got)
		//if got != c.out {
		//	t.Errorf("CheckAddress(%q) == %q, want %q", c.in, got, c.out)
		//}
	}
}
