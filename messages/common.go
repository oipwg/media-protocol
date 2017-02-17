package messages

import (
	"database/sql"
	"errors"
	"math"
)

var ErrBadSignature = errors.New("Bad signature")
var ErrInvalidAddress = errors.New("Not a valid address")
var ErrInvalidReference = errors.New("Invalid reference transaction")
var ErrNotJSON = errors.New("Not a JSON string")
var ErrTooEarly = errors.New("Too early for a valid message")
var ErrWrongPrefix = errors.New("Wrong prefix for message type")
var ErrNoSignatureEnd = errors.New("no end of signature found, malformed tx-comment")
var ErrNotImplemented = errors.New("Not Implemented")

func CalcAvgArtCost(dbtx *sql.Tx) (float64, int, error) {
	stmtStr := `select m.ArtCost from media as m where m.invalidated = 0 and m.artCost > 0
				union all
				select o.ArtCost from oip_artifact as o where o.invalidated = 0 and o.artCost > 0;`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		return 0.0, 0, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return 0.0, 0, err
	}
	defer rows.Close()

	var artCost float64
	var avgArtCost float64
	var artCount int

	for rows.Next() {
		rows.Scan(&artCost)
		avgArtCost += artCost
		artCount++
	}

	if artCount > 0 {
		avgArtCost = avgArtCost / float64(artCount)
	}

	return avgArtCost, artCount, nil
}

func CalcPubFeeUSD(artCost, avgArtCost float64, artSize int, floPerKb, USDperFLO float64) (float64, float64) {
	var pubFeeComUSD float64

	if artCost <= avgArtCost {
		pubFeeComUSD = artCost
	} else {
		pubFeeComUSD = (math.Log(artCost) - math.Log(avgArtCost)) * (avgArtCost / artCost) * (artCost - avgArtCost)
	}

	// pubFeeComFlo := pubFeeComUSD / USDperFLO
	pubFeeFreeFlo := float64(artSize/1024) * floPerKb
	pubFeeFreeUSD := pubFeeFreeFlo * USDperFLO
	pubFeeUSD := math.Max(pubFeeFreeUSD, pubFeeComUSD)

	return pubFeeUSD, pubFeeUSD / USDperFLO
}

func GetArtCount(dbtx *sql.Tx) (int, error) {
	stmtStr := `SELECT COUNT(*) FROM (select m.uid from media as m where m.invalidated = 0
				union all
				select o.uid from oip_artifact as o where o.invalidated = 0);`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	if err != nil {
		return 0, err
	}

	var artCount int
	row.Scan(&artCount)
	// ToDo: handle err

	return artCount, nil
}
