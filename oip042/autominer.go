package oip042

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/oipwg/media-protocol/utility"
	"strconv"
	"strings"
)

type RegisterAutominer struct {
	Alias      string           `json:"alias"`
	Timestamp  int64            `json:"timestamp"`
	FloAddress string           `json:"floAddress"`
	Addresses  []PaymentAddress `json:"addresses"`
	ShortMW    []string         `json:"shortMW"`
	Version    int32            `json:"version"`
	Info       AutominerInfo    `json:"info"`
}

func (ram RegisterAutominer) Store(context OipContext, dbtx *sqlx.Tx) error {
	panic("implement me")
}

type AutominerInfo struct {
	MinShare int32  `json:"minShare"`
	HttpUrl  string `json:"httpURL"`
}

type EditAutominer struct {
	ArtifactID string          `json:"artifactID"`
	Timestamp  int64           `json:"timestamp"`
	Patch      json.RawMessage `json:"patch"`
}

func (autominer *EditAutominer) Store(context OipContext, dbtx *sqlx.Tx) error {
	panic("implement me")
}

type RegisterAutominerPool struct {
	Alias        string            `json:"alias"`
	Timestamp    int64             `json:"timestamp"`
	FloAddress   string            `json:"floAddress"`
	Addresses    []PaymentAddress  `json:"addresses"`
	ShortMW      []string          `json:"shortMW"`
	Verification map[string]string `json:"verification"`
	Version      int32             `json:"version"`
	Info         AutominerInfo     `json:"info"`
}

func (ramp RegisterAutominerPool) Store(context OipContext, dbtx *sqlx.Tx) error {
	panic("implement me")
}

type AutominerPoolInfo struct {
	PoolShare    int32  `json:"poolShare"`
	TargetMargin int32  `json:"targetMargin"`
	HttpUrl      string `json:"httpURL"`
}

type EditAutominerPool struct {
	ArtifactID string          `json:"artifactID"`
	Timestamp  int64           `json:"timestamp"`
	Patch      json.RawMessage `json:"patch"`
}

func (autominerPool *EditAutominerPool) Store(context OipContext, dbtx *sqlx.Tx) error {
	panic("implement me")
}

func (ram RegisterAutominer) Validate(context OipContext) (OipAction, error) {
	v := []string{ram.FloAddress, strconv.FormatInt(int64(ram.Version), 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(ram.FloAddress, context.signature, preImage)
	if !sigOk {
		return ram, ErrBadSignature
	}

	return ram, nil
}

func (ramp RegisterAutominerPool) Validate(context OipContext) (OipAction, error) {
	return ramp, ErrNotImplemented
}

func (autominer *EditAutominer) Validate(context OipContext) (OipAction, error) {
	return autominer, ErrNotImplemented
}
func (autominerPool *EditAutominerPool) Validate(context OipContext) (OipAction, error) {
	return autominerPool, ErrNotImplemented
}
