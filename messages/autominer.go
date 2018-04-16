package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/utility"
)

const AUTOMINER_ROOT_KEY = "alexandria-autominer"

type Autominer struct {
	FLOAddress string `json:"FLOaddress"`
	BTCAddress string `json:"BTCaddress"`
	Version    int64  `json:"version"`
}
type AlexandriaAutominer struct {
	Autominer Autominer `json:"alexandria-autominer"`
	Signature string    `json:"signature"`
}

func StoreAutominer(am AlexandriaAutominer, dbtx *sql.Tx, txid string, block *flojson.BlockResult) error {
	// store in database
	stmtStr := `insert into autominer (txid, block, blockTime, active, version,` +
		` floAddress, btcAddress, signature) values (?, ?, ?, 1, ?, ?, ?, ?)`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		fmt.Printf("ERROR in StoreAutominer: autominer dbtx didn't prepare correctly: %v", err)
		return err
	}

	_, stmterr := stmt.Exec(txid, block.Height, block.Time, am.Autominer.Version, am.Autominer.FLOAddress, am.Autominer.BTCAddress, am.Signature)
	if stmterr != nil {
		fmt.Printf("ERROR in StoreAutominer: autominer dbtx didn't execute correctly: %v", err)
		return stmterr
	}

	stmt.Close()
	return nil
}

func VerifyAutominer(b []byte, block int) (AlexandriaAutominer, error) {

	//fmt.Printf("starting Verify Autominer routine...\n")

	//s := string(b[:len(b)])
	//fmt.Printf("s: %+v\n", s)

	var am AlexandriaAutominer

	if !utility.Testnet() && block < 2205000 {
		return am, ErrTooEarly
	}

	//if !utility.IsJSON(string(b)) {
	//	return am, ErrNotJSON
	//}

	err := json.Unmarshal(b, &am)
	if err != nil {
		return am, err
	}

	//fmt.Printf("am: %+v\n", am)

	// verify signature was created by this address
	// signature pre-image for autominer is <btcaddress>-<version>
	preImage := am.Autominer.BTCAddress + "-" + strconv.FormatInt(am.Autominer.Version, 10)

	//fmt.Printf("pre-image: %v", preImage)
	//fmt.Printf("\n\n\n")
	sigOK, _ := utility.CheckSignature(am.Autominer.FLOAddress, am.Signature, preImage)
	if sigOK == false {
		return am, ErrBadSignature
	}

	// fmt.Println(" -- VERIFIED --")
	return am, nil
}
