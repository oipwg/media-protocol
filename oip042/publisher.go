package oip042

import (
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/oipwg/media-protocol/utility"
	"strconv"
	"strings"
)

type RegisterPub struct {
	Alias        string                 `json:"alias,omitempty"`
	FloAddress   string                 `json:"floAddress,omitempty"`
	Timestamp    int64                  `json:"timestamp,omitempty"`
	Authorized   []string               `json:"authorized,omitempty"`
	Info         *PublisherInfo         `json:"info,omitempty"`
	Verification *PublisherVerification `json:"verification,omitempty"`
	Signature    string                 `json:"signature,omitempty"`
}

type PublisherVerification struct {
	Imdb        string `json:"imdb,omitempty"`
	Musicbrainz string `json:"musicbrainz,omitempty"`
	Twitter     string `json:"twitter,omitempty"`
	Facebook    string `json:"facebook,omitempty"`
}

type PublisherInfo struct {
	Emailmd5           string `json:"emailmd5,omitempty"`
	AvatarNetwork      string `json:"avatarNetwork,omitempty"`
	Avatar             string `json:"avatar,omitempty"`
	HeaderImageNetwork string `json:"headerImageNetwork,omitempty"`
	HeaderImage        string `json:"headerImage,omitempty"`
	Bitmessage         string `json:"bitmessage,omitempty"`
}

func (rp RegisterPub) Store(context OipContext) error {
	j, err := json.Marshal(rp)
	if err != nil {
		return err
	}

	cv := map[string]interface{}{
		"json":        j,
		"alias":       rp.Alias,
		"floAddress":  rp.FloAddress,
		"active":      1,
		"invalidated": 0,
		"validated":   0,
	}

	var q squirrel.Sqlizer
	if context.IsEdit {
		cv["txid"] = context.Reference
		q = squirrel.Update("pub").SetMap(cv).Where(squirrel.Eq{"txid": context.Reference})
	} else {
		// these values are only set on publish
		cv["active"] = 1
		cv["unixtime"] = rp.Timestamp
		cv["timestamp"] = rp.Timestamp
		cv["txid"] = context.TxId
		cv["block"] = context.BlockHeight
		q = squirrel.Insert("pub").SetMap(cv)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = context.DbTx.Exec(sql, args...)
	if err != nil {
		return err
	}

	return nil
}

type EditPub struct {
	Address   string          `json:"address"`
	Timestamp int64           `json:"timestamp"`
	Patch     json.RawMessage `json:"patch"`
	Signature string          `json:"signature"`
}

type DeactivatePub struct {
	FloAddress string `json:"address"`
	Timestamp  int64  `json:"timestamp"`
	Signature  string `json:"signature"`
}

func (ep EditPub) Store(context OipContext) error {
	panic("implement me")
}

func (rp RegisterPub) Validate(context OipContext) (OipAction, error) {
	v := []string{rp.Alias, rp.FloAddress, strconv.FormatInt(rp.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(rp.FloAddress, rp.Signature, preImage)
	if !sigOk {
		return rp, ErrBadSignature
	}

	return rp, nil
}

func (dp DeactivatePub) Validate(context OipContext) (OipAction, error) {
	v := []string{"deactivate", dp.FloAddress, strconv.FormatInt(dp.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(dp.FloAddress, dp.Signature, preImage)
	if !sigOk {
		return dp, ErrBadSignature
	}

	return dp, nil
}

func (dp DeactivatePub) Store(context OipContext) error {
	q := squirrel.Update("pub").
		Set("active", 0).
		Where(squirrel.Eq{"floAddress": dp.FloAddress})

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = context.DbTx.Exec(sql, args...)
	if err != err {
		return err
	}

	return nil
}

func (ep EditPub) Validate(context OipContext) (OipAction, error) {
	return ep, ErrNotImplemented
}

const createPublisherTable = `
CREATE TABLE IF NOT EXISTS pub
(
    alias text DEFAULT '' NOT NULL,
    floAddress text NOT NULL,
    unixtime int NOT NULL,
    active boolean NOT NULL,
    invalidated boolean NOT NULL,
    validated boolean NOT NULL,
    txid text NOT NULL,
    block int NOT NULL,
    uid int PRIMARY KEY,
    timestamp int NOT NULL,
    json text NOT NULL
);
CREATE INDEX IF NOT EXISTS pub_block_index ON pub (block);
CREATE INDEX IF NOT EXISTS pub_unixtime_index ON pub (unixtime);
CREATE INDEX IF NOT EXISTS pub_active_index ON pub (active);
CREATE INDEX IF NOT EXISTS pub_floAddress_index ON pub (floAddress);
`
