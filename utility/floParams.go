package utility

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// ToDo: Find more appropriate location for this file

// Define some of the required parameters for a Florin coin network.
var FloParams = chaincfg.Params{
	Name:             "flo",
	Net:              0xfdc0a5f1,
	PubKeyHashAddrID: 0x23,
	ScriptHashAddrID: 0x08,
}

var FloTestnetParams = chaincfg.Params{
	Name:             "testFlo",
	Net:              0xfdc05af2,
	PubKeyHashAddrID: 115,
	ScriptHashAddrID: 198,
}

func init() {
	chaincfg.Register(&FloParams)
	chaincfg.Register(&FloTestnetParams)
}
