package alexandriaProtocol

import (
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/messages"
	"github.com/oipwg/media-protocol/oip042"
	"github.com/oipwg/media-protocol/utility"
	"strings"
)

const Version string = "0.4.1"
const min_block = 1045632

type ParseErrors struct {
	MessageType string
	Error       error
}

var ErrUnknownType = errors.New("unknown type")

func GetMinBlock() int {
	// TODO: find min block from multiple protocols programmatically
	if utility.Testnet() {
		return 1
	}
	return min_block
}

func Parse(tx *flojson.TxRawResult, txid string, block *flojson.BlockResult, dbtx *sqlx.Tx) (interface{}, map[string]interface{}, error, []ParseErrors) {

	var pe []ParseErrors

	txComment := tx.TxComment
	if strings.HasPrefix(txComment, "text:") {
		txComment = txComment[5:]
	}

	if strings.HasPrefix(txComment, "json:") {
		txComment = txComment[5:]
		res, err := ParseJson(tx, txComment, txid, block, dbtx)
		return res, nil, err, pe
	}

	processingBlock := int(block.Height)

	// check for alexandria-media-multipart single protocol (new media multipart tx-comment)
	mms, VerifyMediaMultipartSingleError := messages.VerifyMediaMultipartSingle(txComment, txid, processingBlock)
	if VerifyMediaMultipartSingleError == nil {
		return mms, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"MediaMultipartSingle", VerifyMediaMultipartSingleError})

	// check for historian messages
	if len(tx.Vin) > 0 && tx.Vin[0].IsCoinBase() {
		hm, err := messages.VerifyHistorianMessage([]byte(txComment), processingBlock, dbtx.Tx)
		if err == nil {
			return hm, nil, nil, pe
		}
		pe = append(pe, ParseErrors{"Historian", err})
	}

	if !utility.IsJSON(txComment) {
		pe = append(pe, ParseErrors{"JSON", messages.ErrNotJSON})
		return nil, nil, messages.ErrNotJSON, pe
	}

	// check for alexandria-media protocol (new media)
	media, jmap, VerifyMediaErr := messages.VerifyMedia([]byte(txComment))
	if VerifyMediaErr == nil {
		return media, jmap, nil, pe
	}
	pe = append(pe, ParseErrors{"AlexandriaMedia", VerifyMediaErr})

	// check for alexandria-media protocol (new publisher)
	publisher, VerifyPublisherErr := messages.VerifyPublisher([]byte(txComment))
	if VerifyPublisherErr == nil {
		return publisher, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"Publisher", VerifyPublisherErr})

	// check for alexandria-media deactivation (new deactivation message)
	deactivation, VerifyDeactivationError := messages.VerifyDeactivation([]byte(txComment))
	if VerifyDeactivationError == nil {
		return deactivation, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"Deactivation", VerifyDeactivationError})

	// Only tests have actually made it to chain, all future items will be of oip042 format
	//// check for alexandria-autominer messages
	//am, err := messages.VerifyAutominer([]byte(txComment), processingBlock)
	//if err == nil {
	//	return am, nil, nil, pe
	//}
	//pe = append(pe, ParseErrors{"Autominer", err})
	//
	//// check for alexandria-autominer-pool messages
	//amp, err := messages.VerifyAutominerPool([]byte(txComment), processingBlock)
	//if err == nil {
	//	return amp, nil, nil, pe
	//}
	//pe = append(pe, ParseErrors{"AutominerPool", err})
	//
	//// check for alexandria-promoter messages
	//promoter, err := messages.VerifyPromoter([]byte(txComment), processingBlock)
	//if err == nil {
	//	return promoter, nil, nil, pe
	//}
	//pe = append(pe, ParseErrors{"Promoter", err})
	//
	//// check for alexandria-retailer messages
	//retailer, err := messages.VerifyRetailer([]byte(txComment), processingBlock)
	//if err == nil {
	//	return retailer, nil, nil, pe
	//}
	//pe = append(pe, ParseErrors{"Retailer", err})

	// check for any oip41 data
	oip041, err := messages.VerifyOIP041(txComment, processingBlock)
	if err == nil {
		return oip041, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"OIP041", err})

	//res, err := ParseJson(tx, txComment, txid, block, dbtx)
	//if err == nil {
	//	return res, nil, err, pe
	//}
	//pe = append(pe, ParseErrors{"OIP042", err})

	return nil, nil, err, pe
}

func ParseJson(tx *flojson.TxRawResult, txComment string, txid string, block *flojson.BlockResult, dbtx *sqlx.Tx) (interface{}, error) {
	type supportedJsonTypes struct {
		Oip042 *oip042.Oip042 `json:"oip042,omitempty"`
	}

	var dec supportedJsonTypes
	err := json.Unmarshal([]byte(txComment), &dec)
	if err != nil {
		return nil, messages.ErrNotJSON
	}

	// only process the first match, disregard remaining
	// otherwise there's order of operations to consider

	if dec.Oip042 != nil {
		return dec.Oip042.ValidateIncoming(tx, txComment, txid, block, dbtx)
	}

	// if dec.type2 != nil {
	//   err := dec.type2.ValidateIncoming(...)
	//   return dec.type2, err
	// }

	return nil, ErrUnknownType
}
