package messages

import (
	"reflect"
	"testing"
)

func TestVerifyMediaMultipartSingle(t *testing.T) {

	s := "alexandria-media-multipart(0,2,FAFRmxDW9an5XLBMixU1ZCF9aB6LShZjyP,0000000000000000000000000000000000000000000000000000000000000000,IKyu0J2jOvMpvPVIXvTmIOzoK5YJT2+SFvkWMPI4Xi9gR6m3iUl0Hi+vYdxrhQDpsDBYZteNcc7stnx5bl6u5R4=,):{\"alexandria-media\":{\"torrent\":\"QmcGgZCtR5RL9QVPfiCx5ACzPGntG1Qk7rxcxLe3gNj4vj\",\"publisher\":\"FAFRmxDW9an5XLBMixU1ZCF9aB6LShZjyP\",\"timestamp\":1477841854,\"type\":\"thing\",\"payment\":{},\"info\":{\"title\":\"Queen Thyra\",\"description\":\"The Danish Queen Thyra -o- CC-by - http://cre"
	mms := MediaMultipartSingle{
		Part:      0,
		Max:       2,
		Reference: "cd6be959f0a90e515ff02f7eacf12dedfaaf2b49fcbe2555b658a53b6e012129",
		Address:   "FAFRmxDW9an5XLBMixU1ZCF9aB6LShZjyP",
		Signature: "IKyu0J2jOvMpvPVIXvTmIOzoK5YJT2+SFvkWMPI4Xi9gR6m3iUl0Hi+vYdxrhQDpsDBYZteNcc7stnx5bl6u5R4=",
		Data:      "{\"alexandria-media\":{\"torrent\":\"QmcGgZCtR5RL9QVPfiCx5ACzPGntG1Qk7rxcxLe3gNj4vj\",\"publisher\":\"FAFRmxDW9an5XLBMixU1ZCF9aB6LShZjyP\",\"timestamp\":1477841854,\"type\":\"thing\",\"payment\":{},\"info\":{\"title\":\"Queen Thyra\",\"description\":\"The Danish Queen Thyra -o- CC-by - http://cre",
		Txid:      "cd6be959f0a90e515ff02f7eacf12dedfaaf2b49fcbe2555b658a53b6e012129",
		Block:     1958430,
	}
	s2 := "alexandria-media-multipart(2,2,FNoqStBA425P1ifZKu8yqUyFPXXFg9D1GK,407334fce77c882a55e3012a7b9497c43dc95215772edc744cda7681001ea800,H9tQN+1qsR4gWr3GVtUSnfny0+MzLzi5mmSo/W8yZ+NKjPDl3NhRiaENkexAFAcsce78bnaLTVuo02gyYecy5Ls=):e\":\"How_to_combine_bitcoin_payment_processing_and_charity.pdf\",\"artist\":\"SamBiohazard\"}} }, \"signature\":\"IBEByQTx68CcFBpYnV0anTAPT05DuipzmIRqjud15zZ586DSxqXA5eFn6MnO95AUi5hKVAw71fG8a+Lvi7A7CWc=\" }"
	mms2 := MediaMultipartSingle{
		Part:      2,
		Max:       2,
		Reference: "407334fce77c882a55e3012a7b9497c43dc95215772edc744cda7681001ea800",
		Address:   "FNoqStBA425P1ifZKu8yqUyFPXXFg9D1GK",
		Signature: "H9tQN+1qsR4gWr3GVtUSnfny0+MzLzi5mmSo/W8yZ+NKjPDl3NhRiaENkexAFAcsce78bnaLTVuo02gyYecy5Ls=",
		Data:      "e\":\"How_to_combine_bitcoin_payment_processing_and_charity.pdf\",\"artist\":\"SamBiohazard\"}} }, \"signature\":\"IBEByQTx68CcFBpYnV0anTAPT05DuipzmIRqjud15zZ586DSxqXA5eFn6MnO95AUi5hKVAw71fG8a+Lvi7A7CWc=\" }",
		Txid:      "682f7c4e0bfc4de3eab8558dbfc1c6ac3a6dfa7fac62e26b880744cf70766452",
		Block:     1535952,
	}
	// Caused existing nodes to crash, TX: 936e2b5a9d1837e913f348677505fdc284630ee0a2d7e9ef53d51df7a7ce1126
	s3 := `alexandria-media-multipart(5,6,FD6qwMcfpnsKmoL2kJSfp1czBMVicmkK1Q,undefined,undefined,):"scale":"1000:1","sugTip":[5,50,100],"tokens":{"mtmcollector":"","mtmproducer":"","happybirthdayep":"","early":"","ltbcoin":"","btc":"1GMMg2J5iUKnDf5PbRr9TcKV3R6KfUiB55"}}},"signature":"{ \"success\": true, \"message\": \"IOS7pWaCsiI6Hu7u1uh+JYOW5ZsoiohMSUW448h3qbHHbMr+`

	cases := []struct {
		in    string
		out   MediaMultipartSingle
		txid  string
		block int
		err   error
	}{
		{s, mms, mms.Txid, mms.Block, nil},
		{s2, mms2, mms2.Txid, mms2.Block, nil},
		{s3, MediaMultipartSingle{}, "", 0, ErrNoSignatureEnd},
	}

	for i, c := range cases {
		got, err := VerifyMediaMultipartSingle(c.in, c.txid, c.block)
		if err != c.err {
			t.Errorf("VerifyMediaMultipartSingle(#%d) | err == %q, want %q", i, err, c.err)
		}
		if err == nil && !reflect.DeepEqual(got, c.out) {
			t.Errorf("VerifyMediaMultipartSingle(#%d) | got == %q, want %q", i, got, c.out)
		}
	}
}
