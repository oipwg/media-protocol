package messages

import (
	"fmt"
	"os"
	"testing"
)

func TestVerifyHistorianMessage(t *testing.T) {
	// signed FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD
	// valid
	s := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// bad version
	s1 := "alexandria-historian-v002:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// bad pool
	s2 := "alexandria-historian-v001:notpool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// wrong length
	s3 := "alexandria-historian-v001:"
	// no signature
	s4 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217"
	// trailing :
	s5 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:"
	// invalid signature
	s6 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwAa"
	// signed FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU
	// valid
	s7 := "alexandria-historian-v001:pool.alexandria.io:0.000104048500:223208386.28518352:2214713879:0.00000429:0.00308:ICyn+Wh4OxKF89+O9u0wkQULeyvJ6CDurGiZACCkNtk8Rl+QpejBmPWKYiuyt6PM5+MrUs/gDcACWjKFTSoYrxA="

	if os.Getenv("F_USER") == "" {
		t.Skip("skipping test; $F_TOKEN not set")
	}
	if os.Getenv("F_TOKEN") == "" {
		t.Skip("skipping test; $F_TOKEN not set")
	}

	// Don't need heavy testing of true address validity
	// The heavy lifting is done by the FlorinCoin daemon
	cases := []struct {
		in    string
		block int
		err   error
	}{
		{s, 1750000, nil},                               // valid
		{s, 1974560, ErrHistorianMessageBadSignature},   // wrong address
		{s1, 1974560, ErrHistorianMessageInvalid},       // bad version
		{s2, 1974560, ErrHistorianMessagePoolUntrusted}, // bad pool
		{s3, 1974560, ErrHistorianMessageInvalid},       // wrong length
		{s4, 1974560, ErrHistorianMessageBadSignature},  // no signature
		{s5, 1974560, ErrHistorianMessageBadSignature},  // trailing :
		{s6, 1974560, ErrHistorianMessageBadSignature},  // invalid signature
		{s7, 1974560, nil},                              // valid
		{s4, 1974559, nil},                              // no signature, but unenforced
	}

	for i, c := range cases {
		got, err := VerifyHistorianMessage([]byte(c.in), c.block)
		if err != c.err {
			t.Errorf("VerifyHistorianMessage(#%d) | err == %q, want %q", i, err, c.err)
		}
		// ToDo: check the decoded result
		fmt.Printf("%v\n", got)
		//if got != c.out {
		//	t.Errorf("CheckAddress(%q) == %q, want %q", c.in, got, c.out)
		//}
	}
}
