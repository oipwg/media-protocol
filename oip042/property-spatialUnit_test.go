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

func TestPublishPropertySpatialUnit_Validate(t *testing.T) {
	old := utility.Testnet()
	utility.SetTestnet(true)
	defer utility.SetTestnet(old)

	a, err := mp.ParseJson(nil, samplePropertySpatialUnitJson, "", &flojson.BlockResult{Height: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ppp := a.(oip042.PublishPropertySpatialUnit)
	j, err := json.Marshal(ppp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Publish Proprty-SpatialUnit:")
	fmt.Println(string(j))
}

const samplePropertySpatialUnitJson = `{
  "oip042": {
    "publish": {
      "artifact": {
        "floAddress": "oVbNpajdAv31g668HfHP6mgVfo9VpL7YAq",
        "timestamp": 1523388623,
        "type": "property",
        "subtype": "spatialUnit",
        "details": {
          "ns": "DS",
          "geometry": {
            "type": "text",
            "data": "Bounded by the alpha and beta rivers, iron hills to flame gulch"
          },
          "spatialType": "text",
          "attrs": []
        },
        "info": {
          "title": "Test 42 Property-SpatialUnit",
          "description": "Test 42 Property-SpatialUnit description",
          "tags": "spatialUnit,test,42"
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
