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

func TestPublish42_Validate(t *testing.T) {
	old := utility.Testnet()
	utility.SetTestnet(true)
	defer utility.SetTestnet(old)

	a, err := mp.ParseJson(nil, publishArtifactTestJson, "", &flojson.BlockResult{Height: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ppp := a.(oip042.PublishArtifact)
	j, err := json.Marshal(ppp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Publish Artifact:")
	fmt.Println(string(j))
}

const publishArtifactTestJson = `{  
    "oip042":{  
        "publish":{  
            "artifact":{  
                "floAddress":"oNAydz5TjkhdP3RPuu3nEirYQf49Jrzm4S",
                "info":{  
                    "title":"oip-account test",
                    "description":"Required Description!"
                },
                "details":{  

                },
                "storage":{  
                    "network":"IPFS",
                    "files":[  
                        {  
                            "fname":"lhuWVA00Vn.png",
                            "fsize":23591,
                            "type":"Image"
                        }
                    ],
                    "location":"QmQh7uTC5YSinJG2FgWLrd8MYSNtr8G5JGAckR5ARwmyET"
                },
                "payment":{  

                },
                "type":"Image",
                "timestamp":1533955139,
                "signature":"H+TBYnY938lJIAkvX3962RwYkJIJI5PMS1tS/68gc27XC0ilwEuL6+DJJ3Jdnr7wn9eqX63jBZPsTOn8X9NyA9Y="
            }
        }
    }
}`
