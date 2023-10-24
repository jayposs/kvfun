// Program client1.go is an example client program that sends requests to a running server.
// It demonstrates all request types supported by the kvf pkg.

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"kvfun/core"
	"kvfun/kvf"
)

var httpClient *http.Client

var bktLocation string = "location" // used for request BktName value

// testRecs used for some put and get funcs below.
// the "location" bucket has 85,000+ records loaded by the loader.go program
var testRecs = []core.Location{
	{
		Id:           "a123",
		Address:      "test rec 1",
		City:         "Howdy",
		St:           "TX",
		Zip:          "88223",
		LocationType: 3,
		LastActionDt: "2023-10-14",
		Notes: []string{
			"Looking Good",
		},
	},
	{
		Id:           "a124",
		Address:      "test rec 2",
		City:         "Bye",
		St:           "TX",
		Zip:          "88419",
		LocationType: 2,
		LastActionDt: "2023-10-14",
		Notes: []string{
			"Talked to Manny about wiring problem",
		},
	},
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	kvf.Debug = false

	httpClient = new(http.Client)

	//kvf.BaseURL = "http://localhost:8000/"  // to override default located in kvf/run.go - where server.go pgm is listening

	put1() // add a single record

	get1() // get the record just added

	put() // add new record and update existing record

	get() // get specific records

	delete() // delete records just added

	log.Println("*** this get should not return any records ***")
	get() // make sure records were deleted

	getAll() // get all records in bucket sorted in key order

	getAllSequence() // get all records from StartKey to EndKey in key order

	qry() // qry using both FindConditions and SortKeys

	qry2() // qry using shorthand funcs in core/util.go

	qry3() // qry using using various options

	kvf.Run(httpClient, "close", "close db")
}

func put1() {
	log.Println("-- put1 --")
	resp, err := core.PutOne(httpClient, bktLocation, &testRecs[0]) // NOTE - core.PutOne accepts single record as Go struct type ([]byte not required)
	if err != nil {
		panic("putone  request failed")
	}
	checkResp(resp, err) // checkResp func is below
}

func get1() {
	log.Println("-- get1 --")
	resp, err := core.Get(httpClient, bktLocation, testRecs[0].Id)
	checkResp(resp, err)

	var locRec core.Location
	json.Unmarshal(resp.Rec, &locRec)

	log.Printf("%+v\n", locRec)
}

func put() {
	log.Println("-- put --")

	testRecs[0].Notes = append(testRecs[0].Notes, "Rec Updated and replaced by put func") // used to verify record is replaced with changed version

	jsonRecs := make([][]byte, 0, len(testRecs))
	for _, rec := range testRecs {
		jsonRec, _ := json.Marshal(rec) // convert each record to []byte
		jsonRecs = append(jsonRecs, jsonRec)
	}
	resp, err := core.Put(httpClient, bktLocation, jsonRecs) // NOTE - unlike core.PutOne, recs must be json.Marshalled to []byte
	checkResp(resp, err)
}

func get() {
	log.Println("-- get: keys 'a123', 'a124' --")
	resp, err := core.Get(httpClient, bktLocation, "a123", "a124")
	checkResp(resp, err)

	locRecs := make([]core.Location, len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
	}
}

func delete() {
	log.Println("-- delete --")
	req := kvf.DeleteRequest{
		BktName: bktLocation,
		Keys:    []string{"a123", "a124"},
	}
	resp, err := kvf.Run(httpClient, "delete", &req)
	checkResp(resp, err)
}

func getAll() {
	log.Println("-- get all --")
	req := kvf.GetAllRequest{BktName: bktLocation} // if only bktLocation parm is specified, all recs in key order are returned
	resp, err := kvf.Run(httpClient, "getall", req)
	checkResp(resp, err)

	log.Println("getAll location bkt response count", len(resp.Recs))

	var locRec core.Location
	for i := 0; i < 10; i++ {
		json.Unmarshal(resp.Recs[i], &locRec)
		log.Printf("%d - %+v\n", i, locRec)
	}
}

func getAllSequence() {
	log.Println("-- get all sequence, should return 4 records --")
	req := kvf.GetAllRequest{
		BktName:  bktLocation,
		StartKey: "59404eebc58bd1a4ee1252ec",
		EndKey:   "59404eebc58bd1a4ee1252ef",
	}
	resp, err := kvf.Run(httpClient, "getall", req)
	checkResp(resp, err)

	locRecs := make([]core.Location, len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%d - %+v\n", i, locRecs[i])
	}
}

func qry() {
	log.Println("-- qry: find st=PA, locationType>1, sort locationType desc, city asc --")
	req := kvf.QryRequest{
		BktName: bktLocation,
		FindConditions: []kvf.FindCondition{
			{Fld: "st", Op: kvf.Matches, ValStr: "PA"},
			{Fld: "locationType", Op: kvf.GreaterThan, ValInt: 1},
		},
		SortFlds: []kvf.SortKey{
			{Fld: "locationType", Dir: kvf.DescInt},
			{Fld: "city", Dir: kvf.AscStr},
		},
	}
	resp, err := kvf.Run(httpClient, "qry", req)
	checkResp(resp, err)

	// Unmarshal response.Recs
	locRecs := make([]core.Location, len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
	}
	// show results
	log.Println("response count", len(resp.Recs))
	for i, rec := range locRecs {
		log.Printf("%d %+v\n", i, rec)
	}
	log.Println("-- qry end --")
}

func qry2() {
	log.Println("-- qry2: find by zip prefix, sortby city asc, lastActionDt desc --")

	find := core.FindStr("zip", kvf.StartsWith, "19") // returns []kvf.FindCondition with 1 condition loaded
	sortBy := core.SortBy("city", kvf.AscStr)         // returns []kvf.SortKey with 1 sortKey loaded
	sortBy = append(sortBy, kvf.SortKey{Fld: "lastActionDt", Dir: kvf.DescStr})

	resp, err := core.Qry(httpClient, bktLocation, find, sortBy)
	checkResp(resp, err)

	var locRec core.Location
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRec)
		log.Printf("%d %+v\n", i, locRec)
	}
	// using range of keys
	log.Println("-- qry2: key range, no sortBy, output should be in key order --")

	startKey := "5b6aed48c58bd1aab33a2d94"
	endKey := "5bd1ea98c58bd13a6156011c"

	resp, err = core.Qry(httpClient, bktLocation, find, nil, startKey, endKey) // no sortBy, should return in key order
	checkResp(resp, err)
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRec)
		log.Printf("%d %+v\n", i, locRec)
	}
}

func qry3() {

	var locRec core.Location

	log.Println("-- qry3: find where address contains 'west', sortby address asc --")
	find := core.FindStr("address", kvf.Contains, "west")
	sortBy := core.SortBy("address", kvf.AscStr)
	resp, err := core.Qry(httpClient, bktLocation, find, sortBy)
	checkResp(resp, err)
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRec)
		log.Printf("%d %+v\n", i, locRec)
	}

	log.Println("-- qry3: find where city startsWith 'lake', sortby address asc --")
	find = core.FindStr("city", kvf.StartsWith, "lake")
	resp, err = core.Qry(httpClient, bktLocation, find, sortBy)
	checkResp(resp, err)
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRec)
		log.Printf("%d %+v\n", i, locRec)
	}

	log.Println("-- qry3: no findConditions, sortby city, address asc --")
	sortBy = core.SortBy("city", kvf.AscStr)
	sortBy = append(sortBy, kvf.SortKey{Fld: "address", Dir: kvf.AscStr})
	resp, err = core.Qry(httpClient, bktLocation, nil, sortBy)
	checkResp(resp, err)
	log.Println("resp count", len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRec)
		log.Printf("%d %+v\n", i, locRec)
		if i > 100 {
			break
		}
	}
}

func checkResp(resp *kvf.Response, err error) bool {
	if err != nil {
		panic(err)
	}
	if resp.Status == kvf.Ok {
		return true
	}
	log.Println(kvf.StatusTxt[resp.Status], resp.Msg) // StatusTxt map in kvf/handlers.go
	return false
}
