package oip042

import (
	"encoding/json"
	"errors"

	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type TomogramDetails struct {
	Date           int64   `json:"date,omitempty"`
	NBCItaxID      int64   `json:"NBCItaxID,omitempty"`
	Etdbid         int64   `json:"etdbid,omitempty"`
	ArtNotes       string  `json:"artNotes,omitempty"`
	ScopeName      string  `json:"scopeName,omitempty"`
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

func (pt PublishTomogram) Store(context OipContext, dbtx *sqlx.Tx) error {

	j, err := json.Marshal(pt)
	if err != nil {
		return err
	}

	q := sq.Insert("artifactsResearchTomogram").
		Columns("ScanDate", "NBCItaxID", "Etdbid", "ArtNotes",
			"ScopeName", "SpeciesName", "TiltSingleDual", "Defocus",
			"Magnification", "Emdb",
			"active", "block", "json", "tags", "timestamp",
			"title", "txid", "type", "subType", "publisher").
		Values(pt.Date, pt.NBCItaxID, pt.Etdbid, pt.ArtNotes,
			pt.ScopeName, pt.SpeciesName, pt.TiltSingleDual, pt.Defocus,
			pt.Magnification, pt.Emdb,
			1, context.BlockHeight, j, pt.Info.Tags, pt.Timestamp,
			pt.Info.Title, context.TxId, pt.Type, pt.SubType, pt.FloAddress)

	sql, args, err := q.ToSql()
	fmt.Println("Storing Tomogram Publish")
	fmt.Println(sql)
	fmt.Println(args)
	if err != nil {
		return err
	}

	_, err = dbtx.Exec(sql, args...)
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

func GetAllTomograms(dbtx *sqlx.Tx) ([]interface{}, error) {
	q := sq.Select("json", "txid", "publisher").From("artifactsResearchTomogram").Where("active = ?", 1)
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

const createTomogramTable = `CREATE TABLE IF NOT EXISTS artifactsResearchTomogram
(
  uid            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,

  -- Tomogram Fields
  ScanDate       INT     NOT NULL,
  NBCItaxID      INT     NOT NULL,
  Etdbid         TEXT    NOT NULL,
  ArtNotes       TEXT,
  ScopeName      TEXT    NOT NULL,
  SpeciesName    TEXT    NOT NULL,
  TiltSingleDual TEXT    NOT NULL,
  Defocus        TEXT,
  Magnification  TEXT,
  SwAcquisition  TEXT,
  SwProcess      TEXT,
  Emdb           TEXT,

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
