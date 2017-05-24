package messages

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/utility"
	"log"
	"math"
	"strconv"
	"strings"
)

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
	address  string
	maxValid int
	minValid int
	url      string
	version  int
}

type hmPoolList []hmPool

var hmPools hmPoolList = hmPoolList{
	{
		// https://github.com/dloa/node-merged-pool/blob/2a3f124/lib/pool.js#L39
		// V1 Alexandria.io is signed with
		"FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU",
		0,
		1974560,
		"pool.alexandria.io",
		1,
	},
	{
		// For a period there was no signature, but they are trusted
		"",
		1974560,
		1887692,
		"pool.alexandria.io",
		1,
	},
	{
		// https://github.com/dloa/node-merged-pool/blob/fcd6ab59/lib/pool.js#L39
		// V1 Alexandria.io is signed with FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD
		"FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD",
		1887692,
		0,
		"pool.alexandria.io",
		1,
	},
}

func StoreHistorianMessage(hm HistorianMessage, dbtx *sql.Tx, txid string, block *flojson.BlockResult) {
	// store in database
	stmtStr := `insert into historian (txid, block, blockTime, active, version,` +
		` url, mrrLast10, poolHashrate, fbdHashrate, fmdWeighted, fmdUSD, signature)` +
		` values (?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		fmt.Println("exit 200")
		log.Fatal(err)
	}

	_, stmterr := stmt.Exec(txid, block.Height, block.Time, hm.Version, hm.URL, hm.Mrr_last_10,
		hm.Pool_hashrate, hm.Fbd_hashrate, hm.Fmd_weighted, hm.Fmd_usd, hm.Signature)
	if err != nil {
		fmt.Println("exit 201")
		log.Fatal(stmterr)
	}

	stmt.Close()
}

func VerifyHistorianMessage(b []byte, block int) (HistorianMessage, error) {
	var hm HistorianMessage
	if strings.HasPrefix(string(b), "alexandria-historian-v001") {
		return parseV1(string(b), block)
	} else {
		return hm, ErrWrongPrefix
	}
}

func parseV1(s string, block int) (HistorianMessage, error) {
	var hm HistorianMessage

	hm.Version = 1
	parts := strings.Split(s, ":")

	if len(parts) < 6 || len(parts) > 9 {
		return hm, ErrHistorianMessageInvalid
	}
	if len(parts) == 8 {
		hm.Signature = parts[7]
	}
	hm.URL = parts[1]

	p, err := hmPools.GetPool(hm.URL, block, 1)
	if err != nil {
		return hm, err
	}

	// If there's no defined address there is no signature to check
	if p.address != "" {
		i := strings.LastIndex(s, ":")
		val, _ := utility.CheckSignature(p.address, s[i+1:], s[:i])
		if !val {
			return hm, ErrBadSignature
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

func (hmp hmPoolList) GetPool(url string, block int, version int) (hmPool, error) {
	var p hmPool
	for _, p := range hmp {
		if p.version == version && p.url == url && p.minValid <= block &&
			(p.maxValid > block || p.maxValid == 0) {
			return p, nil
		}
	}
	return p, ErrHistorianMessagePoolUntrusted
}
