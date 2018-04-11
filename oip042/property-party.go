package oip042

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PartyDetails struct {
	Ns        string          `json:"ns,omitempty"`
	PartyRole string          `json:"partyRole,omitempty"`
	PartyType string          `json:"partyType,omitempty"`
	Tenures   []string        `json:"tenures,omitempty"`
	Groups    []string        `json:"groups,omitempty"`
	Members   []string        `json:"members,omitempty"`
	Attrs     json.RawMessage `json:"attrs,omitempty"`
}

type PublishPropertyParty struct {
	PublishArtifact
	PartyDetails
}

func (ppp PublishPropertyParty) Validate(context OipContext) (OipAction, error) {
	err := json.Unmarshal(ppp.Details, &ppp.PartyDetails)
	if err != nil {
		return nil, err
	}

	return ppp, nil
}

func (ppp PublishPropertyParty) Store(context OipContext, dbtx *sqlx.Tx) error {
	j, err := json.Marshal(ppp)
	if err != nil {
		return err
	}

	q := sq.Insert("artifactPropertyParty").
		Columns("ns", "partyRole", "partyType",
			"active", "block", "json", "tags", "timestamp",
			"title", "txid", "type", "subType", "publisher").
		Values(ppp.Ns, ppp.PartyRole, ppp.PartyType,
			1, context.BlockHeight, j, ppp.Info.Tags, ppp.Timestamp,
			ppp.Info.Title, context.TxId, ppp.Type, ppp.SubType, ppp.FloAddress)

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

func (ppp PublishPropertyParty) MarshalJSON() ([]byte, error) {
	pa := ppp.PublishArtifact
	buf, err := json.Marshal(ppp.PartyDetails)
	if err != nil {
		return nil, err
	}
	pa.Details = buf
	return json.Marshal(pa)
}

func GetAllPropertyParty(dbtx *sqlx.Tx) ([]interface{}, error) {
	// ToDo combine/simplify these GetAll functions similar to GetById
	q := sq.Select("json", "txid", "publisher").
		From("artifactPropertyParty").
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

const createPropertyPartyTable = `CREATE TABLE IF NOT EXISTS artifactPropertyParty
(
  uid            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,

  -- Property-Party Fields
  ns             TEXT NOT NULL,
  partyRole      TEXT NOT NULL,
  partyType      TEXT NOT NULL,

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
