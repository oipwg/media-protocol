package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/metacoin/flojson"
	"github.com/oipwg/media-protocol/utility"
)

const PROMOTER_ROOT_KEY = "alexandria-promoter"

type Promoter_SocialMedias []Promoter_SocialMedia

type Promoter_SocialMedia struct {
	Sequence int    `json:"sequence"`
	Network  string `json:"network"`
	Username string `json:"username"`
}
type Promoter struct {
	FLOAddress  string                `json:"FLOaddress"`
	BTCAddress  string                `json:"BTCaddress"`
	Version     int64                 `json:"version"`
	SocialMedia Promoter_SocialMedias `json:"social-media"`
}
type AlexandriaPromoter struct {
	Promoter  Promoter `json:"alexandria-promoter"`
	Signature string   `json:"signature"`
}

// check for duplicate Sequence numbers in Promoter_SocialMedia slice
// also, remember to check for duplicates before sorting
// TODO: make this better or implement a cap on social media networks
func (promoter Promoter) HasSocialMediaSequenceDuplicate() bool {
	var seen []int
	var numSeen int

	for _, i := range promoter.SocialMedia {
		// return true if there is a negative sequence number
		if i.Sequence < 0 {
			return true
		}

		// find duplicates
		numSeen = 0
		seen = append(seen, i.Sequence)
		for _, j := range seen {
			if i.Sequence == j {
				numSeen++
			}
		}

		// if we find a duplicate we're done
		if numSeen > 1 {
			return true
		}
	}

	// no duplicates - return false
	return false
}

func (slice Promoter_SocialMedias) Len() int {
	return len(slice)
}

func (slice Promoter_SocialMedias) Less(i, j int) bool {
	return slice[i].Sequence < slice[j].Sequence
}

func (slice Promoter_SocialMedias) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func StorePromoter(pr AlexandriaPromoter, dbtx *sql.Tx, txid string, block *flojson.BlockResult) error {
	// store in database
	stmtStr := `insert into promoter (txid, block, blockTime, active, version,` +
		` floAddress, btcAddress, signature) values (?, ?, ?, 1, ?, ?, ?, ?)`

	stmt, err := dbtx.Prepare(stmtStr)
	if err != nil {
		fmt.Printf("ERROR in StorePromoter: promoter dbtx didn't prepare correctly: %v", err)
		return err
	}

	res, stmterr := stmt.Exec(txid, block.Height, block.Time, pr.Promoter.Version, pr.Promoter.FLOAddress, pr.Promoter.BTCAddress, pr.Signature)
	if stmterr != nil {
		fmt.Printf("ERROR in StorePromoter: promoter dbtx didn't execute correctly: %v", err)
		return stmterr
	}

	// optional social media fields
	if len(pr.Promoter.SocialMedia) > 0 {
		lastId, err := res.LastInsertId()
		if err != nil {
			fmt.Printf("ERROR in StorePromoter: couldn't get LastInsertId() from res: %v\n", err)
			return err
		}
		err = StorePromoterSocialMedia(pr.Promoter.SocialMedia, dbtx, txid, block, lastId)
		if err != nil {
			return err
		}
	}

	stmt.Close()
	return nil
}

func StorePromoterSocialMedia(pr []Promoter_SocialMedia, dbtx *sql.Tx, txid string, block *flojson.BlockResult, id int64) error {
	for _, j := range pr {
		stmtStr := `insert into promoter_socialmedia (promoter_uid, network, username) values (?, ?, ?)`
		stmt, err := dbtx.Prepare(stmtStr)
		if err != nil {
			fmt.Printf("ERROR in StorePromoterSocialMedia: promoter_socialmedia dbtx didn't prepare correctly: %v\n", err)
			return err
		}

		_, stmterr := stmt.Exec(id, j.Network, j.Username)
		if stmterr != nil {
			fmt.Printf("ERROR in StorePromoterSocialMedia: promoter_socialmedia dbtx didn't execute correctly: %v\n", err)
			return stmterr
		}
		stmt.Close()
	}
	return nil
}

func VerifyPromoter(b []byte, block int) (AlexandriaPromoter, error) {

	//fmt.Printf("starting VerifyPromoter routine...\n")

	//s := string(b[:len(b)])
	//fmt.Printf("s: %+v\n", s)

	var pr AlexandriaPromoter

	if block < 2205000 {
		return pr, ErrTooEarly
	}

	if !utility.IsJSON(string(b)) {
		return pr, ErrNotJSON
	}

	err := json.Unmarshal(b, &pr)
	if err != nil {
		return pr, err
	}

	//fmt.Printf("pr: %+v\n", pr)

	// verify signature was created by this address
	// signature pre-image for promoter is <btcaddress>-<version>
	preImage := pr.Promoter.BTCAddress + "-" + strconv.FormatInt(pr.Promoter.Version, 10)

	// if there are optional fields, they are included in order of sequence
	// example: <btcaddress>-<version>-<network>-<username>-<network>-...
	if len(pr.Promoter.SocialMedia) > 0 && !pr.Promoter.HasSocialMediaSequenceDuplicate() {
		sort.Sort(pr.Promoter.SocialMedia)
		for _, i := range pr.Promoter.SocialMedia {
			preImage += "-" + i.Network + "-" + i.Username
		}
	}

	//fmt.Printf("\n###### pre-image: %v", preImage)
	////fmt.Printf("\n###### signature: %v\n", pr.Signature)
	//fmt.Printf("\n\n\n")
	sigOK, _ := utility.CheckSignature(pr.Promoter.FLOAddress, pr.Signature, preImage)
	if sigOK == false {
		return pr, ErrBadSignature
	}

	// fmt.Println(" -- VERIFIED --")
	return pr, nil
}
