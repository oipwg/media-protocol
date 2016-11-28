package messages

import (
	"encoding/json"
	"errors"
	"github.com/dloa/media-protocol/utility"
	"regexp"
	"strings"
)

const DEACTIVATION_ROOT_KEY = "alexandria-deactivation"

type AlexandriaDeactivation struct {
	AlexandriaDeactivation struct {
		// required txid and address fields
		Txid    string `json:"txid"`
		Address string `json:"address"`
	} `json:"alexandria-deactivation"`
	Signature string `json:"signature"`
}

func VerifyDeactivation(b []byte) (AlexandriaDeactivation, error) {

	var v AlexandriaDeactivation
	var i interface{}
	var m map[string]interface{}

	if !strings.HasPrefix(string(b), `{"alexandria-deactivation"`) {
		return v, errors.New("Not alexandria-deactivation")
	}

	if !utility.IsJSON(string(b)) {
		return v, ErrNotJSON
	}

	err := json.Unmarshal(b, &v)
	if err != nil {
		return v, err
	}

	errr := json.Unmarshal(b, &i)
	if errr != nil {
		return v, err
	}

	m = i.(map[string]interface{})
	var signature string

	// check the JSON object root key
	// find the signature string
	for key, val := range m {
		if key == "signature" {
			signature = val.(string)
		} else {
			if key != DEACTIVATION_ROOT_KEY {
				return v, errors.New("can't verify deactivation - JSON object root key doesn't match accepted value")
			}
		}
	}

	// verify txid
	rt := regexp.MustCompile("^[a-fA-F0-9]*$")
	if !rt.MatchString(v.AlexandriaDeactivation.Txid) || len(v.AlexandriaDeactivation.Txid) != 64 {
		return v, errors.New("can't verify deactivation - txid in incorrect format")
	}

	// verify address
	ra := regexp.MustCompile("^[a-zA-Z0-9]*$")
	if !ra.MatchString(v.AlexandriaDeactivation.Address) || len(v.AlexandriaDeactivation.Address) != 34 {
		return v, errors.New("can't verify deactivation - address in incorrect format")
	}

	// verify signature
	if v.Signature != signature {
		return v, ErrBadSignature
	}

	// verify signature was created by this address
	// signature pre-image for deactivation is <address>-<txid>
	val, _ := utility.CheckSignature(v.AlexandriaDeactivation.Address, signature, v.AlexandriaDeactivation.Address+"-"+v.AlexandriaDeactivation.Txid)
	if val == false {
		return v, ErrBadSignature
	}

	return v, nil

}
