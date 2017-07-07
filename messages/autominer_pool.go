package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/utility"
)

const AUTOMINER_POOL_ROOT_KEY = "alexandria-autominer-pool"

type AutominerPool struct {
	FLOAddress   string  `json:"FLOaddress"`
	WebURL       string  `json:"WebURL"`
	TargetMargin float64 `json:"TargetMargin"`
	PoolShare    float64 `json:"PoolShare"`
	Version      int64   `json:"version"`
	PoolName     string  `json:"PoolName"` // optional
}
type AlexandriaAutominerPool struct {
	AutominerPool AutominerPool `json:"alexandria-autominer-pool"`
	Signature     string        `json:"signature"`
}

// Verify that the blockchain message is valid for autominer-pool
func VerifyAutominerPool(b []byte, block int) (AlexandriaAutominerPool, error) {

	//fmt.Printf("starting Verify AutominerPool routine...\n")

	//s := string(b[:len(b)])
	//fmt.Printf("s: %+v\n", s)

	var amp AlexandriaAutominerPool

	if !utility.IsJSON(string(b)) {
		return amp, ErrNotJSON
	}

	err := json.Unmarshal(b, &amp)
	if err != nil {
		return amp, err
	}

	fmt.Printf("amp: %+v\n", amp)

	// parse float into string for signature
	targetMarginStr := strconv.FormatFloat(amp.AutominerPool.TargetMargin, 'f', -1, 64)
	poolShareStr := strconv.FormatFloat(amp.AutominerPool.PoolShare, 'f', -1, 64)
	versionStr := strconv.FormatInt(amp.AutominerPool.Version, 10)

	fmt.Printf("preimage components")
	fmt.Printf("targetMarginStr: %v\npoolShareStr: %v\nversionStr: %v\n", targetMarginStr, poolShareStr, versionStr)

	// verify signature was created by this address
	// signature pre-image for autominer is <weburl>-<version>-<targetmargin>-<poolshare>-[poolname]
	// poolname is optional
	preImage := amp.AutominerPool.WebURL + "-" + versionStr + "-" + targetMarginStr + "-" + poolShareStr
	if len(amp.AutominerPool.PoolName) > 0 {
		preImage += "-" + amp.AutominerPool.PoolName
	}

	fmt.Printf("pre-image: %v", preImage)
	fmt.Printf("\n\n\n")
	sigOK, _ := utility.CheckSignature(amp.AutominerPool.FLOAddress, amp.Signature, preImage)
	if sigOK == false {
		return amp, ErrBadSignature
	}

	fmt.Println(" -- VERIFIED --")
	return amp, nil
}

func StoreAutominerPool(amp AlexandriaAutominerPool, dbtx *sql.Tx, txid string, block *flojson.BlockResult) error {
	// store in database
	stmtStr := `insert into autominer_pool (txid, block, blockTime, active, version,` +
		` floAddress, webURL, targetMargin, poolShare, poolName, signature) values (?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		fmt.Printf("ERROR in StoreAutominerPool: autominer-pool dbtx didn't prepare correctly: %v", err)
		return err
	}

	poolName := amp.AutominerPool.PoolName
	// TODO: does empty string always === ""?

	_, stmterr := stmt.Exec(txid, block.Height, block.Time, amp.AutominerPool.Version, amp.AutominerPool.FLOAddress, amp.AutominerPool.WebURL, amp.AutominerPool.TargetMargin, amp.AutominerPool.PoolShare, poolName, amp.Signature)
	if stmterr != nil {
		fmt.Printf("ERROR in StoreAutominerPool: autominer-pool dbtx didn't execute correctly: %v", err)
		return stmterr
	}

	stmt.Close()
	return nil
}
