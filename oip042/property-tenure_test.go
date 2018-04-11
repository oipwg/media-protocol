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

func TestPublishPropertyTenure_Validate(t *testing.T) {
	old := utility.Testnet()
	utility.SetTestnet(true)
	defer utility.SetTestnet(old)

	a, err := mp.ParseJson(nil, samplePropertyTenureJson, "", &flojson.BlockResult{Height: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ppp := a.(oip042.PublishPropertyTenure)
	j, err := json.Marshal(ppp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Publish Proprty-Tenure:")
	fmt.Println(string(j))
}

const samplePropertyTenureJson = `{
  "oip042": {
    "publish": {
      "artifact": {
        "floAddress": "oVbNpajdAv31g668HfHP6mgVfo9VpL7YAq",
        "timestamp": 1523388623,
        "type": "property",
        "subtype": "tenure",
        "details": {
          "ns": "DS",
          "tenureType": "freehold",
          "party": "418bbb72d034ab53ae616870b2ad554b279179323f54f1864780b71535d88e4b",
          "spatialUnit": "bb23d4f9fd0fa0f6e16b9ce59538f8c8cb9ca4c9c2c2992054cd1435de6a03e4",
          "attrs": []
        },
        "info": {
          "title": "Test 42 Property-Tenure",
          "description": "Test 42 Property-Tenure description",
          "tags": "tenure,test,42"
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
