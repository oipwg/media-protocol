package utility

import (
	"encoding/json"
	"github.com/bitspill/bitsig-go"
	"github.com/btcsuite/btcutil"
)

var utilIsTestnet bool = false

func SetTestnet(testnet bool) {
	utilIsTestnet = testnet
}

func Testnet() bool {
	return utilIsTestnet
}

func CheckAddress(address string) bool {
	var err error
	if utilIsTestnet {
		_, err = btcutil.DecodeAddress(address, &FloTestnetParams)
	} else {
		_, err = btcutil.DecodeAddress(address, &FloParams)
	}
	if err != nil {
		return false
	}
	return true
}

func CheckSignature(address string, signature string, message string) (bool, error) {
	if utilIsTestnet {
		return bitsig_go.CheckSignature(address, signature, message, "testFlo", &FloTestnetParams)
	}
	return bitsig_go.CheckSignature(address, signature, message, "flo", &FloParams)
}

// reference: Cory LaNou, Mar 2 '14 at 15:21, http://stackoverflow.com/a/22129435/2576956
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
