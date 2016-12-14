package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleOIP041Edit(o Oip041, txid string, block int, dbtx *sql.Tx) error {
	return ErrNotImplemented

	// ToDo: Check the signature... but first decide what to sign

	stmtstr := `SELECT ('json','txid','publisher')
		FROM oip_artifact WHERE txid=? LIMIT 1`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("RIP Handle Edit 1")
		log.Fatal(err)
	}
	row := stmt.QueryRow(txid)

	var old Oip041
	var sJSON string
	var publisher string
	var txID string

	err = row.Scan(&sJSON, &txID, &publisher)
	if err != nil {
		fmt.Println("RIP Handle Edit 2")
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(sJSON), &old)
	if err != nil {
		fmt.Println("RIP Handle Edit 3")
		return err
	}

	if len(o.Edit.Add) > 0 {
		// we got stuff to add
	}
	if len(o.Edit.Edit) > 0 {
		// we got stuff to edit
		for k, v := range o.Edit.Edit {
			p := strings.Split(k, ".")
			if len(p) == 1 {
				updateField(k, v, txid, dbtx)
				// something along these lines...
			}
		}
	}
	if len(o.Edit.Remove) > 0 {
		// we got stuff to remove
	}

	b, err := json.Marshal(old)
	if err != nil {
		fmt.Println("RIP Handle Edit 4")
		return err
	}

	fmt.Println(string(b))

	return nil
}

/*
var oip041_edit_example_obj Oip041 = Oip041{
	Edit: Oip041Edit{
		Add: map[string]string{
			"payment.tokens": "FREEBIEOFTHEWEEK:\"1\"",
		},
		Edit: map[string]string{
			"files[0].dname": "Throwing Stones",
			"files[0].fname": "1 - Throwing Stones.mp3",
		},
		Remove: []string{
			"tokens.LTBCOIN",
		},
		Timestamp: 1234,
		TxID:      "96bad8e17f908da4c695c58b0f843a03928e338b361b3035ed16a864eafc31a2",
	},
	Signature: "<SignatureOfSomething>",
}
*/

func updateField(key string, value string, txid string, dbtx *sql.Tx) error {
	return ErrNotImplemented

	stmtstr := `UPDATE oip_artifact SET ?=? WHERE txid=?`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("RIP Update Field 1")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(key, value, txid)
	if err != nil {
		fmt.Println("RIP Update Field 1")
		log.Fatal(stmterr)
	}

	stmt.Close()

	return nil
}
