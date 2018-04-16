package oip042

import (
	"encoding/json"
	"github.com/oipwg/media-protocol/utility"
	"strconv"
	"strings"
)

type RegisterAffiliate struct {
	Alias        string            `json:"alias"`
	Timestamp    int64             `json:"timestamp"`
	FloAddress   string            `json:"floAddress"`
	Addresses    []PaymentAddress  `json:"addresses"`
	ShortMW      []string          `json:"shortMW"`
	Verification map[string]string `json:"verification"`
	Version      int32             `json:"version"`
}

type EditAffiliate struct {
	ArtifactID string          `json:"artifactID"`
	Timestamp  int64           `json:"timestamp"`
	Patch      json.RawMessage `json:"patch"`
}

func (affiliate *EditAffiliate) Store(context OipContext) error {
	panic("implement me")
}

func (ra RegisterAffiliate) Validate(context OipContext) error {
	v := []string{ra.FloAddress, strconv.FormatInt(int64(ra.Version), 10), strconv.FormatInt(ra.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(ra.FloAddress, context.signature, preImage)
	if !sigOk {
		return ErrBadSignature
	}

	return nil
}

func (affiliate *EditAffiliate) Validate(context OipContext) (OipAction, error) {
	return affiliate, ErrNotImplemented
}
