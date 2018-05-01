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

	cv := map[string]interface{}{
		"active":     1,
		"block":      context.BlockHeight,
		"json":       j,
		"tags":       ppp.Info.Tags,
		"unixtime":   ppp.Timestamp,
		"title":      ppp.Info.Title,
		"txid":       context.TxId,
		"type":       ppp.Type,
		"subType":    ppp.SubType,
		"publisher":  ppp.FloAddress,
		"hasDetails": 1,
	}

	var q sq.Sqlizer
	if context.IsEdit {
		q = sq.Update("artifact").SetMap(cv)
	} else {
		q = sq.Insert("artifact").SetMap(cv)
	}

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

	cv = map[string]interface{}{
		"artifactId": artifactId,
		"ns":         ppp.Ns,
		"partyType":  ppp.PartyType,
		"partyRole":  ppp.PartyRole,
	}

	if context.IsEdit {
		q = sq.Update("detailsPropertyParty").SetMap(cv)
	} else {
		q = sq.Insert("detailsPropertyParty").SetMap(cv)
	}

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
