package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/oipwg/media-protocol/utility"
	"strings"
)

func VerifyOIP041(s string, block int) (Oip041, error) {
	if block < 1997454 {
		return Oip041{}, ErrTooEarly
	}

	if !utility.IsJSON(s) {
		return Oip041{}, ErrNotJSON
	}

	// make sure signature isn't null in the decoded OIP string
	dec, err := DecodeOIP041(s)
	if err != nil {
		return dec, err
	}
	if dec.Signature == "" {
		return dec, ErrBadSignature
	}
	// ToDo: Validate signature

	if dec.Artifact.Timestamp != 0 {
		err := dec.Artifact.CheckRequiredFields()
		if err != nil {
			return dec, err
		}
	}

	dec.artSize = len(s)

	return dec, nil
}

func DecodeOIP041(s string) (Oip041, error) {
	oip041w := Oip041Wrapper{}
	err := json.Unmarshal([]byte(s), &oip041w)
	return oip041w.Oip041, err
}

func APIGetAllOIP041(dbtx *sql.Tx) ([]Oip041ArtifactAPIResult, error) {
	stmtStr := `select a.block, a.json, a.tags, a.timestamp,
				a.title, a.txid, a.type, a.year, a.publisher, p.name, a.artCost, a.artSize, a.pubFeeUSD
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
		var s string

		rows.Scan(&a.Block, &s, &a.Tags, &a.Timestamp,
			&a.Title, &a.TxID, &a.Type, &a.Year, &a.Publisher, &a.PublisherName, &a.ArtCost, &a.ArtSize, &a.PubFeeUSD)

		json.Unmarshal([]byte(s), &a.OIP041)
		results = append(results, a)
	}

	stmt.Close()
	rows.Close()

	return results, nil
}

func CreateTables(dbTx *sql.Tx) error {
	for _, v := range oip041SqliteCreateStatements {
		fmt.Printf("\nRunning table query:  %s\n", v.name)
		stmt, err := dbTx.Prepare(v.sql)
		if err != nil {
			if !strings.HasPrefix(v.name, "!addcol!") {
				// ToDo: HACK! There is no "add column if not exists"
				// instead the duplicate column error is ignored
				// update to utilize schema versioning and run queries
				// as required instead of all queries every time
				// https://stackoverflow.com/questions/3604310/alter-table-add-column-if-not-exists-in-sqlite
				return err
			} else {
				fmt.Println("Column already added.")
				continue
			}
		}
		_, stmt_err := stmt.Exec()
		if stmt_err != nil {
			return stmt_err
		}
		stmt.Close()
	}
	return nil
}

var oip041SqliteCreateStatements = []struct {
	name string
	sql  string
}{
	{"create oip_artifact table",
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
		'publisher'	TEXT NOT NULL,
		'artCost' FLOAT NOT NULL,
		'artSize' INTEGER NOT NULL,
		'pubFeeUSD' FLOAT NOT NULL
	);`},
	{
		"!addcol! ArtCost",
		"ALTER TABLE 'oip_artifact' ADD COLUMN 'artCost' FLOAT;",
	},
	{
		"!addcol! o.pubFeeUSD",
		"ALTER TABLE 'oip_artifact' ADD COLUMN 'pubFeeUSD';" +
			"ALTER TABLE 'media' ADD COLUMN FLOAT 'artSize' INTEGER;",
	},
	{
		"!addcol! m.pubFeeUSD",
		"ALTER TABLE 'media' ADD COLUMN 'pubFeeUSD' FLOAT;" +
			"ALTER TABLE 'media' ADD COLUMN 'artSize' INTEGER;",
	},
}
