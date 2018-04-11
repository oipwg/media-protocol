package oip042_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/metacoin/flojson"
	mp "github.com/oipwg/media-protocol"
	"github.com/oipwg/media-protocol/oip042"
	"github.com/oipwg/media-protocol/utility"
)

func TestPublishPropertyParty_Validate(t *testing.T) {
	old := utility.Testnet()
	utility.SetTestnet(true)
	defer utility.SetTestnet(old)

	a, err := mp.ParseJson(nil, samplePropertyPartyJson, "", &flojson.BlockResult{Height: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ppp := a.(oip042.PublishPropertyParty)
	j, err := json.Marshal(ppp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Publish Proprty-Party:")
	fmt.Println(string(j))
}

const samplePropertyPartyJson = `{
  "oip042": {
    "publish": {
      "artifact": {
        "floAddress": "oVbNpajdAv31g668HfHP6mgVfo9VpL7YAq",
        "timestamp": 1523388623,
        "type": "property",
        "subtype": "party",
        "details": {
          "ns": "DS",
          "partyRole": "association",
          "partyType": "group",
          "members": [
            "c47e8840f19f20406590799923bf3315a8ad976f2b510f23991aebeca32f6ff8",
            "2118f4ded6b0fa80c6c04cfb5b4e0828dd75cb0b45d64bd7fae1b3197a3e5372"
          ],
          "attrs": []
        },
        "info": {
          "title": "Test 42 Property-Party",
          "description": "Test 42 Property-Party description",
          "tags": "party,test,42"
        },
        "storage": {
          "network": "ipfs",
          "location": "QmUmu8ArQPm8nYN1GqT1nCyEkMyWnupfmp7YZT2DBKoiJ3",
          "files": [
            {
              "fName": "document_104.txt",
              "fType": "text/plain",
              "fSize": 129,
              "dName": "Tenure Entity /ds-property/tmp-tenure-24617"
            }
          ]
        },
        "signature": "IMt9EcnLJCkTzb7YbPBdOe/NX7+hwzxY5hOsjXCG02/8E3+b0uL/nhJ9p4b6EbKmWFPM9HXORAXAK82Wy8AnnvY="
      }
    }
  }
}`
