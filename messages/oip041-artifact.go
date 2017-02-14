package messages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/oipwg/media-protocol/utility"
	"math"
	"strings"
)

func (o Oip041) GetJSON() (string, error) {
	// ToDo: remove redundant Storage items, potentially cache?
	var s string

	b, err := json.Marshal(o)
	s = string(b)

	return s, err
}

func (o Oip041Artifact) CheckRequiredFields() error {
	if !utility.CheckAddress(o.Publisher) {
		return errors.New("Publisher not a valid address")
	}
	if len(o.Type) == 0 {
		return errors.New("Artifact type is required")
	}
	if len(o.Info.Title) == 0 {
		return errors.New("Artifact title is required")
	}
	if len(o.Info.Description) == 0 {
		return errors.New("Artifact title is required")
	}
	if o.Info.Year <= 0 {
		return errors.New("Artifact year is required")
	}
	if len(o.Storage.Network) == 0 {
		return errors.New("Artifact storage network is required")
	}
	if len(o.Storage.Location) == 0 {
		return errors.New("Artifact storage location is required")
	}
	if len(o.Storage.Files) == 0 {
		return errors.New("Artifact must contain at least one file")
	}
	return nil
}

func (o Oip041) GetArtCost() float64 {
	var totMinPlay float64 = 0
	var totSugBuy float64 = 0

	for _, f := range o.Artifact.Storage.Files {
		if f.DisallowPlay != 0 {
			totMinPlay += math.Abs(f.MinPlay)
		}
		if f.DisallowBuy != 0 {
			totSugBuy += math.Abs(f.SugBuy)
		}
	}

	avg := (totMinPlay + totSugBuy) / 2

	return avg
}

func StoreOIP041Artifact(o Oip041, txid string, block int, dbtx *sql.Tx) error {
	// store in database
	stmtStr := `INSERT INTO 'oip_artifact'
		('active','block','json','tags','timestamp',
		'title','txid','type','year','publisher', 'artCost')
		VALUES (?,?,?,?,?,?,?,?,?,?,?);`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	s, err := o.GetJSON()
	if err != nil {
		return nil
	}

	artCost := o.GetArtCost()

	_, err = stmt.Exec(1, block, s, strings.Join(o.Artifact.Info.ExtraInfo.Tags, ","),
		o.Artifact.Timestamp, o.Artifact.Info.Title, txid, o.Artifact.Type,
		o.Artifact.Info.Year, o.Artifact.Publisher, artCost)
	if err != nil {
		return err
	}

	return nil
}
