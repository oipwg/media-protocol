package oip042_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/metacoin/flojson"
	mp "github.com/oipwg/media-protocol"
	"github.com/oipwg/media-protocol/oip042"
)

func TestPublishTomogram_Validate(t *testing.T) {
	a, err := mp.ParseJson(nil, sampleJson, "", &flojson.BlockResult{Height: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	pt := a.(oip042.PublishTomogram)
	j, err := json.Marshal(pt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Publish Research-Tomogram:")
	fmt.Println(string(j))
}

const sampleJson = `{
  "oip042": {
    "publish": {
      "artifact": {
        "floAddress": "FNRCCaR7Y4T4oY5KmUMgjRULsMp7uh6uZY",
        "timestamp": 15,
        "type": "research",
        "subtype": "tomogram",
        "details": {
          "scopeName": "string",
          "speciesName": "string",
          "tiltSingleDual": 1,
          "strain": "string",
          "swAcquisition": "string",
          "swProcess": "string",
          "emdb": "string",
          "magnification": 1.0,
          "defocus": 1.0,
          "date": 13,
          "NBCItaxID": 14,
          "artNotes": "string"
        },
        "info": {
          "title": "string",
          "description": "string",
          "tags": "comma delimited list of search terms",
          "extraInfo": {},
          "roles": [
            {
              "party": "tenure"
            },
            {
              "party": "tenure"
            }
          ]
        },
        "storage": {
          "network": "ipfs",
          "location": "string",
          "files": [
            {
              "fName": "string",
              "fType": "string",
              "fSize": 12,
              "fNotes": "string",
              "dName": "string"
            }
          ]
        },
        "signature": "ILKzRqAhzx8KnyBLJcy/oY4g4gG9rcENPiSjw0jy0BlVTvU42HNz0Bewfqxq0DZaMjB9qskZYgs1RZ0hiVYcgdg="
      }
    }
  }
}`
