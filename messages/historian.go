package messages

import (
	"database/sql"
	"errors"
	"github.com/dloa/media-protocol/utility"
	"math"
	"strconv"
	"strings"
)

var ErrHistorianMessageInvalid = errors.New("Historian message invalid")
var ErrHistorianMessageBadSignature = errors.New("Historian message bad signature")

type HistorianMessage struct {
	Version       int
	URL           string
	Mrr_last_10   float64
	Pool_hashrate float64
	Fbd_hashrate  float64
	Fmd_weighted  float64
	Fmd_usd       float64
	Signature     string
}

func StoreHistorianMessage(hm HistorianMessage, dbtx *sql.Tx, txid string, block int) {
	// ToDo: store the data point in the database
}

func VerifyHistorianMessage(b []byte) (HistorianMessage, error) {
	var hm HistorianMessage
	if strings.HasPrefix(string(b), "alexandria-historian-v001") {
		return parseV1(string(b))
	} else {
		return hm, ErrHistorianMessageInvalid
	}
}

func parseV1(s string) (HistorianMessage, error) {
	// https://github.com/dloa/node-merged-pool/blob/master/lib/pool.js#L39
	// V1 Alexandria.io is signed with FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD

	var hm HistorianMessage
	var isAlexandria bool

	hm.Version = 1
	parts := strings.Split(s, ":")

	if len(parts) != 8 {
		return hm, ErrHistorianMessageInvalid
	}
	hm.Signature = parts[7]

	if parts[1] == "pool.alexandria.io" {
		isAlexandria = true
	} else {
		isAlexandria = false
		// Only trust the alexandria pool for now
		return hm, ErrHistorianMessageInvalid
	}
	hm.URL = parts[1]

	if isAlexandria {
		i := strings.LastIndex(s, ":")
		// ToDo: determine proper signature address
		if !utility.CheckSignature("FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD", s[i+1:], s[:i]) {
			return hm, ErrHistorianMessageBadSignature
		}
	}

	for i := 2; i < 7; i++ {
		f, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			f = math.Inf(-1)
		}
		switch i {
		case 2:
			hm.Mrr_last_10 = f
		case 3:
			hm.Pool_hashrate = f
		case 4:
			hm.Fbd_hashrate = f
		case 5:
			hm.Fmd_weighted = f
		case 6:
			hm.Fmd_usd = f
		}
	}

	return hm, nil
}
