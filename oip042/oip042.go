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

func GetById(dbh *sqlx.DB, artId string) (interface{}, error) {
	// ToDo this function would appreciate the unified table structure
	var err error
	var res interface{}
	res, err = GetByIdFromTable(dbh, artId, "artifactsResearchTomogram")
	if err == nil || err != sql.ErrNoRows {
		return res, err
	}
	res, err = GetByIdFromTable(dbh, artId, "artifactPropertyParty")
	if err == nil || err != sql.ErrNoRows {
		return res, err
	}
	res, err = GetByIdFromTable(dbh, artId, "artifactPropertySpatialUnit")
	if err == nil || err != sql.ErrNoRows {
		return res, err
	}
	res, err = GetByIdFromTable(dbh, artId, "artifactPropertyTenure")
	if err == nil || err != sql.ErrNoRows {
		return res, err
	}
	return nil, sql.ErrNoRows
}

func GetByIdFromTable(dbh *sqlx.DB, artId string, table string) (interface{}, error) {
	q := squirrel.Select("json", "txid", "publisher").
		From(table).Where(squirrel.Eq{"active": 1})
	if len(artId) == 64 {
		q = q.Where(squirrel.Eq{"txid": artId})
	} else {
		q = q.Where("txid LIKE ?", fmt.Sprint(artId, "%"))
	}
	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	row := dbh.QueryRow(query, args...)

	var j json.RawMessage
	var txid string
	var publisher string
	err = row.Scan(&j, &txid, &publisher)
	if err != nil {
		return nil, err
	}

	type OipInner struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	type rWrap struct {
		OipInner  `json:"oip042"`
		Txid      string `json:"txid"`
		Publisher string `json:"publisher"`
	}

	return rWrap{OipInner{j}, txid, publisher}, nil
}
