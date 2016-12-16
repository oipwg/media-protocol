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
	// ToDo: refactor the table deciding logic, it's terribad

	// Check media table
	table := "media"
	stmt, err := dbtx.Prepare(`SELECT publisher FROM ` + table + ` WHERE txid = ? LIMIT 1;`)
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(oip_t.Reference)
	var publisher string
	err = row.Scan(&publisher)
	if err != nil {
		// Close media check
		stmt.Close()

		// Check oip table
		fmt.Println("Might be an OIP")
		table = "oip_artifact"
		stmt, err := dbtx.Prepare(`SELECT publisher FROM ` + table + ` WHERE txid = ? LIMIT 1;`)
		if err != nil {
			log.Fatal(err)
		}
		row := stmt.QueryRow(oip_t.Reference)
		err = row.Scan(&publisher)
		if err != nil {
			if err == sql.ErrNoRows {
				return
			} else {
				log.Fatal(err)
			}
		}
		stmt.Close()
	} else {
		// Close media check
		stmt.Close()
	}

	fmt.Println(publisher)
	if publisher == oip_t.From {
		fmt.Println("Transfer Ok")
	} else {
		fmt.Println("Transfer Denied")
		return
	}

	stmtstr := `UPDATE ` + table + ` SET publisher=? WHERE txid=?`

	stmt, err = dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 600")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(oip_t.To, oip_t.Reference)
	if stmterr != nil {
		fmt.Println("exit 601")
		log.Fatal(stmterr)
	}

	stmt.Close()
}

func VerifyOIP041Transfer(o Oip041) (Oip041, error) {
	oip_t := o.Transfer

	if len(oip_t.Reference) != 64 {
		return o, ErrInvalidReference
	}

	if !utility.CheckAddress(oip_t.To) {
		return o, ErrInvalidAddress
	}

	preImage := oip_t.Reference + "-" + oip_t.To + "-" + oip_t.From + "-" + strconv.FormatInt(oip_t.Timestamp, 10)
	valid, err := utility.CheckSignature(oip_t.From, o.Signature, preImage)
	if !valid {
		return o, err
	}

	return o, nil
}
