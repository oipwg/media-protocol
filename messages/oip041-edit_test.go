package messages

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestUnSquashPatch(t *testing.T) {
	uj := UnSquashPatch(squash_json)

	var expected interface{}
	var actual interface{}
	json.Unmarshal([]byte(unsquash_json), &expected)
	err := json.Unmarshal([]byte(uj), &actual)
	if err != nil {
		t.Errorf("Error during unsquash: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Improper unsquash:\n%s\n", uj)
	}
}

var unsquash_json = `[{"op":"add","path":"/artifact/info/extraInfo","value":{"artist":"Adam B. Levine","composers":["Adam B. Levine"]}},{"op":"replace","path":"/artifact/timestamp","value":1481420000},{"op":"replace","path":"/artifact/storage/files/0/fname","value":"1 - Skipping Stones.mp3"},{"op":"replace","path":"/artifact/storage/files/0/dname","value":"Skipping of the Stones"},{"op":"replace","path":"/artifact/info/title","value":"Title Change Test"},{"op":"remove","path":"/artifact/txid"},{"op":"remove","path":"/artifact/oip-041/signature"},{"op":"remove","path":"/artifact/storage/files/0/sugPlay"},{"op":"remove","path":"/artifact/storage/files/0/sugBuy"},{"op":"remove","path":"/artifact/storage/files/0/storage"},{"op":"remove","path":"/artifact/storage/files/0/retail"},{"op":"remove","path":"/artifact/storage/files/0/promo"},{"op":"remove","path":"/artifact/storage/files/0/minPlay"},{"op":"remove","path":"/artifact/storage/files/0/minBuy"},{"op":"remove","path":"/artifact/storage/files/0/disallowPlay"},{"op":"remove","path":"/artifact/storage/files/0/disallowBuy"},{"op":"remove","path":"/artifact/payment/tokens"},{"op":"remove","path":"/artifact/payment/sug_tip"},{"op":"remove","path":"/artifact/payment/scale"},{"op":"remove","path":"/artifact/payment/fiat"},{"op":"remove","path":"/artifact/info/extra-info"}]`

var squash_json = `{
  "add": [
    {
      "path": "/artifact/info/extraInfo",
      "value": {
        "artist": "Adam B. Levine",
        "composers": [
          "Adam B. Levine"
        ]
      }
    }
  ],
  "replace": [
    {
      "path": "/artifact/timestamp",
      "value": 1481420000
    },
    {
      "path": "/artifact/storage/files/0/fname",
      "value": "1 - Skipping Stones.mp3"
    },
    {
      "path": "/artifact/storage/files/0/dname",
      "value": "Skipping of the Stones"
    },
    {
      "path": "/artifact/info/title",
      "value": "Title Change Test"
    }
  ],
  "remove": [
    {
      "path": "/artifact/txid"
    },
    {
      "path": "/artifact/oip-041/signature"
    },
    {
      "path": "/artifact/storage/files/0/sugPlay"
    },
    {
      "path": "/artifact/storage/files/0/sugBuy"
    },
    {
      "path": "/artifact/storage/files/0/storage"
    },
    {
      "path": "/artifact/storage/files/0/retail"
    },
    {
      "path": "/artifact/storage/files/0/promo"
    },
    {
      "path": "/artifact/storage/files/0/minPlay"
    },
    {
      "path": "/artifact/storage/files/0/minBuy"
    },
    {
      "path": "/artifact/storage/files/0/disallowPlay"
    },
    {
      "path": "/artifact/storage/files/0/disallowBuy"
    },
    {
      "path": "/artifact/payment/tokens"
    },
    {
      "path": "/artifact/payment/sug_tip"
    },
    {
      "path": "/artifact/payment/scale"
    },
    {
      "path": "/artifact/payment/fiat"
    },
    {
      "path": "/artifact/info/extra-info"
    }
  ]
}`

func TestHandleOIP041Edit(t *testing.T) {
	t.Skip("Needs a test DB")

	o, err := VerifyOIP041(example_edit, 21000000)
	if err != nil {
		fmt.Println(err)
	}
	HandleOIP041Edit(o, o.Edit.TxID, 21000000, nil)
	fmt.Println(o)
}

var example_edit = `{
    "oip-041":{
        "editArtifact":{
            "txid":"$artifactID",
            "timestamp":1234567890,
            "patch":{
                "add":[
                    {
                        "path":"/payment/tokens/mtcproducer",
                        "value":""
                    }
                ],
                "replace":[
                    {
                        "path":"/storage/files/3/fname",
                        "value":"birthdayepFirst.jpg"
                    },
                    {
                        "path":"/storage/files/3/dname",
                        "value":"Cover Art 2"
                    },
                    {
                        "path":"/info/title",
                        "value":"Happy Birthday"
                    },
                    {
                        "path":"/timestamp",
                        "value":1481420001
                    }
                ],
                "remove":[
                    {
                        "path":"/payment/tokens/mtmproducer"
                    },
                    {
                        "path":"/storage/files/0/sugBuy"
                    }
                ]
            }
        },
    	"signature":"$txid-$MD5HashOfPatch-$timestamp"
    }
}`
