package oip042

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
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

func (ppp PublishPropertyParty) Store(context OipContext) error {
	j, err := json.Marshal(ppp)
	if err != nil {
		return err
	}

	q := sq.Insert("artifact").
		Columns("active", "block", "json", "tags", "unixtime",
			"title", "txid", "type", "subType", "publisher", "hasDetails").
		Values(1, context.BlockHeight, j, ppp.Info.Tags, ppp.Timestamp,
			ppp.Info.Title, context.TxId, ppp.Type, ppp.SubType, ppp.FloAddress, 1)

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	res, err := context.DbTx.Exec(query, args...)
	if err != nil {
		return err
	}

	artifactId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	q = sq.Insert("detailsPropertyParty").
		Columns("artifactId", "ns", "partyRole", "partyType").
		Values(artifactId, ppp.Ns, ppp.PartyRole, ppp.PartyType)

	query, args, err = q.ToSql()
	if err != nil {
		return err
	}

	res, err = context.DbTx.Exec(query, args...)
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

const createPropertyPartyTable = `
-- Property-Party details
CREATE TABLE IF NOT EXISTS detailsPropertyParty
(
  uid         INTEGER PRIMARY KEY AUTOINCREMENT,
  artifactId  INT  NOT NULL,
  ns          TEXT NOT NULL,
  partyRole   TEXT NOT NULL,
  partyType   TEXT NOT NULL,
  CONSTRAINT detailsPropertyParty_artifactId_uid_fk FOREIGN KEY (artifactId) REFERENCES artifact (uid) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS detailsPropertyParty_artifactId_uindex  ON detailsPropertyParty (artifactId);
CREATE INDEX IF NOT EXISTS detailsPropertyParty_partyRole_index  ON detailsPropertyParty (partyRole);
CREATE INDEX IF NOT EXISTS detailsPropertyParty_partyType_index  ON detailsPropertyParty (partyType);
`
