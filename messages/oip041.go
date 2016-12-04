package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func DecodeOIP041(s string) (Oip041, error) {
	oip041w := Oip041Wrapper{}
	err := json.Unmarshal([]byte(s), &oip041w)
	return oip041w.Oip041, err
}

func CreateTables(dbTx *sql.Tx) error {
	// ToDo: update the triggers for deactivating media
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
		`CREATE TABLE if not exists 'oip_media' (
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
