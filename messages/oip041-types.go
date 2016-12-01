package messages

type Oip041Wrapper struct {
	Oip041 Oip041 `json:"oip-041"`
}

type Oip041 struct {
	Artifact  Oip041Artifact `json:"artifact"`
	Signature string         `json:"signature"`
}

type Oip041Artifact struct {
	Publisher string        `json:"publisher"`
	Timestamp int           `json:"timestamp"`
	Type      string        `json:"type"`
	Info      Oip041Info    `json:"info"`
	Payment   Oip041Payment `json:"payment"`
}

type Oip041Info struct {
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	Year            int                  `json:"year"`
	ExtraInfo       Oip041MusicExtraInfo `json:"extra-info"`
	ExtraInfoString string
}

type Oip041Payment struct {
	Fiat   string       `json:"fiat"`
	Scale  string       `json:"scale"`
	SugTip []int        `json:"sug_tip"`
	Tokens Oip041Tokens `json:"tokens"`
}

type Oip041MusicExtraInfo struct {
	Artist            string        `json:"artist"`
	Company           string        `json:"company"`
	Composers         []string      `json:"composers"`
	Copyright         string        `json:"copyright"`
	UsageProhibitions string        `json:"usageProhibitions"`
	UsageRights       string        `json:"usageRights"`
	Tags              []string      `json:"tags"`
	Storage           Oip041Storage `json:"storage"`
	Files             []Oip041Files `json:"files"`
}

type Oip041Storage struct {
	Network  string `json:"network"`
	Location string `json:"location"`
}

type Oip041Files struct {
	DisallowBuy   bool          `json:"disallowBuy,omitempty"`
	Dname         string        `json:"dname"`
	Duration      int           `json:"duration,omitempty"`
	Fname         string        `json:"fname"`
	Fsize         int           `json:"fsize"`
	MinPlay       string        `json:"minPlay,omitempty"`
	SugPlay       string        `json:"sugPlay,omitempty"`
	Promo         string        `json:"promo,omitempty"`
	Retail        string        `json:"retail,omitempty"`
	PtpFT         int           `json:"ptpFT,omitempty"`
	PtpDT         int           `json:"ptpDT,omitempty"`
	PtpDA         int           `json:"ptpDA,omitempty"`
	Type          string        `json:"type"`
	TokenlyID     string        `json:"tokenly_ID,omitempty"`
	DissallowPlay bool          `json:"dissallowPlay,omitempty"`
	MinBuy        string        `json:"minBuy,omitempty"`
	SugBuy        string        `json:"sugBuy,omitempty"`
	Storage       Oip041Storage `json:"storage,omitempty"`
}

type Oip041Tokens interface{} // dynamic field
