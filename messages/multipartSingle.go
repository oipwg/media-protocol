package messages

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/oipwg/media-protocol/utility"
	"log"
	"strconv"
	"strings"
)

type MediaMultipartSingle struct {
	Part      int
	Max       int
	Reference string
	Address   string
	Signature string
	Data      string
	Txid      string
	Block     int
}

func CheckMediaMultipartComplete(reference string, dbtx *sql.Tx) ([]byte, error) {
	// using the reference tx, check how many different txs we have and determine if we have all transactions
	// if we have a valid media-multipart complete instance, let's return the byte array it consists of
	var ret []byte

	stmtstr := `select part, max, data from media_multipart where active = 1 and complete = 0 and reference = "` + reference + `" order by part asc`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 120")
		log.Fatal(err)
	}

	rows, stmterr := stmt.Query()
	if err != nil {
		fmt.Println("exit 121")
		log.Fatal(stmterr)
	}

	var rowsCount int = 0
	var pmax int
	var fullData string

	for rows.Next() {
		var part int
		var max int
		var data string
		rows.Scan(&part, &max, &data)

		// TODO: require signature verification for multipart messages
		if rowsCount > max {
			return ret, errors.New("too many rows in multipart message - check for reorg/bogus multipart data")
		}
		rowsCount++

		pmax = max
		fullData += data
	}

	if rowsCount != pmax+1 {
		return ret, errors.New("only found " + strconv.Itoa(rowsCount) + "/" + strconv.Itoa(pmax+1) + " multipart messages")
	}

	stmt.Close()
	rows.Close()

	// set complete to 1
	updatestr := `update media_multipart set complete = 1 where reference = "` + reference + `"`
	updatestmt, updateerr := dbtx.Prepare(updatestr)
	if updateerr != nil {
		fmt.Println("exit 122")
		log.Fatal(updateerr)
	}

	_, updatestmterr := updatestmt.Exec()
	if updatestmterr != nil {
		fmt.Println("exit 123")
		log.Fatal(updatestmterr)
	}
	updatestmt.Close()

	return []byte(fullData), nil
}

func StoreMediaMultipartSingle(mms MediaMultipartSingle, dbtx *sql.Tx) {
	// store in database
	stmtstr := `insert into media_multipart (part, max, address, reference, signature, data, txid, block, complete, success, active) values (` + strconv.Itoa(mms.Part) + `, ` + strconv.Itoa(mms.Max) + `, ?, ?, ?, ?, "` + mms.Txid + `", ` + strconv.Itoa(mms.Block) + `, 0, 0, 1)`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 160")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(mms.Address, mms.Reference, mms.Signature, mms.Data)
	if stmterr != nil {
		fmt.Println("exit 161")
		log.Fatal(stmterr)
	}

	stmt.Close()

}

func UpdateMediaMultipartSuccess(reference string, dbtx *sql.Tx) {

	stmtstr := `update media_multipart set success = 1 where reference = "` + reference + `"`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 140")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec()
	if err != nil {
		fmt.Println("exit 141")
		log.Fatal(stmterr)
	}

}

func VerifyMediaMultipartSingle(s string, txid string, block int) (MediaMultipartSingle, error) {
	var ret MediaMultipartSingle
	prefix := "alexandria-media-multipart("
	prefix2 := "oip-mp("

	// check prefix
	checkPrefix := strings.HasPrefix(s, prefix)
	checkPrefix = checkPrefix || strings.HasPrefix(s, prefix2)
	if !checkPrefix {
		return ret, ErrWrongPrefix
	}

	// trim prefix off
	s = strings.TrimPrefix(s, prefix)
	s = strings.TrimPrefix(s, prefix2)

	comChunks := strings.Split(s, "):")
	if len(comChunks) < 2 {
		return ret, errors.New("Malformed multi-part")
	}

	metaString := comChunks[0]
	dataString := strings.Join(comChunks[1:], "):")

	meta := strings.Split(metaString, ",")
	lm := len(meta)
	// 4 if omitting reference, 5 with all fields, 6 if erroneous fluffy-enigma trailing comma
	if lm != 4 && lm != 5 && lm != 6 {
		return ret, errors.New("Malformed multi-part meta")
	}

	// check part and max
	partS := meta[0]
	part, err := strconv.Atoi(partS)
	if err != nil {
		fmt.Println("cannot convert part to int")
		return ret, errors.New("cannot convert part to int")
	}
	maxS := meta[1]
	max, err2 := strconv.Atoi(maxS)
	if err2 != nil {
		fmt.Println("cannot convert max to int")
		return ret, errors.New("cannot convert max to int")
	}

	// get and check address
	address := meta[2]
	if !utility.CheckAddress(address) {
		// fmt.Println("address doesn't check out: \"" + address + "\"")
		return ret, ErrInvalidAddress
	}

	reference := meta[3]
	signature := meta[lm-1]
	if signature == "" {
		// fluffy-enigma for a while appended an erroneous trailing comma
		signature = meta[lm-2]
	}

	// signature pre-image is <part>-<max>-<address>-<txid>-<data>
	// in the case of multipart[0], txid is 64 zeros
	// in the case of multipart[n], where n != 0, txid is the reference txid (from multipart[0])
	preimage := partS + "-" + maxS + "-" + address + "-" + reference + "-" + dataString
	// fmt.Printf("preimage: %v\n", preimage)

	val, _ := utility.CheckSignature(address, signature, preimage)
	if !val {
		// fmt.Println("signature didn't pass checksignature test")
		return ret, ErrBadSignature
	}

	// if part == 0, reference is ignored and may be omitted
	// the local DB will store reference = txid for this transaction after it's submitted
	// in case of a reorg, the publisher must re-publish this multipart message (sorry)
	if part == 0 {
		reference = txid
	}
	// all checks passed, verified!

	//fmt.Printf("data: %v\n", data)
	// fmt.Printf("=== VERIFIED ===\n")
	//fmt.Printf("part: %v\nmax: %v\nreference: %v\naddress: %v\nsignature: %v\ntxid: %v\nblock: %v\n", part, max, reference, address, signature, txid, block)

	ret = MediaMultipartSingle{
		Part:      part,
		Max:       max,
		Reference: reference,
		Address:   address,
		Signature: signature,
		Data:      dataString,
		Txid:      txid,
		Block:     block,
	}

	return ret, nil

}
