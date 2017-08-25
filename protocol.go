package alexandriaProtocol

import (
	"database/sql"
	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/messages"
	"strings"
)

const Version string = "0.4.1"
const min_block = 1045632

type ParseErrors struct {
	MessageType string
	Error       error
}

func GetMinBlock() int {
	// TODO: find min block from multiple protocols programmatically
	return min_block
}

func Parse(txComment string, txid string, block *flojson.BlockResult, dbtx *sql.Tx) (interface{}, map[string]interface{}, error, []ParseErrors) {

	var pe []ParseErrors

	if strings.HasPrefix(txComment, "text:") {
		txComment = txComment[5:]
	}

	processingBlock := int(block.Height)

	// check for alexandria-media-multipart single protocol (new media multipart tx-comment)
	mms, VerifyMediaMultipartSingleError := messages.VerifyMediaMultipartSingle(txComment, txid, processingBlock)
	if VerifyMediaMultipartSingleError == nil {
		return mms, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"MediaMultipartSingle", VerifyMediaMultipartSingleError})

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

	// check for historian messages
	hm, err := messages.VerifyHistorianMessage([]byte(txComment), processingBlock, dbtx)
	if err == nil {
		return hm, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"Historian", err})

	// check for alexandria-autominer messages
	am, err := messages.VerifyAutominer([]byte(txComment), processingBlock)
	if err == nil {
		return am, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"Autominer", err})

	// check for alexandria-autominer-pool messages
	amp, err := messages.VerifyAutominerPool([]byte(txComment), processingBlock)
	if err == nil {
		return amp, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"AutominerPool", err})

	// check for alexandria-promoter messages
	promoter, err := messages.VerifyPromoter([]byte(txComment), processingBlock)
	if err == nil {
		return promoter, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"Promoter", err})

	// check for alexandria-retailer messages
	retailer, err := messages.VerifyRetailer([]byte(txComment), processingBlock)
	if err == nil {
		return retailer, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"Retailer", err})

	// check for any oip41 data
	oip041, err := messages.VerifyOIP041(txComment, processingBlock)
	if err == nil {
		return oip041, nil, nil, pe
	}
	pe = append(pe, ParseErrors{"OIP041", err})

	return nil, nil, err, pe
}
