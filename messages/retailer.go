package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/utility"
)

const RETAILER_ROOT_KEY = "alexandria-retailer"

type Retailer_OptionalFields struct {
	Name         string  `json:"name"`
	MinimumShare float64 `json:"minimum-share"`
}
type Retailer struct {
	FLOAddress     string                  `json:"FLOaddress"`
	BTCAddress     string                  `json:"BTCaddress"`
	WebURL         string                  `json:"web-url"`
	Version        int64                   `json:"version"`
	OptionalFields Retailer_OptionalFields `json:"optional-fields"`
}
type AlexandriaRetailer struct {
	Retailer  Retailer `json:"alexandria-retailer"`
	Signature string   `json:"signature"`
}

func StoreRetailer(ar AlexandriaRetailer, dbtx *sql.Tx, txid string, block *flojson.BlockResult) error {
	// store in database
	stmtStr := `insert into retailer (txid, block, blockTime, active, version,` +
		` floAddress, btcAddress, webURL, signature) values (?, ?, ?, 1, ?, ?, ?, ?, ?)`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		fmt.Printf("ERROR in StoreRetailer: retailer dbtx didn't prepare correctly: %v\n", err)
		return err
	}

	res, stmterr := stmt.Exec(txid, block.Height, block.Time, ar.Retailer.Version, ar.Retailer.FLOAddress, ar.Retailer.BTCAddress, ar.Retailer.WebURL, ar.Signature)
	if stmterr != nil {
		fmt.Printf("ERROR in StoreRetailer: retailer dbtx didn't execute correctly: %v", err)
		return stmterr
	}

	// optional fields
	if ar.Retailer.OptionalFields != (Retailer_OptionalFields{}) {
		lastId, err := res.LastInsertId()
		if err != nil {
			fmt.Printf("ERROR in StoreRetailer: couldn't get LastInsertId() from res: %v\n", err)
			return err
		}
		err = StoreRetailerOptionalFields(ar.Retailer.OptionalFields, dbtx, txid, block, lastId)
		if err != nil {
			return err
		}
	}

	stmt.Close()
	return nil
}

func StoreRetailerOptionalFields(rof Retailer_OptionalFields, dbtx *sql.Tx, txid string, block *flojson.BlockResult, id int64) error {
	stmtStr := `insert into retailer_optionalfields (retailer_uid, name, minimumShare) values (?, ?, ?)`
	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		fmt.Printf("ERROR in StoreRetailerOptionalFields: retailer_optionalfields dbtx didn't prepare correctly: %v\n", err)
		return err
	}

	_, stmterr := stmt.Exec(id, rof.Name, rof.MinimumShare)
	if stmterr != nil {
		fmt.Printf("ERROR in StoreRetailerOptionalFields: retailer_optionalfields dbtx didn't execute correctly: %v\n", err)
		return stmterr
	}
	stmt.Close()
	return nil
}

func VerifyRetailer(b []byte, block int) (AlexandriaRetailer, error) {

	//fmt.Printf("starting VerifyRetailer routine...\n")

	//s := string(b[:len(b)])
	//fmt.Printf("s: %+v\n", s)

	var ar AlexandriaRetailer

	if !utility.IsJSON(string(b)) {
		return ar, ErrNotJSON
	}

	err := json.Unmarshal(b, &ar)
	if err != nil {
		return ar, err
	}

	fmt.Printf("ar: %+v\n", ar)

	// verify signature was created by this address
	// signature pre-image for retailer is <btcaddress>-<weburl>-<version>
	preImage := ar.Retailer.BTCAddress + "-" + ar.Retailer.WebURL + "-" + strconv.FormatInt(ar.Retailer.Version, 10)

	// if there are optional fields, they are included in order of sequence
	// example: <btcaddress>-<weburl>-<version>-<name>-<minimumshare>-<name>-...
	if ar.Retailer.OptionalFields != (Retailer_OptionalFields{}) {
		preImage += "-" + ar.Retailer.OptionalFields.Name + "-" + strconv.FormatFloat(ar.Retailer.OptionalFields.MinimumShare, 'f', -1, 64)
	}

	fmt.Printf("\n###### pre-image: %v", preImage)
	fmt.Printf("\n###### signature: %v\n", ar.Signature)
	fmt.Printf("\n\n\n")
	sigOK, _ := utility.CheckSignature(ar.Retailer.FLOAddress, ar.Signature, preImage)
	if sigOK == false {
		return ar, ErrBadSignature
	}

	// fmt.Println(" -- VERIFIED --")
	return ar, nil
}
