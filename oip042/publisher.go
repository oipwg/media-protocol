package oip042

import (
	"encoding/json"
	"github.com/oipwg/media-protocol/utility"
	"strconv"
	"strings"
)

type RegisterPub struct {
	Alias      string   `json:"alias"`
	FloAddress string   `json:"floAddress"`
	Timestamp  int64    `json:"timestamp"`
	Authorized []string `json:"authorized"`
	Info       struct {
		Emailmd5           string `json:"emailmd5"`
		AvatarNetwork      string `json:"avatarNetwork"`
		Avatar             string `json:"avatar"`
		HeaderImageNetwork string `json:"headerImageNetwork"`
		HeaderImage        string `json:"headerImage"`
		Bitmessage         string `json:"bitmessage"`
	} `json:"info"`
	Verification struct {
		Imdb        string `json:"imdb"`
		Musicbrainz string `json:"musicbrainz"`
		Twitter     string `json:"twitter"`
		Facebook    string `json:"facebook"`
	} `json:"verification"`
	Signature string `json:"signature"`
}

func (rp RegisterPub) Store(context OipContext) error {
	panic("implement me")
}

type EditPub struct {
	Address   string          `json:"address"`
	Timestamp int64           `json:"timestamp"`
	Patch     json.RawMessage `json:"patch"`
	Signature string          `json:"signature"`
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

func (ep EditPub) Validate(context OipContext) (OipAction, error) {
	return ep, ErrNotImplemented
}
