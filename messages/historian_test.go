package messages

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"reflect"
	"testing"
)

var (
	DBH *sql.DB
)

func TestVerifyHistorianMessage(t *testing.T) {
	fmt.Printf("***TestVerifyHistorianMessage***")

	createTestDB(t)
	dbtx, err := DBH.Begin()
	if err != nil {
		t.Fatal("Couldn't initialize db tx")
	}

	// signed FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD
	// valid
	s := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	hm := HistorianMessage{
		1, "pool.alexandria.io", 0.0001360085, 0, 3.163064456533333e+08, math.Inf(-1), 5e-06, 0.00217, 0, "IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA=",
	}
	// bad version
	s1 := "alexandria-historian-v002:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// bad pool
	s2 := "alexandria-historian-v001:notpool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="
	// wrong length
	s3 := "alexandria-historian-v001:"
	// no signature
	s4 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217"
	hm4 := HistorianMessage{
		1, "pool.alexandria.io", 0.0001360085, 0, 3.163064456533333e+08, math.Inf(-1), 5e-06, 0.00217, 0, "",
	}
	// trailing :
	s5 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:"
	// invalid signature
	s6 := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwAa"
	// signed FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU
	// valid
	s7 := "alexandria-historian-v001:pool.alexandria.io:0.000104048500:223208386.28518352:2214713879:0.00000429:0.00308:ICyn+Wh4OxKF89+O9u0wkQULeyvJ6CDurGiZACCkNtk8Rl+QpejBmPWKYiuyt6PM5+MrUs/gDcACWjKFTSoYrxA="
	hm7 := HistorianMessage{
		1, "pool.alexandria.io", 0.0001040485, 0, 2.2320838628518352e+08, 2.214713879e+09, 4.29e-06, 0.00308, 0, "ICyn+Wh4OxKF89+O9u0wkQULeyvJ6CDurGiZACCkNtk8Rl+QpejBmPWKYiuyt6PM5+MrUs/gDcACWjKFTSoYrxA=",
	}
	// invalid signature
	s8 := "oip-historian-1:FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU:0.000111054110:186009592.24127597:13858880968:0.00001983:0.04655:signature"
	// signed with FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU
	// valid
	s9 := "oip-historian-1:FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU:0.000111054110:186009592.24127597:13858880968:0.00001983:0.04655:HxsHFYL+ROlGABTxTXFVtV/g+bI8M/vimB77xvs2iEdMEiKKZ819UQeos6uBa5XNuBZspccVI6PCfB6OeJoIB/8="
	hm9 := HistorianMessage{
		1, "alexandria.io", 0.00011105411, 0, 1.8600959224127597e+08, 1.3858880968e+10, 1.983e-05, 0.04655, 0, "HxsHFYL+ROlGABTxTXFVtV/g+bI8M/vimB77xvs2iEdMEiKKZ819UQeos6uBa5XNuBZspccVI6PCfB6OeJoIB/8=",
	}
	// signed with FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU
	// valid
	s10 := "oip-historian-2:FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU:0.001084668945:nr:415415641.6450302:3045555424:0.00002092:0.05743:IOO0GGausgoY2d38vQFr9anU1x1k7MkMZirBD9b3t+VOVjTU5tpGoYSqW8+Yb1+o/UqfiSYDZ0PaNJGIfE85+bw="
	hm10 := HistorianMessage{
		1, "alexandria.io", 0.001084668945, math.Inf(-1), 415415641.6450302, 3045555424, 0.00002092, 0.05743, 0, "IOO0GGausgoY2d38vQFr9anU1x1k7MkMZirBD9b3t+VOVjTU5tpGoYSqW8+Yb1+o/UqfiSYDZ0PaNJGIfE85+bw=",
	}

	s11 := "oip-historian-3:FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU:0.000002599888:0.00013285258034166893:4748423487.878089:95118575981:0.00000653:0.10900:182.029:II2qMMcDjY69FVuOiqMDr0rg0wRAmIV6HIxOqjdSIVMrepal+FXT1uboKn9CvF4ar8SIz8BIs9Gml8iRk2d/Wts="
	hm11 := HistorianMessage{
		3, "alexandria.io", 0.000002599888, 0.00013285258034166893, 4748423487.878089, 95118575981, 0.00000653, 0.10900, 182.029, "II2qMMcDjY69FVuOiqMDr0rg0wRAmIV6HIxOqjdSIVMrepal+FXT1uboKn9CvF4ar8SIz8BIs9Gml8iRk2d/Wts=",
	}

	nilHM := HistorianMessage{}

	cases := []struct {
		in    string
		out   HistorianMessage
		block int
		err   error
	}{
		{s, hm, 1750000, nil},                                  // valid
		{s, nilHM, 1974560, ErrBadSignature},                   // wrong address
		{s1, nilHM, 1974560, ErrWrongPrefix},                   // bad version
		{s2, nilHM, 1974560, ErrHistorianMessagePoolUntrusted}, // bad pool
		{s3, nilHM, 1974560, ErrHistorianMessageInvalid},       // wrong length
		{s4, nilHM, 1974560, ErrBadSignature},                  // no signature
		{s4, hm4, 1974559, nil},                                // no signature, but unenforced
		{s5, nilHM, 1974560, ErrBadSignature},                  // trailing :
		{s6, nilHM, 1974560, ErrBadSignature},                  // invalid signature
		{s7, hm7, 1974560, nil},                                // valid
		{s8, nilHM, 2000000, ErrBadSignature},                  //invalid signature
		{s9, hm9, 2000000, nil},                                // valid
		{s10, hm10, 2000000, nil},                              // valid
		{s11, hm11, 2480672, nil},                              // valid
	}

	for i, c := range cases {
		got, err := VerifyHistorianMessage([]byte(c.in), c.block, dbtx)
		if err != c.err {
			t.Errorf("VerifyHistorianMessage(#%d) | err == %v, want %v", i, err, c.err)
		}
		if err == nil && !reflect.DeepEqual(got, c.out) {
			t.Errorf("VerifyMediaMultipartSingle(#%d) | got == %v, want %v", i, got, c.out)
		}
	}
}

var (
	// to prevent the compiler deleting the benchmark
	hmTestErr error
	hmTestHM  HistorianMessage
)

// This benchmark was really just for curiosity sake, maybe later it will
// actually be adapted to serve a purpose
func BenchmarkVerifyHistorianMessage(b *testing.B) {
	// signed FL4Ty99iBsGu3aPrGx6rwUtWwyNvUjb7ZD
	// valid
	s := "alexandria-historian-v001:pool.alexandria.io:0.000136008500:316306445.6533333:nr:0.00000500:0.00217:IN9OrF1Kpd5S0x36nXWI0lFjhnS1Z9I9k7cxWJrFUlsfcgwJytZ+GlKP1/tHCijAdGAX6LnOgOtcvI/vMQgVcwA="

	for n := 0; n < b.N; n++ {
		// ToDo: benchmark oip-historian with database lookup
		hmTestHM, hmTestErr = VerifyHistorianMessage([]byte(s), 1750000, nil)
	}
}

func createTestDB(t *testing.T) {
	var err error
	DBH, err = sql.Open("sqlite3", ":memory:")

	if err != nil || DBH == nil {
		t.Fatal("Database couldn't be opened.")
		return
	}

	createTable := `CREATE TABLE autominer_pool (
                        uid integer not null primary key AUTOINCREMENT,
                        txid TEXT not null,
                        block int not null,
                        blockTime int not null,
                        active int not null,
                        version int not null,
                        floAddress TEXT not null,
                        webURL TEXT not null,
                        targetMargin FLOAT not null,
                        poolShare FLOAT not null,
                        poolName TEXT,
                        signature TEXT not null,
                        invalidated int default 0
                );`
	insertData := `INSERT INTO "autominer_pool" VALUES(1,'a73b37b07b96f4f5d8d23a941a0eedc2747d79f1d5e7a6cab1a1acdd00d7a4f1',2243836,1500521093,1,1,'FLmic78oU6eqXsTAaHGGdrFyY7FznjHfPU','alexandria.io',20.0,0.0,'','IGridS+M9cNrwFpitJmESfSQ33NXGCTVhoB4P9H3u+gODwnWIHG4711cxSohDx9/d700RDfYvPiOrqW3zh3Y7XI=',0);`

	tx, err := DBH.Begin()
	if err != nil {
		t.Fatal("Couldn't initialize db tx")
	}

	tx.Exec(createTable)
	tx.Exec(insertData)
	tx.Commit()
}
