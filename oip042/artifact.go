package oip042

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bitspill/json-patch"
	"github.com/oipwg/media-protocol/utility"
	"strconv"
	"strings"
)

type ArtifactInfo struct {
	Title       string `json:"title,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Description string `json:"description,omitempty"`
	Year        int    `json:"year,omitempty"`
	NSFW        bool   `json:"nsfw,omitempty"`
}

type ArtifactPayment struct {
	Fiat        string           `json:"fiat,omitempty"`
	Scale       string           `json:"scale,omitempty"`
	SugTip      []int            `json:"sugTip,omitempty"`
	Tokens      *PaymentTokens   `json:"tokens,omitempty"`
	Addresses   []PaymentAddress `json:"addresses,omitempty"`
	Platform    int              `json:"platform,omitempty"`
	Affiliate   int              `json:"affiliate,omitempty"`
	MaxDiscount float64          `json:"maxdisc,omitempty"`
	ShortMW     []string         `json:"shortMW,omitempty"`
}

type ArtifactStorage struct {
	Network  string          `json:"network,omitempty"`
	Location string          `json:"location,omitempty"`
	Files    []ArtifactFiles `json:"files,omitempty"`
}

type ArtifactFiles struct {
	Software     string  `json:"software,omitempty"`
	DisallowBuy  bool    `json:"disBuy,omitempty"`
	Dname        string  `json:"dname,omitempty"`
	Duration     float64 `json:"duration,omitempty"`
	Fname        string  `json:"fname,omitempty"`
	Fsize        int64   `json:"fsize,omitempty"`
	MinPlay      float64 `json:"minPlay,omitempty"`
	SugPlay      float64 `json:"sugPlay,omitempty"`
	Promo        float64 `json:"promo,omitempty"`
	Retail       float64 `json:"retail,omitempty"`
	PtpFT        int     `json:"ptpFT,omitempty"`
	PtpDT        int     `json:"ptpDT,omitempty"`
	PtpDA        int     `json:"ptpDA,omitempty"`
	Type         string  `json:"type,omitempty"`
	TokenlyID    string  `json:"tokenlyID,omitempty"`
	DisallowPlay bool    `json:"disPlay,omitempty"`
	MinBuy       float64 `json:"minBuy,omitempty"`
	SugBuy       float64 `json:"sugBuy,omitempty"`
	SubType      string  `json:"subtype,omitempty"`
	CType        string  `json:"cType,omitempty"`
	FNotes       string  `json:"fNotes,omitempty"`
}

type PaymentAddress map[string]string

type PaymentTokens map[string]int

type PublishArtifact struct {
	FloAddress string           `json:"floAddress,omitempty"`
	Timestamp  int64            `json:"timestamp,omitempty"`
	Type       string           `json:"type,omitempty"`
	SubType    string           `json:"subtype,omitempty"`
	Info       *ArtifactInfo    `json:"info,omitempty"`
	Details    json.RawMessage  `json:"details,omitempty"`
	Storage    *ArtifactStorage `json:"storage,omitempty"`
	Payment    *ArtifactPayment `json:"payment,omitempty"`
	Signature  string           `json:"signature,omitempty"`
}

func (pa PublishArtifact) Store(context OipContext) error {
	// ToDo store generic publishes without indexing details
	fmt.Println("Attempted to store unknown PublishArtifact type")
	fmt.Println("Disregarding for now.")
	return ErrNotImplemented
}

type EditArtifact struct {
	ArtifactID string          `json:"artifactID"`
	Timestamp  int64           `json:"timestamp"`
	RawPatch   json.RawMessage `json:"patch"`
	Patch      jsonpatch.Patch
	Signature  string `json:"signature"`
}

func (ea EditArtifact) Store(context OipContext) error {
	sq := squirrel.Select("uid", "publisher", "json", "title", "type", "subType", "nsfw", "hasDetails").
		From("artifact").
		Where(squirrel.Eq{"txid": ea.ArtifactID})

	query, args, err := sq.ToSql()
	if err != nil {
		return err
	}

	type data struct {
		Uid        int64  `db:"uid"`
		Publisher  string `db:"publisher"`
		Json       []byte `db:"json"`
		Title      string `db:"title"`
		Type       string `db:"type"`
		SubType    string `db:"subType"`
		Nsfw       bool   `db:"nsfw"`
		HasDetails bool   `db:"hasDetails"`
	}
	var d data
	row := context.DbTx.QueryRowx(query, args...)
	if err := row.StructScan(&d); err != nil {
		return err
	}

	pp, err := ea.Patch.Apply(d.Json)
	if err != nil {
		return err
	}

	var pa PublishArtifact
	err = json.Unmarshal(pp, &pa)
	if err != nil {
		return err
	}

	context.IsEdit = true
	context.Reference = ea.ArtifactID
	context.ArtifactId = d.Uid
	ep, err := pa.Validate(context)
	if err != nil {
		return err
	}

	err = ep.Store(context)
	if err != nil {
		return err
	}

	return nil
}

type TransferArtifact struct {
	ArtifactID     string `json:"artifactID"`
	ToFloAddress   string `json:"toFloAddress"`
	FromFloAddress string `json:"fromFloAddress"`
	Timestamp      int64  `json:"timestamp"`
	Signature      string `json:"signature"`
}

func (ta TransferArtifact) Store(context OipContext) error {
	panic("implement me")
}

type DeactivateArtifact struct {
	ArtifactID string `json:"artifactID"`
	Timestamp  int64  `json:"timestamp"`
	Signature  string `json:"signature"`
}

func (da DeactivateArtifact) Store(context OipContext) error {
	q := squirrel.Update("artifact").
		Set("invalidated", 1).
		Where(squirrel.Eq{"txid": da.ArtifactID})

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = context.DbTx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

var ErrDescriptionMissing = errors.New("artifact missing description")
var ErrTypeMissing = errors.New("artifact missing type")

func (pa PublishArtifact) Validate(context OipContext) (OipAction, error) {
	if !context.IsEdit {
		// only validate signatures if it's the first go 'round, edits may have changed signed values
		var loc string
		if pa.Storage != nil {
			loc = pa.Storage.Location
		} else {
			loc = ""
		}
		v := []string{loc, pa.FloAddress, strconv.FormatInt(pa.Timestamp, 10)}
		preImage := strings.Join(v, "-")
		sigOk, _ := utility.CheckSignature(pa.FloAddress, pa.Signature, preImage)
		if !sigOk {
			return nil, ErrBadSignature
		}
	}
	if len(strings.TrimSpace(pa.Info.Description)) == 0 {
		return nil, ErrDescriptionMissing
	}
	if len(strings.TrimSpace(pa.Type)) == 0 {
		return nil, ErrTypeMissing
	}
	if pa.Storage != nil && strings.ToLower(pa.Storage.Network) != "ipfs" {
		return nil, errors.New("artifact: only IPFS network is supported")
	}
	if pa.Timestamp <= 0 {
		return nil, errors.New("artifact: invalid timestamp")
	}

	if pa.Type == "research" && pa.SubType == "tomogram" {
		return PublishTomogram{PublishArtifact: pa}.Validate(context)
	}
	if pa.Type == "property" {
		if pa.SubType == "party" {
			return PublishPropertyParty{PublishArtifact: pa}.Validate(context)
		}
		if pa.SubType == "tenure" {
			return PublishPropertyTenure{PublishArtifact: pa}.Validate(context)
		}
		if pa.SubType == "spatialUnit" {
			return PublishPropertySpatialUnit{PublishArtifact: pa}.Validate(context)
		}
	}

	return pa, nil
}

func (ea EditArtifact) Validate(context OipContext) (OipAction, error) {
	q := squirrel.Select("publisher").
		From("artifact").
		Where(squirrel.Eq{"txid": ea.ArtifactID})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var publisher string
	row := context.DbTx.QueryRow(query, args...)
	if err := row.Scan(&publisher); err != nil {
		return nil, err
	}

	v := []string{ea.ArtifactID, publisher, strconv.FormatInt(ea.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(publisher, ea.Signature, preImage)
	if !sigOk {
		return nil, ErrBadSignature
	}

	ea.Patch, err = utility.UnSquashPatch(ea.RawPatch)
	if err != nil {
		return ea, err
	}

	return ea, nil
}

func (ta TransferArtifact) Validate(context OipContext) (OipAction, error) {
	v := []string{ta.ArtifactID, ta.ToFloAddress, ta.FromFloAddress, strconv.FormatInt(ta.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(ta.FromFloAddress, ta.Signature, preImage)
	if !sigOk {
		return nil, ErrBadSignature
	}

	return ta, nil
}

func (da DeactivateArtifact) Validate(context OipContext) (OipAction, error) {
	q := squirrel.Select("publisher").
		From("artifact").
		Where(squirrel.Eq{"txid": da.ArtifactID})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	var publisher string
	row := context.DbTx.QueryRow(query, args...)
	if err := row.Scan(&publisher); err != nil {
		return nil, err
	}

	v := []string{da.ArtifactID, publisher, strconv.FormatInt(da.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	fmt.Println(preImage)
	sigOk, _ := utility.CheckSignature(publisher, da.Signature, preImage)
	if !sigOk {
		return nil, ErrBadSignature
	}

	return da, nil
}
