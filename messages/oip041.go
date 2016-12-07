package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

func VerifyOIP041(s string, block int) (Oip041, error) {
	if block < 1997454 {
		return Oip041{}, ErrTooEarly
	}
	return DecodeOIP041(s)
}

func DecodeOIP041(s string) (Oip041, error) {
	oip041w := Oip041Wrapper{}
	err := json.Unmarshal([]byte(s), &oip041w)
	return oip041w.Oip041, err
}

func StoreOIP041(o Oip041, txid string, block int, dbtx *sql.Tx) error {
	// store in database
	stmtStr := `INSERT INTO 'oip_artifact'
		('active','block','json','tags','timestamp',
		'title','txid','type','year','publisher')
		VALUES (?,?,?,?,?,?,?,?,?,?);`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	s, err := o.GetJSON()
	if err != nil {
		return nil
	}

	_, err = stmt.Exec(1, block, s, strings.Join(o.Artifact.Info.ExtraInfo.Tags, ","),
		o.Artifact.Timestamp, o.Artifact.Info.Title, txid, o.Artifact.Type,
		o.Artifact.Info.Year, o.Artifact.Publisher)
	if err != nil {
		return err
	}

	return nil
}

func APIGetAllOIP041(dbtx *sql.Tx) ([]Oip041ArtifactAPIResult, error) {
	stmtStr := `select a.block, a.json, a.tags, a.timestamp,
				a.title, a.txid, a.type, a.year, a.publisher
				from oip_artifact as a join publisher as p
				where p.address = a.publisher and a.invalidated = 0`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		return []Oip041ArtifactAPIResult{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []Oip041ArtifactAPIResult{}, err
	}

	var results []Oip041ArtifactAPIResult

	for rows.Next() {
		var a Oip041ArtifactAPIResult
		rows.Scan(&a.Block, &a.JSON, &a.Tags, &a.Timestamp,
			&a.Title, &a.TxID, &a.Type, &a.Year, &a.Publisher)

		results = append(results, a)
	}

	stmt.Close()
	rows.Close()

	return results, nil
}

func CreateTables(dbTx *sql.Tx) error {
	for _, v := range oip041SqliteCreateStatements {
		fmt.Printf("Creating %s\n", v.name)
		stmt, err := dbTx.Prepare(v.sql)
		if err != nil {
			return err
		}
		_, stmt_err := stmt.Exec()
		if stmt_err != nil {
			return stmt_err
		}
		dbTx.Commit()
		stmt.Close()
	}
	return nil
}

var oip041SqliteCreateStatements = []struct {
	name string
	sql  string
}{
	{"oip_media table",
		`CREATE TABLE if not exists 'oip_artifact' (
		'uid'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		'active'	INTEGER NOT NULL,
		'block'	INTEGER NOT NULL,
		'invalidated' INTEGER default 0,
		'isAlbum'	INTEGER,
		'isFree'	INTEGER,
		'json'	INTEGER NOT NULL,
		'tags'	TEXT NOT NULL,
		'timestamp'	INTEGER NOT NULL,
		'title'	TEXT NOT NULL,
		'txid'	TEXT NOT NULL,
		'type'	TEXT NOT NULL,
		'validated' INTEGER default 0,
		'year'	INTEGER NOT NULL,
		'publisher'	TEXT NOT NULL
	);`},
}
