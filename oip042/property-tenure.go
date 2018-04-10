package oip042

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type TenureDetails struct {
	Ns           string          `json:"ns"`
	TenureType   string          `json:"tenureType"`
	Tenures      []string        `json:"tenures"`
	SpatialUnits []string        `json:"spatialUnits"`
	Attrs        json.RawMessage `json:"attrs"`
}

type PublishPropertyTenure struct {
	PublishArtifact
	TenureDetails
}

func (ppt PublishPropertyTenure) Validate(context OipContext) (OipAction, error) {
	err := json.Unmarshal(ppt.Details, &ppt.TenureDetails)
	if err != nil {
		return nil, err
	}

	return ppt, nil
}

func (ppt PublishPropertyTenure) Store(context OipContext, dbtx *sqlx.Tx) error {
	j, err := json.Marshal(ppt)
	if err != nil {
		return err
	}

	q := sq.Insert("artifactPropertyTenure").
		Columns("ns", "tenureType",
			"active", "block", "json", "tags", "timestamp",
			"title", "txid", "type", "subType", "publisher").
		Values(ppt.Ns, ppt.TenureType,
			1, context.BlockHeight, j, ppt.Info.Tags, ppt.Timestamp,
			ppt.Info.Title, context.TxId, ppt.Type, ppt.SubType, ppt.FloAddress)

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = dbtx.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (ppt PublishPropertyTenure) MarshalJSON() ([]byte, error) {
	pa := ppt.PublishArtifact
	buf, err := json.Marshal(ppt.TenureDetails)
	if err != nil {
		return nil, err
	}
	pa.Details = buf
	return json.Marshal(pa)
}

func GetAllPropertyTenure(dbtx *sqlx.Tx) ([]interface{}, error) {
	// ToDo combine/simplify these GetAll functions similar to GetById
	q := sq.Select("json", "txid", "publisher").
		From("artifactPropertyTenure").
		Where("active = ?", 1)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dbtx.Queryx(sql, args...)
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

const createPropertyTenureTable = `CREATE TABLE IF NOT EXISTS artifactPropertyTenure
(
  uid            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,

  -- Property-Tenure Fields
  ns             TEXT NOT NULL,
  tenureType     TEXT NOT NULL,

  -- General OIP Fields
  active         INTEGER NOT NULL,
  block          INTEGER NOT NULL,
  invalidated    INTEGER                      DEFAULT 0,
  json           INTEGER NOT NULL,
  tags           TEXT    NOT NULL,
  timestamp      INTEGER NOT NULL,
  title          TEXT    NOT NULL,
  txid           TEXT    NOT NULL,
  type           TEXT    NOT NULL,
  subType        TEXT    NOT NULL,
  validated      INTEGER                      DEFAULT 0,
  publisher      TEXT    NOT NULL,
  nsfw           BOOLEAN                      DEFAULT 0
)`
