package messages

import (
	"reflect"
	"testing"
)

func TestVerifyOIPTransfer(t *testing.T) {
	s := `
			{
				"oip-041": {
					"transferArtifact": {
						"txid": "d5fa5f01038afb6537ea517fcb107eaaee2a6834997b7b7265f580beaec5a1b4",
						"to": "FLuiVU5iDQ4a6ztcpBLwBNjBisyY2DvUTV",
						"from": "FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q",
						"timestamp": 1481738812
					},
					"signature": "Hx6PXThfI1OTguwZHLbJ64BmMWPs2n1hYTzEjljXhwaIRBD+uNGILNAB50CJBupx19331tZUMuPRkLKrf4YFapc="
				}
			}`
	oip_t := Oip041{
		Transfer: Oip041Transfer{
			Reference: "d5fa5f01038afb6537ea517fcb107eaaee2a6834997b7b7265f580beaec5a1b4",
			Timestamp: 1481738812,
			From:      "FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q",
			To:        "FLuiVU5iDQ4a6ztcpBLwBNjBisyY2DvUTV",
		},
		artSize:   len(s),
		Signature: "Hx6PXThfI1OTguwZHLbJ64BmMWPs2n1hYTzEjljXhwaIRBD+uNGILNAB50CJBupx19331tZUMuPRkLKrf4YFapc=",
	}
	s2 := `aye`

	cases := []struct {
		in    string
		out   Oip041
		block int
		err   error
	}{
		{s, Oip041{}, 1, ErrTooEarly},
		{s, oip_t, 2010000, nil},
		{s2, Oip041{}, 2010000, ErrNotJSON},
	}

	for i, c := range cases {
		o, err := VerifyOIP041(c.in, c.block)
		if err == nil {
			o, err = VerifyOIP041Transfer(o)
		}
		if err != c.err {
			t.Errorf("VerifyOIP041Transfer(#%d) | err == %v, want %v", i, err, c.err)
		}
		if err == nil && !reflect.DeepEqual(o, c.out) {
			t.Errorf("VerifyOIP041Transfer(#%d) | got == %#v, want %#v", i, o, c.out)
		}
	}
}
