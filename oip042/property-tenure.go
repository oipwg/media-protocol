package oip042

import (
	"encoding/json"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"strings"
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

func (ppt PublishPropertyTenure) Store(context OipContext) error {
	index := false
	if len(context.IndexTypes) == 0 {
		index = true
	} else {
		for _, t := range context.IndexTypes {
			if strings.ToLower(ppt.Type) == t {
				index = true
				break
			}
		}
	}
	if !index {
		return errors.New("not indexed due to IndexedTypes config")
	}

	j, err := json.Marshal(ppt)
	if err != nil {
		return err
	}

	cv := map[string]interface{}{
		"json":      j,
		"tags":      ppt.Info.Tags,
		"unixtime":  ppt.Timestamp,
		"title":     ppt.Info.Title,
		"type":      ppt.Type,
		"subType":   ppt.SubType,
		"publisher": ppt.FloAddress,
	}

	var q sq.Sqlizer
	if context.IsEdit {
		cv["txid"] = context.Reference
		q = sq.Update("artifact").SetMap(cv).Where(sq.Eq{"txid": context.Reference})
	} else {
		// these values are only set on publish
		cv["active"] = 1
		cv["txid"] = context.TxId
		cv["block"] = context.BlockHeight
		cv["hasDetails"] = 1
		q = sq.Insert("artifact").SetMap(cv)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	res, err := context.DbTx.Exec(sql, args...)
	if err != nil {
		return err
	}

	cv = map[string]interface{}{
		"ns":         ppt.Ns,
		"tenureType": ppt.TenureType,
	}

	if context.IsEdit {
		q = sq.Update("detailsPropertyTenure").SetMap(cv).Where(sq.Eq{"artifactId": context.ArtifactId})
	} else {
		artifactId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		cv["artifactId"] = artifactId
		context.ArtifactId = artifactId
		q = sq.Insert("detailsPropertyTenure").SetMap(cv)
	}

	sql, args, err = q.ToSql()
	if err != nil {
		return err
	}

	_, err = context.DbTx.Exec(sql, args...)
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

const createPropertyTenureTable = `
-- Property-Tenure details
CREATE TABLE IF NOT EXISTS detailsPropertyTenure
(
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    artifactId INT NOT NULL,
    ns TEXT NOT NULL,
    tenureType TEXT NOT NULL,
    CONSTRAINT detailsPropertyTenure_artifactId_uid_fk FOREIGN KEY (artifactId) REFERENCES artifact (uid) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS detailsPropertyTenure_artifactId_uindex  ON detailsPropertyTenure (artifactId);
CREATE INDEX IF NOT EXISTS detailsPropertyTenure_tenureType_index  ON detailsPropertyTenure (tenureType);
`
