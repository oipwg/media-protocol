package messages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oipwg/media-protocol/utility"
	"log"
	"strconv"
	"strings"
)

const PUBLISHER_ROOT_KEY = "alexandria-publisher"

type AlexandriaPublisher struct {
	AlexandriaPublisher struct {
		// required publisher metadata
		Name      string `json:"name"`
		Address   string `json:"address"`
		Timestamp int64  `json:"timestamp"`

		// optional fields
		Emailmd5   string `json:"emailmd5"`
		Bitmessage string `json:"bitmessage"`
	} `json:"alexandria-publisher"`
	Signature string `json:"signature"`
}

func CheckPublisherAddressExists(address string, dbtx *sql.Tx) bool {
	// check if this publisher address is already in-use
	stmtstr := `select name from publisher where address = ?`

	rows, stmterr := dbtx.Query(stmtstr, address)
	if stmterr != nil {
		fmt.Println("exit 91248")
		log.Fatal(stmterr)
	}

	var rowsCount int = 0
	for rows.Next() {
		rowsCount++
	}

	rows.Close()
	return rowsCount > 0

}

func CreateNewPublisherTxComment(b []byte) {
	// given some JSON, post it to the blockchain using either a tx-comment or multipart tx-comment

}

func StorePublisher(publisher AlexandriaPublisher, dbtx *sql.Tx, txid string, block int, hash string) {
	// store in database
	stmtstr := `insert into publisher (name, address, timestamp, txid, block, emailmd5, bitmessage, hash, signature, active) values (?, ?, ?, "` + txid + `", ` + strconv.Itoa(block) + `, ?, ?, "` + hash + `", ?, 1)`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 100")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(publisher.AlexandriaPublisher.Name, publisher.AlexandriaPublisher.Address, publisher.AlexandriaPublisher.Timestamp, publisher.AlexandriaPublisher.Emailmd5, publisher.AlexandriaPublisher.Bitmessage, publisher.Signature)
	if err != nil {
		fmt.Println("exit 101")
		log.Fatal(stmterr)
	}

	stmt.Close()

}

func VerifyPublisher(b []byte) (AlexandriaPublisher, error) {

	var v AlexandriaPublisher
	var i interface{}
	var m map[string]interface{}

	if !strings.HasPrefix(string(b), `{ "alexandria-publisher"`) &&
		!strings.HasPrefix(string(b), `{"alexandria-publisher"`) {
		return v, ErrWrongPrefix
	}

	// fmt.Printf("Attempting to verify alexandria-publisher JSON...")

	if !utility.IsJSON(string(b)) {
		return v, ErrNotJSON
	}

	err := json.Unmarshal(b, &v)
	if err != nil {
		return v, err
	}

	errr := json.Unmarshal(b, &i)
	if errr != nil {
		return v, err
	}

	m = i.(map[string]interface{})
	var signature string

	// check the JSON object root key
	// find the signature string
	for key, val := range m {
		if key == "signature" {
			signature = val.(string)
		} else {
			if key != PUBLISHER_ROOT_KEY {
				return v, errors.New("can't verify publisher - JSON object root key doesn't match accepted value")
			}
		}
	}

	// verify signature
	if v.Signature != signature {
		return v, ErrBadSignature
	}

	// verify signature was created by this address
	// signature pre-image for publisher is <name>-<address>-<timestamp>
	val, _ := utility.CheckSignature(v.AlexandriaPublisher.Address, signature, v.AlexandriaPublisher.Name+"-"+v.AlexandriaPublisher.Address+"-"+strconv.FormatInt(v.AlexandriaPublisher.Timestamp, 10))
	if val == false {
		return v, ErrBadSignature
	}

	// fmt.Println(" -- VERIFIED --")
	return v, nil

}
