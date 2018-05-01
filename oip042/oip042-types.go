package oip042

import (
	"github.com/jmoiron/sqlx"
	"github.com/metacoin/flojson"
)

type Oip042 struct {
	Register   *Register   `json:"register,omitempty"`
	Publish    *Publish    `json:"publish,omitempty"`
	Edit       *Edit       `json:"edit,omitempty"`
	Deactivate *Deactivate `json:"deactivate,omitempty"`
	Transfer   *Transfer   `json:"transfer,omitempty"`
}

type Register struct {
	Pub           *RegisterPub           `json:"pub,omitempty"`
	Platform      *RegisterPlatform      `json:"platform,omitempty"`
	Affiliate     *RegisterAffiliate     `json:"affiliate,omitempty"`
	Autominer     *RegisterAutominer     `json:"autominer,omitempty"`
	AutominerPool *RegisterAutominerPool `json:"autominerPool,omitempty"`
	//Signature     string                 `json:"signature"`
}

type Publish struct {
	Artifact *PublishArtifact `json:"artifact,omitempty"`
	//Signature string           `json:"signature"`
}

type Edit struct {
	Pub           *EditPub           `json:"pub,omitempty"`
	Platform      *EditPlatform      `json:"platform,omitempty"`
	Affiliate     *EditAffiliate     `json:"affiliate,omitempty"`
	Autominer     *EditAutominer     `json:"autominer,omitempty"`
	AutominerPool *EditAutominerPool `json:"autominerPool,omitempty"`
	Artifact      *EditArtifact      `json:"artifact,omitempty"`
	//Signature     string             `json:"signature"`
}

type Transfer struct {
	Artifact *TransferArtifact `json:"artifact,omitempty"`
	//Signature string            `json:"signature"`
}

type Deactivate struct {
	Artifact *DeactivateArtifact `json:"artifact,omitempty"`
	//Signature string              `json:"signature"`
}

type OipContext struct {
	Tx          *flojson.TxRawResult
	TxComment   string
	TxId        string
	BlockHeight int64
	Block       *flojson.BlockResult
	DbTx        *sqlx.Tx
	IsEdit      bool
	Reference   string
	signature   string
}
