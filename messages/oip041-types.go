package messages

type Oip041Wrapper struct {
	Oip041 Oip041 `json:"oip-041"`
}

type Oip041 struct {
	Artifact  Oip041Artifact `json:"artifact"`
	Edit      Oip041Edit     `json:"edit"`
	Transfer  Oip041Transfer `json:"transferArtifact"`
	Signature string         `json:"signature"`
}

type Oip041ArtifactAPIResult struct {
	Block     int         `json:"block"`
	OIP041    interface{} `json:"oip041"`
	Tags      string      `json:"tags"`
	Timestamp int         `json:"timestamp"`
	Title     string      `json:"title"`
	TxID      string      `json:"txid"`
	Type      string      `json:"type"`
	Year      int         `json:"year"`
	Publisher string      `json:"publisher"`
}

type Oip041Transfer struct {
	Reference string `json:"txid"`
	To        string `json:"to"`
	From      string `json:"from"`
	Timestamp int64  `json:"timestamp"`
}

type Oip041Edit struct {
	Add       map[string]string `json:"add"`
	Edit      map[string]string `json:"edit"`
	Remove    []string          `json:"remove"`
	Timestamp int               `json:"timestamp"`
	TxID      string            `json:"txid"`
}

type Oip041Artifact struct {
	Publisher string        `json:"publisher"`
	Timestamp int           `json:"timestamp"`
	Type      string        `json:"type"`
	Info      Oip041Info    `json:"info"`
	Storage   Oip041Storage `json:"storage"`
	Files     []Oip041Files `json:"files"`
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
	Artist            string   `json:"artist"`
	Company           string   `json:"company"`
	Composers         []string `json:"composers"`
	Copyright         string   `json:"copyright"`
	UsageProhibitions string   `json:"usageProhibitions"`
	UsageRights       string   `json:"usageRights"`
	Tags              []string `json:"tags"`
}

type Oip041Storage struct {
	Network  string `json:"network,omitempty"`
	Location string `json:"location,omitempty"`
}

type Oip041Files struct {
	DisallowBuy  int           `json:"disallowBuy"`
	Dname        string        `json:"dname"`
	Duration     float64       `json:"duration,omitempty"`
	Fname        string        `json:"fname"`
	Fsize        int           `json:"fsize"`
	MinPlay      float64       `json:"minPlay"`
	SugPlay      float64       `json:"sugPlay"`
	Promo        float64       `json:"promo"`
	Retail       float64       `json:"retail"`
	PtpFT        int           `json:"ptpFT,omitempty"`
	PtpDT        int           `json:"ptpDT,omitempty"`
	PtpDA        int           `json:"ptpDA,omitempty"`
	Type         string        `json:"type"`
	TokenlyID    string        `json:"tokenly_ID,omitempty"`
	DisallowPlay int           `json:"disallowPlay"`
	MinBuy       float64       `json:"minBuy"`
	SugBuy       float64       `json:"sugBuy"`
	Storage      Oip041Storage `json:"storage"`
}

type Oip041Tokens map[string]string
