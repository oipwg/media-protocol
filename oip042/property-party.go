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
		"json":      j,
		"tags":      ppp.Info.Tags,
		"unixtime":  ppp.Timestamp,
		"title":     ppp.Info.Title,
		"type":      ppp.Type,
		"subType":   ppp.SubType,
		"publisher": ppp.FloAddress,
	}

	var q sq.Sqlizer
	if context.IsEdit {
		q = sq.Update("artifact").SetMap(cv).Where(sq.Eq{"txid": context.Reference})
	} else {
		// these values are only set on publish
		cv["active"] = 1
		cv["txid"] = context.TxId
		cv["block"] = context.BlockHeight
		cv["hasDetails"] = 1
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

	cv = map[string]interface{}{
		"ns":        ppp.Ns,
		"partyType": ppp.PartyType,
		"partyRole": ppp.PartyRole,
	}

	if context.IsEdit {
		q = sq.Update("detailsPropertyParty").SetMap(cv).Where(sq.Eq{"artifactId": context.ArtifactId})
	} else {
		artifactId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		cv["artifactId"] = artifactId
		context.ArtifactId = artifactId
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
