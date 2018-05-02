package oip042

import (
	"encoding/json"
	//"errors"
	sq "github.com/Masterminds/squirrel"
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

func (ppsu PublishPropertySpatialUnit) Store(context OipContext) error {
	j, err := json.Marshal(ppsu)
	if err != nil {
		return err
	}

	cv := map[string]interface{}{
		"json":      j,
		"tags":      ppsu.Info.Tags,
		"unixtime":  ppsu.Timestamp,
		"title":     ppsu.Info.Title,
		"type":      ppsu.Type,
		"subType":   ppsu.SubType,
		"publisher": ppsu.FloAddress,
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
		"ns":          ppsu.Ns,
		"spatialType": ppsu.SpatialType,
	}

	if context.IsEdit {
		q = sq.Update("detailsPropertySpatialUnit").SetMap(cv).Where(sq.Eq{"artifactId": context.ArtifactId})
	} else {
		artifactId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		cv["artifactId"] = artifactId
		context.ArtifactId = artifactId
		q = sq.Insert("detailsPropertySpatialUnit").SetMap(cv)
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

func (ppsu PublishPropertySpatialUnit) MarshalJSON() ([]byte, error) {
	pa := ppsu.PublishArtifact
	buf, err := json.Marshal(ppsu.SpatialUnitDetails)
	if err != nil {
		return nil, err
	}
	pa.Details = buf
	return json.Marshal(pa)
}

const createPropertySpatialUnitTable = `
-- Property-SpatialUnit details
CREATE TABLE IF NOT EXISTS detailsPropertySpatialUnit
(
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    artifactId INT NOT NULL,
    ns TEXT NOT NULL,
    spatialType TEXT NOT NULL,
    CONSTRAINT detailsPropertySpatialUnit_artifactId_uid_fk FOREIGN KEY (artifactId) REFERENCES artifact (uid) ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS detailsPropertySpatialUnit_artifactId_uindex  ON detailsPropertySpatialUnit (artifactId);
CREATE INDEX IF NOT EXISTS detailsPropertySpatialUnit_spatialType_index  ON detailsPropertySpatialUnit (spatialType);
`
