package oip042

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
)

type RegisterPlatform struct {
	Alias        string            `json:"alias"`
	Timestamp    int64             `json:"timestamp"`
	FloAddress   string            `json:"floAddress"`
	Addresses    []PaymentAddress  `json:"addresses"`
	ShortMW      []string          `json:"shortMW"`
	Verification map[string]string `json:"verification"`
	Version      int32             `json:"version"`
	Info         PlatformInfo      `json:"info"`
}

func (rp RegisterPlatform) Store(context OipContext, dbtx *sqlx.Tx) error {
	panic("implement me")
}

type PlatformInfo struct {
	MinShare int32  `json:"minShare"`
	HttpUrl  string `json:"httpURL"`
}

type EditPlatform struct {
	ArtifactID string          `json:"artifactID"`
	Timestamp  int64           `json:"timestamp"`
	Patch      json.RawMessage `json:"patch"`
}

func (platform *EditPlatform) Store(context OipContext, dbtx *sqlx.Tx) error {
	panic("implement me")
}

func (rp RegisterPlatform) Validate(context OipContext) (OipAction, error) {
	return nil, ErrNotImplemented
}

func (platform *EditPlatform) Validate(context OipContext) (OipAction, error) {
	return nil, ErrNotImplemented
}
