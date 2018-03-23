package oip042

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/metacoin/flojson"
)

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
		fmt.Println("publish 42")
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
	_, err := dbtx.Exec(createTomogramTable)
	return err
}

func GetById(dbh *sqlx.DB, artId string) (interface{}, error) {
	sql, args, err := squirrel.Select("json").
		From("artifactsResearchTomogram").
		Where(squirrel.Eq{"txid": artId}).Where(squirrel.Eq{"active": 1}).ToSql()

	if err != nil {
		return nil, err
	}
	row := dbh.QueryRow(sql, args...)

	var j json.RawMessage
	err = row.Scan(&j)
	if err != nil {
		return nil, err
	}

	type OipInner struct {
		Artifact json.RawMessage `json:"artifact"`
	}
	type rWrap struct {
		OipInner `json:"oip042"`
	}

	return rWrap{OipInner{j}}, nil
}
