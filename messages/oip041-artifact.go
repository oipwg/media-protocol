package messages

import (
	"database/sql"
	"encoding/json"
	"strings"
)

func (o Oip041) GetJSON() (string, error) {
	// ToDo: remove redundant Storage items, potentially cache?
	var s string

	b, err := json.Marshal(o)
	s = string(b)

	return s, err
}

func StoreOIP041Artifact(o Oip041, txid string, block int, dbtx *sql.Tx) error {
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
