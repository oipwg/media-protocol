package messages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/oipwg/media-protocol/oip042"
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
	//s, err := o.GetJSON()
	//if err != nil {
	//	return err
	//}

	artCost := 0   // o.GetArtCost()
	pubFeeUSD := 0 // o.GetPubFeeUSD(dbtx)

	//a := o.Artifact
	//pa := oip042.PublishArtifact{
	//	Type:    a.Type,
	//	SubType: "",
	//	Info: &oip042.ArtifactInfo{
	//		Title:       a.Info.Title,
	//		Tags:        a.Info.Tags,
	//		Description: a.Info.Description,
	//		Year:        a.Info.Year,
	//		NSFW:        a.Info.NSFW,
	//	},
	//	FloAddress: a.Publisher,
	//	Timestamp:  a.Timestamp,
	//	Storage: &oip042.ArtifactStorage{
	//		Location: a.Storage.Location,
	//		Network:  a.Storage.Network,
	//	},
	//	Payment: &oip042.ArtifactPayment{
	//		Scale:a.Payment.Scale,
	//		MaxDiscount: a.Payment.MaxDiscount,
	//		Tokens: nil,
	//		Fiat: a.Payment.Fiat,
	//		Addresses: nil,
	//
	//		},
	//	Signature: o.Signature,
	//}
	//
	//for _, f := range a.Storage.Files {
	//	pa.Storage.Files = append(pa.Storage.Files, oip042.ArtifactFiles{
	//		DisallowBuy:  f.DisallowBuy,
	//		Dname:        f.Dname,
	//		Duration:     f.Duration,
	//		Fname:        f.Fname,
	//		Fsize:        f.Fsize,
	//		MinPlay:      f.MinPlay,
	//		SugPlay:      f.SugPlay,
	//		Promo:        f.Promo,
	//		Retail:       f.Retail,
	//		PtpFT:        f.PtpFT,
	//		PtpDT:        f.PtpDT,
	//		PtpDA:        f.PtpDA,
	//		Type:         f.Type,
	//		TokenlyID:    f.TokenlyID,
	//		DisallowPlay: f.DisallowPlay,
	//		MinBuy:       f.MinBuy,
	//		SugBuy:       f.SugBuy,
	//		SubType:      f.SubType,
	//	})
	//}

	var pa oip042.PublishArtifact

	jt1, err := json.MarshalIndent(o.Artifact, " ", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(jt1))
	err = json.Unmarshal(jt1, &pa)

	j, err := json.MarshalIndent(pa, " ", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(j))

	return nil
	q := squirrel.Insert("artifact").
		Columns("active", "block", "json", "tags", "unixtime",
			"title", "txid", "type", "subType", "publisher", "hasDetails",
			"artCost", "pubFeeUsd", "artSize").
		Values(1, block, string(j), o.Artifact.Info.Tags, o.Artifact.Timestamp,
			o.Artifact.Info.Title, txid, o.Artifact.Type, "", o.Artifact.Publisher, 0,
			artCost, pubFeeUSD, o.artSize)

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = dbtx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
