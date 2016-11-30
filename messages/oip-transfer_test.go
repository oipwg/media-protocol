package messages

import (
	"reflect"
	"testing"
)

func TestVerifyOIPTransfer(t *testing.T) {

	s := `{"oip-transfer":{"tx":"x","to":"x","fro":"F8gFhCVvcBv18fQNf5U3RZ6Zotgcjy8JnF","ts":66,"sig":"INIbrjpLIF4dnsMAUtEu5ETBvbqKVbBs+UCASs8oiLp+RFIg7xsaRebPgzEMBYYHIeWMrM0mmLbiHb3OfJCZBxw="}}`
	s2 := `{"oip-transfer":{"tx":"6f9c23edabc92e5738491a269a66aa469e03fc61156084909b646131e92ab985","to":"FLHP1SVdUSWWWnU43qBotWmHUBXYbwS4Ds","fro":"FNa3C96zuEtA5Zra54wkLpMZ6mRvTCo5uG","ts":1480321281,"sig":"IBdSrcPP6NLcKH5zvlPivy9M5p8O4tWSRVXrx5CCPACnKko761MhlkOoZE35XJcblhbEQx173Qg8AELBVDu1+V4="}}`
	oip_t2 := OIPTransfer{
		"6f9c23edabc92e5738491a269a66aa469e03fc61156084909b646131e92ab985",
		"FLHP1SVdUSWWWnU43qBotWmHUBXYbwS4Ds",
		"FNa3C96zuEtA5Zra54wkLpMZ6mRvTCo5uG",
		1480321281,
		"IBdSrcPP6NLcKH5zvlPivy9M5p8O4tWSRVXrx5CCPACnKko761MhlkOoZE35XJcblhbEQx173Qg8AELBVDu1+V4=",
	}
	s3 := `aye`
	nilOIPT := OIPTransfer{}

	cases := []struct {
		in    string
		out   OIPTransfer
		block int
		err   error
	}{
		{s, nilOIPT, 1, ErrTooEarly},
		{s, nilOIPT, 2000000, ErrInvalidReference},
		{s2, oip_t2, 2000000, nil},
		{s3, nilOIPT, 2000000, ErrWrongPrefix},
	}

	for i, c := range cases {
		got, err := VerifyOIPTransfer(c.in, c.block)
		if err != c.err {
			t.Errorf("VerifyOIPTransfer(#%d) | err == %q, want %q", i, err, c.err)
		}
		if err == nil && !reflect.DeepEqual(got, c.out) {
			t.Errorf("VerifyOIPTransfer(#%d) | got == %q, want %q", i, got, c.out)
		}
	}
}
