package messages

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var oip041_music_example = `{
  "oip-041": {
    "artifact": {
      "publisher": "F97Tp8LYnw94CpXmAhqACXWTT36jyvLCWx",
      "timestamp": 1470269387,
      "type": "music",
      "storage":{
        "network": "IPFS",
        "location": "QmPukCZKeJD4KZFtstpvrguLaq94rsWfBxLU1QoZxvgRxA"
      },
      "files": [
	    {
	      "dname": "Skipping Stones",
	      "fname": "1 - Skipping Stones.mp3",
	      "fsize": 6515667,
	      "type": "album track",
	      "duration": 1533.603293,
	      "sugPlay": 100,
	      "minPlay": null,
	      "sugBuy": 750,
	      "minBuy": 500,
	      "promo": 10,
	      "retail": 15,
	      "ptpFT": 10,
	      "ptpDT": 20,
	      "ptpDA": 50
	    },
	    {
	      "dname": "Lessons",
	      "fname": "2 - Lessons with intro.mp3",
	      "fsize": 6515667,
	      "type": "album track",
	      "duration": 1231.155243,
	      "disallowPlay": 1,
	      "sugBuy": 750,
	      "minBuy": 500,
	      "promo": 10,
	      "retail": 15,
	      "ptpFT": 10,
	      "ptpDT": 20,
	      "ptpDA": 50
	    },
	    {
	      "dname": "Born to Roam",
	      "fname": "3 - Born to Roam.mp3",
	      "fsize": 6515667,
	      "type": "album track",
	      "duration": 2374.550714,
	      "sugPlay": 100,
	      "minPlay": 50,
	      "disallowBuy": 1,
	      "promo": 10,
	      "retail": 15,
	      "ptpFT": 10,
	      "ptpDT": 20,
	      "ptpDA": 50
	    },
	    {
	      "dname": "Cover Art",
	      "fname": "birthdayepFINAL.jpg",
	      "type": "coverArt",
	      "disallowBuy": 1
	    }
	  ],
      "info": {
        "title": "Happy Birthday EP",
        "description": "this is the second organically grown, gluten free album released by Adam B. Levine - contact adam@tokenly.com with questions or comments or discuss collaborations.",
        "year": 2016,
        "extra-info": {
          "artist": "Adam B. Levine",
          "company": "",
          "composers": [
            "Adam B. Levine"
          ],
          "copyright": "",
          "tokenly_ID": "",
          "usageProhibitions": "",
          "usageRights": "",
          "tags": []
        }
      },
      "payment": {
        "fiat": "USD",
        "scale": "1000:1",
        "sug_tip": [
          5,
          50,
          100
        ],
        "tokens": {
          "MTMCOLLECTOR": "",
          "MTMPRODUCER": "",
          "HAPPYBDAYEP": "",
          "EARLY": "",
          "LTBCOIN": "",
          "BTC": "1GMMg2J5iUKnDf5PbRr9TcKV3R6KfUiB55"
        }
      }
    },
    "signature": "H27r7UxUb8BozjEvV0v++nCyRI7S6yyroeKCJQpgU5NO3CP6FpXWs5kCxy8vhmMhbtpj/FMj+8s3+updw7g+bmE="
  }
}`

func TestDecodeOIP041(t *testing.T) {
	oip041w := Oip041Wrapper{}
	err := json.Unmarshal([]byte(oip041_music_example), &oip041w)
	fmt.Println(err)
	fmt.Println(oip041w.Oip041)
	if err != nil {
		t.Error("error")
	}
	b, e := json.Marshal(oip041w)
	fmt.Println(e)
	fmt.Println(string(b))
	if !reflect.DeepEqual(oip041w.Oip041, oip041_example_obj) {
		t.Error("not equal")
	}
}

var oip041_example_obj Oip041 = Oip041{
	Oip041Artifact{
		"F97Tp8LYnw94CpXmAhqACXWTT36jyvLCWx",
		1470269387,
		"music",
		Oip041Info{
			"Happy Birthday EP",
			"this is the second organically grown, gluten free album released by Adam B. Levine - contact adam@tokenly.com with questions or comments or discuss collaborations.",
			2016,
			Oip041MusicExtraInfo{
				"Adam B. Levine",
				"",
				[]string{
					"Adam B. Levine"},
				"",
				"",
				"",
				[]string{}},
			""},
		Oip041Storage{
			"IPFS",
			"QmPukCZKeJD4KZFtstpvrguLaq94rsWfBxLU1QoZxvgRxA"},
		[]Oip041Files{
			{
				0,
				"Skipping Stones",
				1533.603293,
				"1 - Skipping Stones.mp3",
				6515667,
				0,
				100,
				10,
				15,
				10,
				20,
				50,
				"album track",
				"",
				0,
				500,
				750,
				Oip041Storage{}},
			{
				0,
				"Lessons",
				1231.155243,
				"2 - Lessons with intro.mp3",
				6515667,
				0,
				0,
				10,
				15,
				10,
				20,
				50,
				"album track",
				"",
				0,
				500,
				750,
				Oip041Storage{}},

			{
				1,
				"Born to Roam",
				2374.550714,
				"3 - Born to Roam.mp3",
				6515667,
				50,
				100,
				10,
				15,
				10,
				20,
				50,
				"album track",
				"",
				0,
				0,
				0,
				Oip041Storage{}},
			{
				1,
				"Cover Art",
				0,
				"birthdayepFINAL.jpg",
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				"coverArt",
				"",
				0,
				0,
				0,
				Oip041Storage{}}},

		Oip041Payment{
			"USD",
			"1000:1",
			[]int{5,
				50,
				100},
			map[string]string{
				"BTC":          "1GMMg2J5iUKnDf5PbRr9TcKV3R6KfUiB55",
				"MTMCOLLECTOR": "",
				"MTMPRODUCER":  "",
				"HAPPYBDAYEP":  "",
				"EARLY":        "",
				"LTBCOIN":      ""}}},
	"H27r7UxUb8BozjEvV0v++nCyRI7S6yyroeKCJQpgU5NO3CP6FpXWs5kCxy8vhmMhbtpj/FMj+8s3+updw7g+bmE="}
