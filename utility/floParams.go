package utility

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// ToDo: Find more appropriate location for this file

// Define some of the required parameters for a Florin coin network.
var FloParams = chaincfg.Params{
	Name:             "Florincoin",
	PubKeyHashAddrID: 0x23,
	ScriptHashAddrID: 0x08,
}

func init() {
	chaincfg.Register(&FloParams)
}
