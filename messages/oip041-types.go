package messages

import (
	"encoding/json"
)

type Oip041Wrapper struct {
	Oip041 Oip041 `json:"oip-041"`
}

type Oip041 struct {
	Artifact   Oip041Artifact   `json:"artifact"`
	Edit       Oip041Edit       `json:"editArtifact"`
	Deactivate Oip041Deactivate `json:"deactivateArtifact"`
	Transfer   Oip041Transfer   `json:"transferArtifact"`
	Signature  string           `json:"signature"`
	artSize    int
}

type Oip041ArtifactAPIResult struct {
	Block         int     `json:"block"`
	OIP041        Oip041  `json:"oip-041"`
	Tags          string  `json:"tags"`
	Timestamp     int64   `json:"timestamp"`
	Title         string  `json:"title"`
	TxID          string  `json:"txid"`
	Type          string  `json:"type"`
	Year          int     `json:"year"`
	Publisher     string  `json:"publisher"`
	PublisherName string  `json:"publisherName"`
	ArtCost       float64 `json:"artCost"`
	ArtSize       int     `json:"artSize"`
	PubFeeUSD     float64 `json:"pubFeeUSD"`
	NSFW          bool    `json:"nsfw"`
}

type Oip041Transfer struct {
	Reference string `json:"txid"`
	To        string `json:"to"`
	From      string `json:"from"`
	Timestamp int64  `json:"timestamp"`
}

type Oip041Deactivate struct {
	Reference string `json:"txid"`
	Timestamp int64  `json:"timestamp"`
}

type Oip041Edit struct {
	Patch     json.RawMessage `json:"patch"`
	Timestamp int64           `json:"timestamp"`
	TxID      string          `json:"txid"`
}

type Oip041Artifact struct {
	Publisher string        `json:"publisher"`
	Timestamp int64         `json:"timestamp"`
	Type      string        `json:"type"`
	Info      Oip041Info    `json:"info"`
	Storage   Oip041Storage `json:"storage"`
	Payment   Oip041Payment `json:"payment"`
}

type Oip041Info struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Tags            string          `json:"tags"`
	Year            int             `json:"year"`
	NSFW            bool            `json:"nsfw"`
	ExtraInfo       json.RawMessage `json:"extraInfo"`
	ExtraInfoString string
}

type Oip041Payment struct {
	Fiat        string          `json:"fiat"`
	Scale       string          `json:"scale"`
	SugTip      []int           `json:"sugTip"`
	Tokens      Oip041Tokens    `json:"tokens"`
	Addresses   []Oip041Address `json:"addresses"`
	Retailer    int             `json:"retailer"`
	Promoter    int             `json:"promoter"`
	MaxDiscount float64         `json:"maxdisc"`
}

type Oip041MusicExtraInfo struct {
	Artist            string   `json:"artist"`
	Company           string   `json:"company"`
	Composers         []string `json:"composers"`
	Copyright         string   `json:"copyright"`
	UsageProhibitions string   `json:"usageProhibitions"`
	UsageRights       string   `json:"usageRights"`
	Genre             string   `json:"genre"`
	Tags              []string `json:"tags"`
	ISRC              string   `json:"ISRC"`
}

type Oip041Storage struct {
	Network  string        `json:"network,omitempty"`
	Location string        `json:"location,omitempty"`
	Files    []Oip041Files `json:"files"`
}

type Oip041Files struct {
	DisallowBuy  bool    `json:"disBuy,omitempty"`
	Dname        string  `json:"dname,omitempty"`
	Duration     float64 `json:"duration,omitempty"`
	Fname        string  `json:"fname,omitempty"`
	Fsize        int     `json:"fsize,omitempty"`
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
	// ToDo: Add per file granularity back, requires custom json marshalling to omit
	// Storage Oip041Storage `json:"storage"`
}

type Oip041Address struct {
	Token   string `json:"token"`
	Address string `json:"address"`
}

type Oip041Tokens map[string]int
