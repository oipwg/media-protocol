package oip042

import (
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
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
	j, err := json.Marshal(ppt)
	if err != nil {
		return err
	}

	cv := map[string]interface{}{
		"active":     1,
		"block":      context.BlockHeight,
		"json":       j,
		"tags":       ppt.Info.Tags,
		"unixtime":   ppt.Timestamp,
		"title":      ppt.Info.Title,
		"txid":       context.TxId,
		"type":       ppt.Type,
		"subType":    ppt.SubType,
		"publisher":  ppt.FloAddress,
		"hasDetails": 1,
	}

	var q sq.Sqlizer
	if context.IsEdit {
		q = sq.Update("artifact").SetMap(cv)
	} else {
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

	artifactId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	cv = map[string]interface{}{
		"artifactId":  artifactId,
		"ns":          ppt.Ns,
		"spatialType": ppt.TenureType,
	}

	if context.IsEdit {
		q = sq.Update("detailsPropertyTenure").SetMap(cv)
	} else {
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
