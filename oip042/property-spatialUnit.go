package oip042

import (
	"encoding/json"
	//"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SpatialUnitDetails struct {
	Ns           string          `json:"ns,omitempty"`
	Geometry     *Geometry       `json:"geometry,omitempty"`
	SpatialType  string          `json:"spatialType,omitempty"`
	SpatialUnits []string        `json:"spatialUnits,omitempty"`
	BBox         []float64       `json:"bbox,omitempty"`
	Attrs        json.RawMessage `json:"attrs,omitempty"`
}

type DecimalDegrees struct {
	// ToDo
}
type DegreesMinutesSeconds struct {
	// ToDo
}

type Geometry struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
	// ToDo reconsider how this is structured, should Data be an interface?
	dd   *DecimalDegrees
	dms  *DegreesMinutesSeconds
	text string
}

// ToDo: Currently has infinite recursion,
//var ErrUnknownGeometryType = errors.New("unknown geometry type")
//
//func (u *Geometry) UnmarshalJSON(data []byte) error {
//	var err error
//
// // ToDo as is need a temp struct to unmarshal to then copy values back to Geometry
//	err = json.Unmarshal(data, &u)
//	if err != nil {
//		return err
//	}
//
//	switch u.Type {
//	case "dd":
//		err = json.Unmarshal(u.Data, &u.dd)
//	case "dms":
//		err = json.Unmarshal(u.Data, &u.dms)
//	case "text":
//		u.text = string(u.Data)
//	default:
//		err = ErrUnknownGeometryType
//	}
//	if err != nil {
//		return err
//	}
//	return nil
//}

type PublishPropertySpatialUnit struct {
	PublishArtifact
	SpatialUnitDetails
}

func (ppsu PublishPropertySpatialUnit) Validate(context OipContext) (OipAction, error) {
	err := json.Unmarshal(ppsu.Details, &ppsu.SpatialUnitDetails)
	if err != nil {
		return nil, err
	}

	return ppsu, nil
}

func (ppsu PublishPropertySpatialUnit) Store(context OipContext, dbtx *sqlx.Tx) error {
	j, err := json.Marshal(ppsu)
	if err != nil {
		return err
	}

	q := sq.Insert("artifactPropertySpatialUnit").
		Columns("ns", "spatialType",
			"active", "block", "json", "tags", "timestamp",
			"title", "txid", "type", "subType", "publisher").
		Values(ppsu.Ns, ppsu.SpatialType,
			1, context.BlockHeight, j, ppsu.Info.Tags, ppsu.Timestamp,
			ppsu.Info.Title, context.TxId, ppsu.Type, ppsu.SubType, ppsu.FloAddress)

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

func (ppsu PublishPropertySpatialUnit) MarshalJSON() ([]byte, error) {
	pa := ppsu.PublishArtifact
	buf, err := json.Marshal(ppsu.SpatialUnitDetails)
	if err != nil {
		return nil, err
	}
	pa.Details = buf
	return json.Marshal(pa)
}

func GetAllPropertySpatialUnit(dbtx *sqlx.Tx) ([]interface{}, error) {
	// ToDo combine/simplify these GetAll functions similar to GetById
	q := sq.Select("json", "txid", "publisher").
		From("artifactPropertySpatialUnit").
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

const createPropertySpatialUnitTable = `CREATE TABLE IF NOT EXISTS artifactPropertySpatialUnit
(
  uid            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,

  -- Property-SpatialUnit Fields
  ns             TEXT NOT NULL,
  spatialType    TEXT NOT NULL,

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
