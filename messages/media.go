package messages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dloa/media-protocol/utility"
	"log"
	"strconv"
	"strings"
)

const MEDIA_ROOT_KEY = "alexandria-media"

type AlexandriaMedia struct {
	AlexandriaMedia struct {
		// required media metadata
		Torrent   string `json:"torrent"`
		Publisher string `json:"publisher"`
		Timestamp int64  `json:"timestamp"`
		Type      string `json:"type"`

		Info struct {
			// required file information
			Title       string `json:"title"`
			Description string `json:"description"`
			Year        int    `json:"year"`

			// optional extra-info field
			ExtraInfo interface{} `json:"extra-info"`
		} `json:"info"`

		// optional fields
		Payment interface{} `json:"payment"`
		Extras  string      `json:"extras"`
	} `json:"alexandria-media"`
	Signature string `json:"signature"`
}

func extractMediaExtraInfo(jmap map[string]interface{}) ([]byte, error) {
	// find the "extra info" json object
	var ret []byte
	for k, v := range jmap {
		if k == "alexandria-media" {
			vm := v.(map[string]interface{})
			for k2, v2 := range vm {
				if k2 == "info" {
					v2m := v2.(map[string]interface{})
					for k3, v3 := range v2m {
						if k3 == "extra-info" {
							// fmt.Printf("v3(%v): %v\n\n", reflect.TypeOf(v3), v3)
							v3json, err := json.Marshal(v3)
							if err != nil {
								return ret, err
							}
							return v3json, nil
						}
					}

				}
			}
		}
	}
	return ret, errors.New("no media extra info found")
}

func extractMediaPayment(jmap map[string]interface{}) ([]byte, error) {
	// find the "payment" json object
	var ret []byte
	for k, v := range jmap {
		if k == "alexandria-media" {
			vm := v.(map[string]interface{})
			for k2, v2 := range vm {
				if k2 == "payment" {
					// fmt.Printf("v3(%v): %v\n\n", reflect.TypeOf(v3), v3)
					v2json, err := json.Marshal(v2)
					if err != nil {
						return ret, err
					}
					return v2json, nil

				}
			}
		}
	}
	return ret, errors.New("no payment extra info found")
}

func DeactivateMedia(deactiv AlexandriaDeactivation, dbtx *sql.Tx) error {
	stmtstr := `update media set invalidated = 1 where publisher = "` + deactiv.AlexandriaDeactivation.Address + `" and txid = "` + deactiv.AlexandriaDeactivation.Txid + `"`
	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		return err
	}

	_, stmterr := stmt.Exec()
	if stmterr != nil {
		return stmterr
	}

	stmt.Close()
	return nil
}

func StoreMedia(media AlexandriaMedia, jmap map[string]interface{}, dbtx *sql.Tx, txid string, block int, multipart int) {
	// check for media payment data
	payment, payment_err := extractMediaPayment(jmap)
	paymentString := ""
	if payment_err != nil {
		fmt.Printf("payment data not found/failed - error returned: %v\n", payment_err)
	} else {
		paymentString = string(payment)
	}

	// check for media info extras
	extraInfo, ei_err := extractMediaExtraInfo(jmap)
	extraInfoString := ""
	if ei_err != nil {
		fmt.Printf("extra info not found/failed - error returned: %v\n", ei_err)
	} else {
		extraInfoString = string(extraInfo)
	}

	// make sure extras is stored as an empty string if it doesn't exist
	if len(media.AlexandriaMedia.Extras) < 1 {
		media.AlexandriaMedia.Extras = ""
	}

	stmtstr := `insert into media (publisher, torrent, timestamp, type, info_title, info_description, info_year, info_extra, payment, extras, txid, block, signature, multipart, active) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "` + txid + `", ` + strconv.Itoa(block) + `, ?, ` + strconv.Itoa(multipart) + `, 1)`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 102")
		log.Fatal(err)
	}

	// fmt.Printf("stmt: %v\n", stmt)

	_, stmterr := stmt.Exec(media.AlexandriaMedia.Publisher, media.AlexandriaMedia.Torrent, media.AlexandriaMedia.Timestamp, media.AlexandriaMedia.Type, media.AlexandriaMedia.Info.Title, media.AlexandriaMedia.Info.Description, media.AlexandriaMedia.Info.Year, extraInfoString, paymentString, media.AlexandriaMedia.Extras, media.Signature)
	if stmterr != nil {
		fmt.Println("exit 103")
		log.Fatal(stmterr)
	}

	stmt.Close()

}

func VerifyMedia(b []byte) (AlexandriaMedia, map[string]interface{}, error) {

	var v AlexandriaMedia
	var i interface{}
	var m map[string]interface{}

	if !strings.HasPrefix(string(b), `{ "alexandria-media"`) {
		return v, nil, errors.New("Not alexandria-media")
	}

	if !utility.IsJSON(string(b)) {
		return v, m, ErrNotJSON
	}

	// fmt.Printf("Attempting to verify alexandria-media JSON...")
	err := json.Unmarshal(b, &v)
	if err != nil {
		return v, m, err
	}

	errr := json.Unmarshal(b, &i)
	if errr != nil {
		return v, m, err
	}

	m = i.(map[string]interface{})
	var signature string

	// check the JSON object root key
	// find the signature string
	for key, val := range m {
		if key == "signature" {
			signature = val.(string)
		} else {
			if key != MEDIA_ROOT_KEY {
				return v, m, errors.New("can't verify media - JSON object root key doesn't match accepted value")
			}
		}
	}

	// fmt.Printf("*** debug: JSON object root matches, printing v:\n%v\n*** /debug ***\n", v)
	// verify torrent hash length
	if len(v.AlexandriaMedia.Torrent) <= 1 {
		return v, m, errors.New("can't verify media - invalid torrent hash length")
	}

	// verify signature
	if v.Signature != signature {
		return v, m, errors.New("can't verify media - signature mismatch")
	}

	// verify timestamp length
	if v.AlexandriaMedia.Timestamp <= 0 {
		return v, m, errors.New("can't verify media - invalid timestamp")
	}

	// verify type length
	if len(v.AlexandriaMedia.Type) <= 1 {
		return v, m, errors.New("can't verify media - invalid type length")
	}

	// verify media info lengths
	if len(v.AlexandriaMedia.Info.Title) <= 0 {
		return v, m, errors.New("can't verify media - invalid info title length")
	}
	if len(v.AlexandriaMedia.Info.Description) <= 0 {
		return v, m, errors.New("can't verify media - invalid info description length")
	}
	if v.AlexandriaMedia.Info.Year <= 0 {
		return v, m, errors.New("can't verify media - invalid info year")
	}

	// verify signature was created by this address
	// signature pre-image for media is <torrenthash>-<publisher>-<timestamp>
	val, _ := utility.CheckSignature(v.AlexandriaMedia.Publisher, signature, v.AlexandriaMedia.Torrent+"-"+v.AlexandriaMedia.Publisher+"-"+strconv.FormatInt(v.AlexandriaMedia.Timestamp, 10))
	if val == false {
		return v, m, ErrBadSignature
	}

	// fmt.Println(" -- VERIFIED --")
	return v, m, nil

}
