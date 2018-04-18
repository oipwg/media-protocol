package oip042

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
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
	Patch      json.RawMessage `json:"patch"`
}

func (ea EditArtifact) Store(context OipContext) error {
	panic("implement me")
}

type TransferArtifact struct {
	ArtifactID     string `json:"artifactID"`
	ToFloAddress   string `json:"toFloAddress"`
	FromFloAddress string `json:"fromFloAddress"`
	Timestamp      int64  `json:"timestamp"`
}

func (ta TransferArtifact) Store(context OipContext) error {
	panic("implement me")
}

type DeactivateArtifact struct {
	ArtifactID string `json:"artifactID"`
	Timestamp  int64  `json:"timestamp"`
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
	v := []string{pa.Storage.Location, pa.FloAddress, strconv.FormatInt(pa.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(pa.FloAddress, pa.Signature, preImage)
	if !sigOk {
		return nil, ErrBadSignature
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
	return nil, ErrNotImplemented

	v := []string{ea.ArtifactID, "ToDo", strconv.FormatInt(ea.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature("ToDo", context.signature, preImage)
	if !sigOk {
		return nil, ErrBadSignature
	}

	return ea, nil
}

func (ta TransferArtifact) Validate(context OipContext) (OipAction, error) {
	v := []string{ta.ArtifactID, ta.ToFloAddress, ta.FromFloAddress, strconv.FormatInt(ta.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(ta.FromFloAddress, context.signature, preImage)
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
	if err := row.Scan(&publisher); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	v := []string{da.ArtifactID, publisher, strconv.FormatInt(da.Timestamp, 10)}
	preImage := strings.Join(v, "-")
	sigOk, _ := utility.CheckSignature(publisher, context.signature, preImage)
	if !sigOk {
		return nil, ErrBadSignature
	}

	return da, nil
}
