package messages

import (
	"database/sql"
	"fmt"
	"github.com/dloa/media-protocol/utility"
	"log"
	"strconv"
)

func StoreOIP041Deactivate(o Oip041, dbtx *sql.Tx) error {
	// ToDo: perhaps save a record of this happening
	// ToDo: refactor the table deciding logic, it's terribad

	oip_d := o.Deactivate

	// Check media table
	table := "media"
	stmt, err := dbtx.Prepare(`SELECT publisher FROM ` + table + ` WHERE txid = ? LIMIT 1;`)
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(oip_d.Reference)
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
		row := stmt.QueryRow(oip_d.Reference)
		err = row.Scan(&publisher)
		if err != nil {
			if err == sql.ErrNoRows {
				return err
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

	preImage := oip_d.Reference + "-" + strconv.FormatInt(oip_d.Timestamp, 10)
	valid, err := utility.CheckSignature(publisher, o.Signature, preImage)
	if !valid {
		fmt.Println("Signature check failed.")
		return err
	}

	stmtstr := `update ` + table + ` set invalidated = 1 where publisher = ? and txid = ?`

	stmt, err = dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 606")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(publisher, oip_d.Reference)
	if stmterr != nil {
		fmt.Println("exit 607")
		log.Fatal(stmterr)
	}

	stmt.Close()
	return nil
}

func VerifyOIP041Deactivate(o Oip041) (Oip041, error) {
	oip_d := o.Deactivate

	if len(oip_d.Reference) != 64 {
		return o, ErrInvalidReference
	}
	return o, nil
}
