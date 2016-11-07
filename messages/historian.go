package messages

import (
	"database/sql"
	"errors"
)

// Sample:
// alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA=

type HistorianMessage struct {
	Version       string
	UTL           string
	Mrr_last_10   string
	Pool_hashrate string
	Fbd_hashrate  string
	Fmd_weighted  string
	Fmd_usd       string
}

func StoreHistorianMessage(hm HistorianMessage, dbtx *sql.Tx, txid string, block int, multipart int) {
	// ToDo: store the data point in the database
}

func VerifyHistorianMessage(b []byte) (HistorianMessage, error) {
	// ToDo: Parse/Validate and return a HistorianMessage
	return nil, errors.New("Not implemented")
}
