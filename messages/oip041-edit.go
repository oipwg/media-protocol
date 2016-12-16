package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/bitspill/json-patch"
	"log"
	"strings"
)

func HandleOIP041Edit(o Oip041, txid string, block int, dbtx *sql.Tx) error {

	// ToDo: This is super ugly
	patch := UnSquashPatch(strings.Replace(string(o.Edit.Patch), `"path":"/`, `"path":"/artifact/`, -1))
	fmt.Printf("Patch:\n%s\n", patch)
	obj, err := jsonpatch.DecodePatch([]byte(patch))
	if err != nil {
		log.Fatalf("Failed to decode patch:\n%v", err)
	}

	// ToDo: Check the signature... but first decide what to sign

	stmtstr := `SELECT json,txid,publisher
		FROM oip_artifact WHERE txid=? LIMIT 1`

	stmt, err := dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("RIP Handle Edit 1")
		log.Fatal(err)
	}
	row := stmt.QueryRow(o.Edit.TxID)

	var sJSON string
	var publisher string
	var txID string

	err = row.Scan(&sJSON, &txID, &publisher)
	if err != nil {
		fmt.Println("RIP Handle Edit 2")
		log.Fatal(err)
	}
	stmt.Close()

	fmt.Printf("Pre-Patch:\n%s\n", sJSON)
	out, err := obj.Apply([]byte(sJSON))
	if err != nil {
		log.Fatalf("Failed to apply patch:\n%v", err)
	}
	fmt.Printf("Post-Patch Result:\n%s\n", string(out))

	stmtstr = `UPDATE oip_artifact SET json=? WHERE txid=?`
	// ToDo: apply update to searchable meta-data
	stmt, err = dbtx.Prepare(stmtstr)
	if err != nil {
		fmt.Println("exit 600")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(string(out), o.Edit.TxID)
	if stmterr != nil {
		fmt.Println("exit 601")
		log.Fatal(stmterr)
	}

	stmt.Close()

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

func UnSquashPatch(sp string) string {
	var p map[string][]map[string]*json.RawMessage // yikes
	var up jsonpatch.Patch

	err := json.Unmarshal([]byte(sp), &p)
	fmt.Println(err)

	// op="Add", arr="Array of actions"
	for op, arr := range p {
		// _="index", act="Action object"
		for _, act := range arr {
			o := json.RawMessage([]byte(`"` + op + `"`))
			act["op"] = &o
			up = append(up, act)
		}
	}

	fmt.Println(up)
	usp, err := json.Marshal(&up)
	fmt.Println(err)
	return string(usp)
}
