package oip042

import (
	"encoding/json"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

type TomogramDetails struct {
	Date           int64   `json:"date,omitempty"`
	NBCItaxID      int64   `json:"NBCItaxID,omitempty"`
	ArtNotes       string  `json:"artNotes,omitempty"`
	ScopeName      string  `json:"scopeName,omitempty"`
	Roles          string  `json:"roles,omitempty"`
	SpeciesName    string  `json:"speciesName,omitempty"`
	Strain         string  `json:"strain,omitempty"`
	TiltSingleDual int64   `json:"tiltSingleDual,omitempty"`
	Defocus        float64 `json:"defocus,omitempty"`
	Dosage         float64 `json:"dosage,omitempty"`
	TiltConstant   float64 `json:"tiltConstant,omitempty"`
	TiltMin        float64 `json:"tiltMin,omitempty"`
	TiltMax        float64 `json:"tiltMax,omitempty"`
	TiltStep       float64 `json:"tiltStep,omitempty"`
	Magnification  float64 `json:"magnification,omitempty"`
	Emdb           string  `json:"emdb,omitempty"`
	Microscopist   string  `json:"microscopist,omitempty"`
	Institution    string  `json:"institution,omitempty"`
	Lab            string  `json:"lab,omitempty"`
	Sid            string  `json:"sid,omitempty"`
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

	if len(pt.TomogramDetails.SpeciesName) == 0 {
		return nil, errors.New("tomogram: missing Species Name")
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

	cv := map[string]interface{}{
		"active":     1,
		"block":      context.BlockHeight,
		"json":       j,
		"tags":       pt.Info.Tags,
		"unixtime":   pt.Timestamp,
		"title":      pt.Info.Title,
		"txid":       context.TxId,
		"type":       pt.Type,
		"subType":    pt.SubType,
		"publisher":  pt.FloAddress,
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
		"artifactId":     artifactId,
		"ScanDate":       pt.Date,
		"NBCItaxID":      pt.NBCItaxID,
		"ArtNotes":       pt.ArtNotes,
		"ScopeName":      pt.ScopeName,
		"SpeciesName":    pt.SpeciesName,
		"TiltSingleDual": pt.TiltSingleDual,
		"Defocus":        pt.Defocus,
		"Magnification":  pt.Magnification,
		"Emdb":           pt.Emdb,
		"SwAcquisition":  "",
		"SwProcess":      "",
		"Institution":    pt.Institution,
		"Lab":            pt.Lab,
		"sid":            pt.Sid,
	}

	if context.IsEdit {
		q = sq.Update("detailsResearchTomogram").SetMap(cv)
	} else {
		q = sq.Insert("detailsResearchTomogram").SetMap(cv)
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
  ArtNotes       TEXT    NOT NULL,
  ScopeName      TEXT    NOT NULL,
  SpeciesName    TEXT    NOT NULL,
  TiltSingleDual TEXT    NOT NULL,
  Defocus        TEXT    NOT NULL,
  Magnification  TEXT    NOT NULL,
  SwAcquisition  TEXT    NOT NULL,
  SwProcess      TEXT    NOT NULL,
  Emdb           TEXT    NOT NULL,
  Institution    TEXT    NOT NULL,
  Lab            TEXT    NOT NULL,
  sid            TEXT    NOT NULL,
  CONSTRAINT detailsResearchTomogram_artifactId_uid_fk FOREIGN KEY (artifactId) REFERENCES artifact (uid) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS detailsResearchTomogram_artifactId_uindex ON detailsResearchTomogram (artifactId);
CREATE INDEX IF NOT EXISTS detailsResearchTomogram_speciesName_index ON detailsResearchTomogram (SpeciesName);
CREATE INDEX IF NOT EXISTS detailsResearchTomogram_emdb_index ON detailsResearchTomogram (Emdb);
`
