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

var locationData []core.Location

func main() {

	httpClient = new(http.Client)

	//kvf.BaseURL = "http://localhost:8000/"  // to override default located in kvf/run.go - where server.go pgm is listening

	//put1() // add a single record

	//get1() // get the record just added

	// putAdd()  // add records

	// putReplace() // replace records

	//get()  // get specific records

	// delete() // delete records just added

	getAll() // get all records in bucket sorted in key order

	//getAllSequence()  // get all records from StartKey to EndKey in key order

	//qry()  // qry using both FindConditions and SortKeys

	//qrySequence()  // qry using StartKey and EndKey

	//qryNoFilter() // qry using on SortKeys

	//qryNoSort()  // qry using on FindConditions

}

func get1() {
	req := kvf.GetOneRequest{
		BktName: bktName,
		Key:     "5c7ee97dc58bd17d64b36c66",
	}
	resp, err := kvf.Run(httpClient, "getone", req)
	if err != nil {
		panic("getone  request failed")
	}
	log.Println("-- get1 --")
	var locRec core.Location
	json.Unmarshal(resp.Rec, &locRec)
	log.Printf("%+v\n", locRec)
}

func get() {
	req := kvf.GetRequest{
		BktName: bktName,
		Keys: []string{
			"5c7ee97dc58bd17d64b36c66",
			"5c16f6f9c58bd111797556db",
			"5ff357d8362a3335e5fc2ee3",
		},
	}
	resp, err := kvf.Run(httpClient, "get", req)
	if err != nil {
		panic("get  request failed")
	}
	log.Println("-- get --")
	locRecs := make([]core.Location, len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
	}
}

func getAll() {
	req := kvf.GetAllRequest{ // if only BktName parm is specified, all recs in key order are returned
		BktName: bktName,
	}
	resp, err := kvf.Run(httpClient, "getall", req)
	if err != nil {
		panic("get all request failed")
	}
	log.Println("-- get all --")
	locRecs := make([]core.Location, len(resp.Recs))
	log.Println("response count", len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
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
