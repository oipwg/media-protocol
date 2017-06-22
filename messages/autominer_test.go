package messages

import (
	"reflect"
	"testing"
)

func TestVerifyAutominer(t *testing.T) {

	// valid signature and valid json (should work)
	s := []byte(`{ "alexandria-autominer": { "FLOAddress": "FH96569yCZ4wqK5npC4TT2RVtkB14DHs4V", "BTCAddress": "1helloworldtest", "Version": 1, "Signature": "H3VQDzcglDT9sSvGBmceiwpaXtY6dhgHjZ65QXCL/tUIdp2374+cdgpeHtJRc3ElVk64+7iJjnfruDDs/b9mk9k=" } }`)
	am := AlexandriaAutominer{}

	cases := []struct {
		in    []byte
		out   AlexandriaAutominer
		block int
		err   error
	}{
		{s, am, 1750000, nil}, // valid
		/*
			{s, am, 1974560, ErrBadSignature},                   // wrong address
			{s, am, 1974560, ErrWrongPrefix},                   // bad version
			{s, am, 1974560, ErrHistorianMessagePoolUntrusted}, // bad pool
			{s, am, 1974560, ErrHistorianMessageInvalid},       // wrong length
			{s, am, 1974560, ErrBadSignature},                  // no signature
			{s, am, 1974559, nil},                                // no signature, but unenforced
			{s, am, 1974560, ErrBadSignature},                  // trailing :
			{s, am, 1974560, ErrBadSignature},                  // invalid signature
			{s, am, 1974560, nil},                                // valid
		*/
	}

	for i, c := range cases {
		got, err := VerifyAutominer([]byte(c.in), c.block)
		if err != c.err {
			t.Errorf("VerifyAutominer(#%d) | err == %q, want %q", i, err, c.err)
		}
		if err == nil && !reflect.DeepEqual(got, c.out) {
			t.Errorf("VerifyMediaMultipartSingle(#%d) | got == %v, want %v", i, got, c.out)
		}
	}
}

/*
func TestVerifyHistorianMessage(t *testing.T) {
	// signed FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD
	// valid
	s := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	hm := HistorianMessage{
		1, "pool.alexandria.io", 0.0001360085, 3.163064456533333e+08, math.Inf(-1), 5e-06, 0.00217, "IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA=",
	}
	// bad version
	s1 := "alexandria-historian-v002:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// bad pool
	s2 := "alexandria-historian-v001:notpool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// wrong length
	s3 := "alexandria-historian-v001:"
	// no signature
	s4 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217"
	hm4 := HistorianMessage{
		1, "pool.alexandria.io", 0.0001360085, 3.163064456533333e+08, math.Inf(-1), 5e-06, 0.00217, "",
	}
	// trailing :
	s5 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:"
	// invalid signature
	s6 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwAa"
	// signed FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU
	// valid
	s7 := "alexandria-historian-v001:pool.alexandria.io:0.000104048500:223208386.28518352:2214713879:0.00000429:0.00308:ICyn+Wh4OxKF89+O9u0wkQULeyvJ6CDurGiZACCkNtk8Rl+QpejBmPWKYiuyt6PM5+MrUs/gDcACWjKFTSoYrxA="
	hm7 := HistorianMessage{
		1, "pool.alexandria.io", 0.0001040485, 2.2320838628518352e+08, 2.214713879e+09, 4.29e-06, 0.00308, "ICyn+Wh4OxKF89+O9u0wkQULeyvJ6CDurGiZACCkNtk8Rl+QpejBmPWKYiuyt6PM5+MrUs/gDcACWjKFTSoYrxA=",
	}
	nilHM := HistorianMessage{}

	cases := []struct {
		in    string
		out   HistorianMessage
		block int
		err   error
	}{
		{s, hm, 1750000, nil},                                  // valid
		{s, nilHM, 1974560, ErrBadSignature},                   // wrong address
		{s1, nilHM, 1974560, ErrWrongPrefix},                   // bad version
		{s2, nilHM, 1974560, ErrHistorianMessagePoolUntrusted}, // bad pool
		{s3, nilHM, 1974560, ErrHistorianMessageInvalid},       // wrong length
		{s4, nilHM, 1974560, ErrBadSignature},                  // no signature
		{s4, hm4, 1974559, nil},                                // no signature, but unenforced
		{s5, nilHM, 1974560, ErrBadSignature},                  // trailing :
		{s6, nilHM, 1974560, ErrBadSignature},                  // invalid signature
		{s7, hm7, 1974560, nil},                                // valid
	}

	for i, c := range cases {
		got, err := VerifyHistorianMessage([]byte(c.in), c.block)
		if err != c.err {
			t.Errorf("VerifyHistorianMessage(#%d) | err == %q, want %q", i, err, c.err)
		}
		if err == nil && !reflect.DeepEqual(got, c.out) {
			t.Errorf("VerifyMediaMultipartSingle(#%d) | got == %v, want %v", i, got, c.out)
		}
	}
}
*/
