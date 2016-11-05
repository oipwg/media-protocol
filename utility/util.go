package utility

import (
	"encoding/json"
	"fmt"
	"github.com/metacoin/flojson"
	"github.com/metacoin/foundation"
)

// reference: Cory LaNou, Mar 2 '14 at 15:21, http://stackoverflow.com/a/22129435/2576956
func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func CheckAddress(address string) bool {
	reply, err := foundation.RPCCall("validateaddress", address)
	if err != nil {
		fmt.Println("foundation error: " + err.Error())
		return false
	}
	result, ok := reply.(*flojson.ValidateAddressResult)
	if !ok {
		return false
	}
	if result.IsValid == true {
		return true
	}
	return false
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
