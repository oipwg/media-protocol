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

	q := sq.Insert("artifact").
		Columns("active", "block", "json", "tags", "unixtime",
			"title", "txid", "type", "subType", "publisher", "hasDetails").
		Values(1, context.BlockHeight, j, ppt.Info.Tags, ppt.Timestamp,
			ppt.Info.Title, context.TxId, ppt.Type, ppt.SubType, ppt.FloAddress, 1)

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	res, err := dbtx.Exec(sql, args...)
	if err != nil {
		return err
	}

	artifactId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	q = sq.Insert("detailsPropertyTenure").
		Columns("artifactId", "ns", "tenureType").
		Values(artifactId, ppt.Ns, ppt.TenureType)

	sql, args, err = q.ToSql()
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
