package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dloa/media-protocol/utility"
	"log"
	"strconv"
	"strings"
)

type OIPTransfer struct {
	Reference string `json:"tx"`
	To        string `json:"to"`
	From      string `json:"fro"`
	Timestamp int64  `json:"ts"`
	Signature string `json:"sig"`
}

type OIPTransferWrapper struct {
	OIPTransfer OIPTransfer `json:"oip-transfer"`
}

func StoreOIPTransfer(oip_t OIPTransfer, dbtx *sql.Tx) {
	// ToDo: perhaps save a record of this happening

	stmtstr := `UPDATE media SET publisher=? WHERE txid=?`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 600")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(oip_t.To, oip_t.Reference)
	if err != nil {
		fmt.Println("exit 601")
		log.Fatal(stmterr)
	}

	stmt.Close()
}

func VerifyOIPTransfer(s string, block int) (OIPTransfer, error) {
	var oip_tw OIPTransferWrapper
	var oip_t OIPTransfer

	if block < 1988636 {
		return oip_t, ErrTooEarly
	}

	if !strings.HasPrefix(s, `{"oip-transfer"`) {
		return oip_t, ErrWrongPrefix
	}

	if !utility.IsJSON(s) {
		return oip_t, ErrNotJSON
	}

	err := json.Unmarshal([]byte(s), &oip_tw)
	if err != nil {
		return oip_t, err
	}
	oip_t = oip_tw.OIPTransfer

	if len(oip_t.Reference) != 64 {
		return oip_t, ErrInvalidReference
	}

	if !utility.CheckAddress(oip_t.To) {
		return oip_t, ErrInvalidAddress
	}

	preImage := oip_t.Reference + "-" + oip_t.To + "-" + strconv.FormatInt(oip_t.Timestamp, 10)
	valid, err := utility.CheckSignature(oip_t.From, oip_t.Signature, preImage)
	if !valid {
		return oip_t, err
	}

	return oip_t, nil
}
