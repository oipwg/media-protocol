package oip042

import (
	"github.com/jmoiron/sqlx"
)

type IOip042 interface {
	IOip042()
}

type OipAction interface {
	Validate(context OipContext) (OipAction, error)
	Store(context OipContext, dbtx *sqlx.Tx) error
}
