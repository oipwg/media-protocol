package utility

import (
	"encoding/json"
	"github.com/bitspill/bitsig-go"
	"github.com/btcsuite/btcutil"
)

func CheckAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, &FloParams)
	if err != nil {
		return false
	}
	return true
}

func CheckSignature(address string, signature string, message string) (bool, error) {
	return bitsig_go.CheckSignature(address, signature, message, "Florincoin", &FloParams)
}

// reference: Cory LaNou, Mar 2 '14 at 15:21, http://stackoverflow.com/a/22129435/2576956
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
