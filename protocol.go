package alexandriaProtocol

import (
	"github.com/dloa/media-protocol/messages"
	"strings"
)

const Version string = "0.4.1"
const min_block = 1045632

func GetMinBlock() int {
	// TODO: find min block from multiple protocols programmatically
	return min_block
}

func Parse(txComment string, txid string, processingBlock int) (interface{}, map[string]interface{}, error) {

	if strings.HasPrefix(txComment, "text:") {
		txComment = txComment[5:]
	}

	// check for alexandria-media-multipart single protocol (new media multipart tx-comment)
	mms, VerifyMediaMultipartSingleError := messages.VerifyMediaMultipartSingle(txComment, txid, processingBlock)
	if VerifyMediaMultipartSingleError == nil {
		return mms, nil, nil
	}

	// check for alexandria-media protocol (new media)
	media, jmap, VerifyMediaErr := messages.VerifyMedia([]byte(txComment))
	if VerifyMediaErr == nil {
		return media, jmap, nil
	}

	// check for alexandria-media protocol (new publisher)
	publisher, VerifyPublisherErr := messages.VerifyPublisher([]byte(txComment))
	if VerifyPublisherErr == nil {
		return publisher, nil, nil
	}

	// check for alexandria-media deactivation (new deactivation message)
	deactivation, VerifyDeactivationError := messages.VerifyDeactivation([]byte(txComment))
	if VerifyDeactivationError == nil {
		return deactivation, nil, nil
	}

	// check for alexandria-historian messages
	hm, err := messages.VerifyHistorianMessage([]byte(txComment), processingBlock)
	if err == nil {
		return hm, nil, nil
	}

	// check for any oip41 data
	oip041, err := messages.VerifyOIP041(txComment, processingBlock)
	if err == nil {
		return oip041, nil, nil
	}

	return nil, nil, err
}
