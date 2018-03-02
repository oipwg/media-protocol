package messages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/oipwg/media-protocol/utility"
	"math"
	"strconv"
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
	if len(o.Storage.Network) == 0 {
		return errors.New("Artifact storage network is required")
	}
	if len(o.Storage.Location) == 0 {
		return errors.New("Artifact storage location is required")
	}
	if len(o.Storage.Files) == 0 {
		return errors.New("Artifact must contain at least one file")
	}
	if o.Payment.MaxDiscount < 0 {
		return errors.New("maxdisc must be >= 0")
	}
	return nil
}

func (o Oip041) GetArtCost() float64 {
	var totMinPlay float64 = 0
	var totSugPlay float64 = 0
	var totMinBuy float64 = 0
	var totSugBuy float64 = 0

	for _, f := range o.Artifact.Storage.Files {
		if !f.DisallowPlay {
			totMinPlay += math.Abs(f.MinPlay)
		}
		if !f.DisallowPlay {
			totSugPlay += math.Abs(f.SugPlay)
		}
		if !f.DisallowBuy {
			totMinBuy += math.Abs(f.MinBuy)
		}
		if !f.DisallowBuy {
			totSugBuy += math.Abs(f.SugBuy)
		}
	}

	avg := (totMinPlay + totSugPlay + totMinBuy + totSugBuy) / 4

	splitScale := strings.Split(o.Artifact.Payment.Scale, ":")

	if len(splitScale) == 2 {
		scales := [2]float64{}
		if s, err := strconv.ParseFloat(splitScale[0], 64); err == nil {
			scales[0] = s
		}
		if s, err := strconv.ParseFloat(splitScale[1], 64); err == nil {
			scales[1] = s
		}

		avg = avg * scales[1] / scales[0]
	}

	return avg
}

func (o Oip041) GetPubFeeUSD(dbtx *sql.Tx) float64 {
	artCost := o.GetArtCost()

	// ToDo: Fetch proper values.
	avgArtCost, _, _ := CalcAvgArtCost(dbtx)
	floPerKb := 0.01
	USDperFLO := 0.004564

	pubFeeUSD, _ := CalcPubFeeUSD(artCost, avgArtCost, o.artSize, floPerKb, USDperFLO)

	return pubFeeUSD
}

func StoreOIP041Artifact(o Oip041, txid string, block int, dbtx *sql.Tx) error {
	// store in database
	stmtStr := `INSERT INTO 'oip_artifact'
		('active','block','json','tags','timestamp',
		'title','txid','type','year','publisher', 'artCost', 'artSize', 'pubFeeUSD', 'nsfw')
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

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
	pubFeeUSD := o.GetPubFeeUSD(dbtx)

	_, err = stmt.Exec(1, block, s, o.Artifact.Info.Tags,
		o.Artifact.Timestamp, o.Artifact.Info.Title, txid, o.Artifact.Type,
		o.Artifact.Info.Year, o.Artifact.Publisher, artCost, o.artSize, pubFeeUSD, o.Artifact.Info.NSFW)
	if err != nil {
		return err
	}

	return nil
}
