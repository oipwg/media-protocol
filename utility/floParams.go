package utility

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// ToDo: Find more appropriate location for this file

// Define some of the required parameters for a Florin coin network.
var FloParams = chaincfg.Params{
	Name:             "flo",
	PubKeyHashAddrID: 0x23,
	ScriptHashAddrID: 0x08,
}

var FloTestnetParams = chaincfg.Params{
	Name:             "testFlo",
	PubKeyHashAddrID: 0x73,
	ScriptHashAddrID: 0xC6,
}

func init() {
	chaincfg.Register(&FloParams)
	chaincfg.Register(&FloTestnetParams)
}
