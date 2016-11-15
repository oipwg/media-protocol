package messages

import (
	"database/sql"
	"errors"
	"github.com/dloa/media-protocol/utility"
	"math"
	"strconv"
	"strings"
)

var ErrHistorianMessageBadSignature = errors.New("Historian message bad signature")
var ErrHistorianMessageInvalid = errors.New("Historian message invalid")
var ErrHistorianMessagePoolUntrusted = errors.New("Historian message pool untrusted")

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

type hmPool struct {
	address string
	url     string
	version int
}

var hmPools []hmPool = []hmPool{
	{
		// https://github.com/dloa/node-merged-pool/blob/master/lib/pool.js#L39
		// V1 Alexandria.io is signed with FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD
		"FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD",
		"pool.alexandria.io",
		1,
	},
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
	var hm HistorianMessage

	hm.Version = 1
	parts := strings.Split(s, ":")

	if len(parts) != 8 {
		return hm, ErrHistorianMessageInvalid
	}
	hm.Signature = parts[7]

	hm.URL = parts[1]

	p, err := getPool(hm.URL, 1)
	if err != nil {
		return hm, ErrHistorianMessagePoolUntrusted
	}

	i := strings.LastIndex(s, ":")
	if !utility.CheckSignature(p.address, s[i+1:], s[:i]) {
		return hm, ErrHistorianMessageBadSignature
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

func getPool(url string, version int) (hmPool, error) {
	var p hmPool
	for _, p := range hmPools {
		if p.version == version && p.url == url {
			return p, nil
		}
	}
	return p, ErrHistorianMessagePoolUntrusted
}
