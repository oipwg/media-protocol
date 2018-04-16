package oip042

import (
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

type TomogramDetails struct {
	Date           int64   `json:"date,omitempty"`
	NBCItaxID      int64   `json:"NBCItaxID,omitempty"`
	Etdbid         int64   `json:"etdbid,omitempty"`
	ArtNotes       string  `json:"artNotes,omitempty"`
	ScopeName      string  `json:"scopeName,omitempty"`
	Roles          string  `json:"roles,omitempty"`
	SpeciesName    string  `json:"speciesName,omitempty"`
	Strain         string  `json:"strain,omitempty"`
	TiltSingleDual int64   `json:"tiltSingleDual,omitempty"`
	Defocus        float64 `json:"defocus,omitempty"`
	Magnification  float64 `json:"magnification,omitempty"`
	Emdb           string  `json:"emdb,omitempty"`
	Microscopist   string  `json:"microscopist,omitempty"`
	Institution    string  `json:"institution,omitempty"`
	Lab            string  `json:"lab,omitempty"`
}

type PublishTomogram struct {
	PublishArtifact
	TomogramDetails
}

func (pt PublishTomogram) Validate(context OipContext) (OipAction, error) {
	err := json.Unmarshal(pt.Details, &pt.TomogramDetails)
	if err != nil {
		return nil, err
	}

	//if len(pt.TomogramDetails.ScopeName) == 0 {
	//	return nil, errors.New("tomogram: missing Scope Name")
	//}
	if len(pt.TomogramDetails.SpeciesName) == 0 {
		return nil, errors.New("tomogram: missing Species Name")
	}
	if len(pt.TomogramDetails.Strain) == 0 {
		return nil, errors.New("tomogram: missing Strain")
	}
	if pt.Date <= 0 {
		return nil, errors.New("tomogram: invalid Date")
	}
	if pt.NBCItaxID <= 0 {
		return nil, errors.New("tomogram: invalid NBCItaxID")
	}

	return pt, nil
}

func (pt PublishTomogram) Store(context OipContext) error {

	j, err := json.Marshal(pt)
	if err != nil {
		return err
	}

	q := sq.Insert("artifact").
		Columns("active", "block", "json", "tags", "unixtime",
			"title", "txid", "type", "subType", "publisher", "hasDetails").
		Values(1, context.BlockHeight, j, pt.Info.Tags, pt.Timestamp,
			pt.Info.Title, context.TxId, pt.Type, pt.SubType, pt.FloAddress, 1)

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

	q = sq.Insert("detailsResearchTomogram").
		Columns("artifactId", "ScanDate", "NBCItaxID", "Etdbid", "ArtNotes",
			"ScopeName", "SpeciesName", "TiltSingleDual", "Defocus",
			"Magnification", "Emdb", "SwAcquisition", "SwProcess").
		Values(artifactId, pt.Date, pt.NBCItaxID, pt.Etdbid, pt.ArtNotes,
			pt.ScopeName, pt.SpeciesName, pt.TiltSingleDual, pt.Defocus,
			pt.Magnification, pt.Emdb, "", "")

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

func (pt PublishTomogram) MarshalJSON() ([]byte, error) {
	pa := pt.PublishArtifact
	buf, err := json.Marshal(pt.TomogramDetails)
	if err != nil {
		return nil, err
	}
	pa.Details = buf
	return json.Marshal(pa)
}

const createTomogramTable = `
-- Research-Tomogram details
CREATE TABLE IF NOT EXISTS detailsResearchTomogram
(
  uid            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  artifactId     INT     NOT NULL,
  ScanDate       INT     NOT NULL,
  NBCItaxID      INT     NOT NULL,
  Etdbid         TEXT    NOT NULL,
  ArtNotes       TEXT    NOT NULL,
  ScopeName      TEXT    NOT NULL,
  SpeciesName    TEXT    NOT NULL,
  TiltSingleDual TEXT    NOT NULL,
  Defocus        TEXT    NOT NULL,
  Magnification  TEXT    NOT NULL,
  SwAcquisition  TEXT    NOT NULL,
  SwProcess      TEXT    NOT NULL,
  Emdb           TEXT    NOT NULL,
  CONSTRAINT detailsResearchTomogram_artifactId_uid_fk FOREIGN KEY (artifactId) REFERENCES artifact (uid) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS detailsResearchTomogram_artifactId_uindex ON detailsResearchTomogram (artifactId);
CREATE INDEX IF NOT EXISTS detailsResearchTomogram_speciesName_index ON detailsResearchTomogram (SpeciesName);
CREATE INDEX IF NOT EXISTS detailsResearchTomogram_emdb_index ON detailsResearchTomogram (Emdb);
`
