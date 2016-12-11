package messages

import (
	"database/sql"
	"fmt"
	"github.com/dloa/media-protocol/utility"
	"log"
	"strconv"
)

func StoreOIPTransfer(oip_t Oip041Transfer, dbtx *sql.Tx) {
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

func HandleOIP041Transfer(o Oip041, txid string, processingBlock int, dbtx *sql.Tx) {
	// ToDo: Verify rights to transfer
	// leedle leedle leedle lee

	oip_t := o.Transfer

	if len(oip_t.Reference) != 64 {
		return //oip_t, ErrInvalidReference
	}

	if !utility.CheckAddress(oip_t.To) {
		return //oip_t, ErrInvalidAddress
	}

	preImage := oip_t.Reference + "-" + oip_t.To + "-" + oip_t.From + "-" + strconv.FormatInt(oip_t.Timestamp, 10)
	valid, err := utility.CheckSignature(oip_t.From, o.Signature, preImage)
	if !valid {
		return //oip_t, err
	}

	stmt, err := dbtx.Prepare(`SELECT publisher FROM media WHERE txid = ? LIMIT 1;`)
	if err != nil {
		log.Fatal(err)
	}

	row := stmt.QueryRow(oip_t.Reference)
	var publisher string
	err = row.Scan(&publisher)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()

	fmt.Println(publisher)
	if publisher == oip_t.From {
		fmt.Println("Transfered")
		StoreOIPTransfer(oip_t, dbtx)
	} else {
		fmt.Println("Denied")
	}
	fmt.Println()
}
