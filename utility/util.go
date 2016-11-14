package utility

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/metacoin/foundation"
)

func CheckAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, &FloParams)
	if err != nil {
		return false
	}
	return true
}

func CheckSignature(address string, signature string, message string) bool {
	reply, err := foundation.RPCCall("verifymessage", address, signature, message)
	if err != nil {
		fmt.Println("foundation error: " + err.Error())
		return false
	}
	if reply == true {
		return true
	}
	return false
}

// reference: Cory LaNou, Mar 2 '14 at 15:21, http://stackoverflow.com/a/22129435/2576956
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
