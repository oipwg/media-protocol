package oip042

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/metacoin/flojson"
)

// Oip 042 ToDo list
// - Split tables such that all common info goes into artifact042 table
// -- type specific details get split off to separate tables
// -- allows for better searching/getById, etc
// -- creates a unified listing instead of segmented
// - Edits
// - Deactivates
// - Transfers
// - More artifact detail types

var ErrMissingOipAction = errors.New("missing oip action")
var ErrBadSignature = errors.New("bad signature")

func (o Oip042) ValidateIncoming(tx *flojson.TxRawResult, txComment string, txid string, block *flojson.BlockResult, dbtx *sqlx.Tx) (OipAction, error) {

	// only process the first match, disregard remaining
	// otherwise there's order of operations to consider

	ctx := OipContext{
		Tx:          tx,
		TxComment:   txComment,
		TxId:        txid,
		BlockHeight: block.Height,
		Block:       block,
		DbTx:        dbtx,
	}

	if o.Publish != nil {
		return o.Publish.Validate(ctx)
	}

	if o.Register != nil {
		return o.Register.Validate(ctx)
	}

	if o.Transfer != nil {
		return o.Transfer.Validate(ctx)
	}

	if o.Edit != nil {
		return o.Edit.Validate(ctx)
	}

	if o.Deactivate != nil {
		return o.Deactivate.Validate(ctx)
	}

	return nil, ErrMissingOipAction
}

func SetupTables(dbtx *sqlx.Tx) error {
	if _, err := dbtx.Exec(createArtifactTable); err != nil {
		return err
	}
	if _, err := dbtx.Exec(createTomogramTable); err != nil {
		return err
	}
	if _, err := dbtx.Exec(createPropertyPartyTable); err != nil {
		return err
	}
	if _, err := dbtx.Exec(createPropertySpatialUnitTable); err != nil {
		return err
	}
	if _, err := dbtx.Exec(createPropertyTenureTable); err != nil {
		return err
	}

	return nil
}

func GetAllArtifacts(dbtx *sqlx.Tx) ([]interface{}, error) {
	q := squirrel.Select("json", "txid", "publisher").
		From("artifact").
		Where(squirrel.Eq{"active": 1}).
		Where(squirrel.Eq{"invalidated": 0})
	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dbtx.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type OipInner struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	type rWrap struct {
		OipInner  `json:"oip042"`
		Txid      string `json:"txid"`
		Publisher string `json:"publisher"`
	}
	var res []interface{}
	for rows.Next() {
		var j json.RawMessage
		var txid string
		var publisher string
		err := rows.Scan(&j, &txid, &publisher)
		if err != nil {
			return nil, err
		}
		res = append(res, rWrap{OipInner{j}, txid, publisher})
	}

	return res, nil
}

func GetById(dbh *sqlx.DB, artId string) (interface{}, error) {
	q := squirrel.Select("a.json", "a.txid", "a.publisher", "p.name").
		From("artifact AS a").
		LeftJoin("publisher AS p ON p.address = a.publisher").
		Where(squirrel.Eq{"a.active": 1}).
		Where(squirrel.Eq{"a.invalidated": 0})
	if len(artId) == 64 {
		q = q.Where(squirrel.Eq{"a.txid": artId})
	} else {
		q = q.Where("a.txid LIKE ?", fmt.Sprint(artId, "%"))
	}
	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	row := dbh.QueryRow(query, args...)

	var j json.RawMessage
	var txid string
	var publisher string
	var publisherNameN sql.NullString
	err = row.Scan(&j, &txid, &publisher, &publisherNameN)
	if err != nil {
		return nil, err
	}

	publisherName := ""
	if publisherNameN.Valid {
		publisherName = publisherNameN.String
	}

	type OipInner struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	type rWrap struct {
		OipInner      `json:"oip042"`
		Txid          string `json:"txid"`
		Publisher     string `json:"publisher"`
		PublisherName string `json:"publisherName"`
	}

	return rWrap{OipInner{j}, txid, publisher, publisherName}, nil
}

func GetByType(dbtx *sqlx.Tx, t string, st string, page uint64, results uint64, pub string) (interface{}, error) {
	q := squirrel.Select("a.json", "a.txid", "a.publisher").
		From("artifact as a").
		Where(squirrel.Eq{"a.active": 1}).
		Where(squirrel.Eq{"a.invalidated": 0})
	qc := squirrel.Select("count(*)").
		From("artifact as a").
		Where(squirrel.Eq{"a.active": 1}).
		Where(squirrel.Eq{"a.invalidated": 0})

	if t != "*" && t != "" {
		if t == "-" {
			t = ""
		}
		q = q.Where(squirrel.Eq{"a.type": t})
		qc = qc.Where(squirrel.Eq{"a.type": t})
	}
	if st != "*" && st != "" {
		if st == "-" {
			st = ""
		}
		q = q.Where(squirrel.Eq{"a.subType": st})
		qc = qc.Where(squirrel.Eq{"a.subType": st})
	}
	if pub != "*" && pub != "" {
		q = q.Where(squirrel.Eq{"a.publisher": pub})
		qc = qc.Where(squirrel.Eq{"a.publisher": pub})
	}
	if results != 0 {
		q = q.Limit(results)
		// page only makes sense if there are result limits
		if page != 0 {
			q = q.Offset((page - 1) * results)
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dbtx.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	query, args, err = qc.ToSql()
	if err != nil {
		return nil, err
	}
	var total uint64
	err = dbtx.QueryRow(query, args...).Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			total = 0
		} else {
			return nil, err
		}
	}

	type OipInner struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	type rWrap struct {
		OipInner  `json:"oip042"`
		Txid      string `json:"txid"`
		Publisher string `json:"publisher"`
	}
	var res []interface{}
	for rows.Next() {
		var j json.RawMessage
		var txid string
		var publisher string
		err := rows.Scan(&j, &txid, &publisher)
		if err != nil {
			return nil, err
		}
		res = append(res, rWrap{OipInner{j}, txid, publisher})
	}

	type ret struct {
		Total   uint64        `json:"total"`
		Count   uint64        `json:"count"`
		Pages   uint64        `json:"pages"`
		Results []interface{} `json:"results"`
	}

	count := uint64(len(res))
	pages := uint64(1)
	if count != 0 {
		pages = total / count
		if total%count != 0 {
			pages++ // trailing partial page
		}
	}
	return ret{Total: total, Count: count, Pages: pages, Results: res}, nil
}

func GetByPublisher(dbtx *sqlx.Tx, publisher string) ([]interface{}, error) {
	q := squirrel.Select("a.json", "a.txid", "a.publisher").
		From("artifact as a").
		Where(squirrel.Eq{"active": 1}).
		Where(squirrel.Eq{"invalidated": 0})

	if publisher != "*" && publisher != "" {
		if publisher == "-" {
			publisher = ""
		}
		q = q.Where(squirrel.Eq{"a.publisher": publisher})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dbtx.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type OipInner struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	type rWrap struct {
		OipInner  `json:"oip042"`
		Txid      string `json:"txid"`
		Publisher string `json:"publisher"`
	}
	var res []interface{}
	for rows.Next() {
		var j json.RawMessage
		var txid string
		var publisher string
		err := rows.Scan(&j, &txid, &publisher)
		if err != nil {
			return nil, err
		}
		res = append(res, rWrap{OipInner{j}, txid, publisher})
	}

	return res, nil
}

const createArtifactTable = `
-- General artifact information
CREATE TABLE IF NOT EXISTS artifact
(
  uid         INTEGER PRIMARY KEY AUTOINCREMENT,
  block       INT  NOT NULL,
  txid        TEXT NOT NULL,
  json        TEXT NOT NULL,
  unixtime    INT  NOT NULL,
  title       TEXT NOT NULL,
  type        TEXT NOT NULL,
  subType     TEXT NOT NULL,
  publisher   TEXT NOT NULL,
  tags        TEXT NOT NULL,
  artCost	  FLOAT NOT NULL DEFAULT 0,
  pubFeeUsd   FLOAT NOT NULL DEFAULT 0,
  artSize     INT NOT NULL DEFAULT 0,
  invalidated BOOLEAN             DEFAULT 0 NOT NULL,
  active      BOOLEAN             DEFAULT 0 NOT NULL,
  nsfw        BOOLEAN             DEFAULT 0 NOT NULL,
  hasDetails  BOOLEAN             DEFAULT 0 NOT NULL
);
CREATE INDEX IF NOT EXISTS artifact_txid_uindex  ON artifact (txid);
CREATE INDEX IF NOT EXISTS artifact_subtype_index  ON artifact (subtype);
CREATE INDEX IF NOT EXISTS artifact_type_index  ON artifact (type);
CREATE INDEX IF NOT EXISTS artifact_publisher_index  ON artifact (publisher);
CREATE INDEX IF NOT EXISTS artifact_block_index ON artifact (block);
`
