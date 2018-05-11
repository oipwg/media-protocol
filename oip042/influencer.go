package oip042

import (
	"encoding/json"
	"github.com/oipwg/media-protocol/utility"
	"strconv"
	"strings"
)

type RegisterInfluencer struct {
	Alias        string            `json:"alias"`
	Timestamp    int64             `json:"timestamp"`
	FloAddress   string            `json:"floAddress"`
	Addresses    []PaymentAddress  `json:"addresses"`
	ShortMW      []string          `json:"shortMW"`
	Verification map[string]string `json:"verification"`
	Version      int32             `json:"version"`
	Signature    string            `json:"signature"`
}

type EditInfluencer struct {
	ArtifactID string          `json:"artifactID"`
	Timestamp  int64           `json:"timestamp"`
	Patch      json.RawMessage `json:"patch"`
	Signature  string          `json:"signature"`
}

func (ei *EditInfluencer) Store(context OipContext) error {
	panic("implement me")
}
func (ri *RegisterInfluencer) Store(context OipContext) error {
	panic("implement me")
}

func (ri *RegisterInfluencer) Validate(context OipContext) (OipAction, error) {
	v := []string{ri.FloAddress, strconv.FormatInt(int64(ri.Version), 10), strconv.FormatInt(ri.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(ri.FloAddress, ri.Signature, preImage)
	if !sigOk {
		return ri, ErrBadSignature
	}

	return ri, nil
}

func (ei *EditInfluencer) Validate(context OipContext) (OipAction, error) {
	return ei, ErrNotImplemented
}
