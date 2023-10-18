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

var bktName string = "location"

// testRecs used for some put and get funcs below.
// the "location" bucket has 85,000 records loaded by the loader.go program
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

	//qry()  // qry using both FindConditions and SortKeys

	//qrySequence()  // qry using StartKey and EndKey

	//qryNoFilter() // qry using on SortKeys

	//qryNoSort()  // qry using on FindConditions

}

func put1() {
	log.Println("-- put1 --")
	resp, err := core.PutOne(httpClient, bktName, &testRecs[0])
	if err != nil {
		panic("putone  request failed")
	}
	checkResp(resp)
}

func get1() {
	log.Println("-- get1 --")
	resp, err := core.Get(httpClient, bktName, testRecs[0].Id)
	if err != nil {
		panic("getone  request failed")
	}
	checkResp(resp)

	var locRec core.Location
	json.Unmarshal(resp.Rec, &locRec)

	log.Printf("%+v\n", locRec)
}

func put() {
	log.Println("-- put --")

	testRecs[0].Notes = append(testRecs[0].Notes, "Rec Updated and replaced by put func") // verify record is replaced with changed version

	jsonRecs := make([][]byte, 0, len(testRecs))
	for _, rec := range testRecs {
		jsonRec, _ := json.Marshal(rec) // convert each record to []byte
		jsonRecs = append(jsonRecs, jsonRec)
	}
	resp, err := core.Put(httpClient, bktName, jsonRecs)
	if err != nil {
		panic("put request failed")
	}
	checkResp(resp)
}

func get() {
	log.Println("-- get --")
	resp, err := core.Get(httpClient, bktName, "a123", "a124")
	if err != nil {
		panic("get  request failed")
	}
	checkResp(resp)

	locRecs := make([]core.Location, len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
	}
}

func delete() {
	log.Println("-- delete --")
	req := kvf.DeleteRequest{
		BktName: bktName,
		Keys:    []string{"a123", "a124"},
	}
	resp, err := kvf.Run(httpClient, "delete", &req)
	if err != nil {
		panic("delete  request failed")
	}
	checkResp(resp)
}

func getAll() {
	log.Println("-- get all --")
	req := kvf.GetAllRequest{BktName: bktName} // if only BktName parm is specified, all recs in key order are returned
	resp, err := kvf.Run(httpClient, "getall", req)
	if err != nil {
		panic("get all request failed")
	}
	log.Println("getAll location bkt response count", len(resp.Recs))

	var locRec core.Location
	for i := 0; i < 10; i++ {
		json.Unmarshal(resp.Recs[i], &locRec)
		log.Printf("%d - %+v\n", i, locRec)
	}
}

func getAllSequence() {
	log.Println("-- get all sequence --")
	req := kvf.GetAllRequest{BktName: bktName, StartKey: "59404eebc58bd1a4ee1252ec", EndKey: "59404eebc58bd1a4ee1252ef"}
	resp, err := kvf.Run(httpClient, "getall", req)
	if err != nil {
		panic("get all sequence request failed")
	}
	locRecs := make([]core.Location, len(resp.Recs))
	log.Println("getAllSequence location bkt response count", len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%d - %+v\n", i, locRecs[i])
		if i > 10 {
			break
		}
	}
}

func qry() {
	req := kvf.QryRequest{
		BktName: bktName,
		FindConditions: []kvf.FindCondition{
			{Fld: "st", Op: kvf.Matches, ValStr: "PA"},
			{Fld: "locationType", Op: kvf.EqualTo, ValInt: 2},
		},
		SortFlds: []kvf.SortKey{
			{Fld: "city", Dir: kvf.DescStr},
			{Fld: "address", Dir: kvf.AscStr},
		},
	}
	resp, err := kvf.Run(httpClient, "qry", req)
	if err != nil {
		panic("qry request failed")
	}
	log.Println("-- qry --")
	locRecs := make([]core.Location, len(resp.Recs))
	log.Println("response count", len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
		if i > 20 {
			break
		}
	}
}

func checkResp(resp *kvf.Response) bool {
	if resp.Status == kvf.Ok {
		return true
	}
	log.Println(kvf.StatusTxt[resp.Status], resp.Msg) // StatusTxt map in kvf/handlers.go
	return false
}
