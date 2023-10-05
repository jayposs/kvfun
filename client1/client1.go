/*
 */
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"kvfun/kvf"
)

var baseURL string = "http://localhost:8000/"

var httpClient *http.Client

var bktName string = "location"

type Location struct {
	Id        string `json:"id"`
	Address   string `json:"address"`
	City      string `json:"city"`
	St        string `json:"st"`
	Zip       string `json:"zip"`
	CompanyId int    `json:"companyId"`
}

var locationData []Location

func main() {

	httpClient = new(http.Client)

	get1()

	get()

	//qryAll()

	qry()

	// put1

	// put

}

func get1() {
	get1Req := kvf.GetOneRequest{
		BktName: bktName,
		Key:     "5c7ee97dc58bd17d64b36c66",
	}
	resp, err := run("getone", get1Req)
	if err != nil {
		panic("getone  request failed")
	}
	log.Println("-- get1 --")
	var locRec Location
	json.Unmarshal(resp.Rec, &locRec)
	log.Printf("%+v\n", locRec)
}

func get() {
	getReq := kvf.GetRequest{
		BktName: bktName,
		Keys: []string{
			"5c7ee97dc58bd17d64b36c66",
			"5c16f6f9c58bd111797556db",
			"5ff357d8362a3335e5fc2ee3",
		},
	}
	resp, err := run("get", getReq)
	if err != nil {
		panic("get  request failed")
	}
	log.Println("-- get --")
	locRecs := make([]Location, len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
	}
}

func qryAll() {
	qryReq := kvf.QryRequest{ // if only BktName parm is specified, all recs in key order are returned
		BktName: bktName,
	}
	resp, err := run("qry", qryReq)
	if err != nil {
		panic("qry all request failed")
	}
	log.Println("-- qry all --")
	locRecs := make([]Location, len(resp.Recs))
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
	qryReq := kvf.QryRequest{
		BktName: bktName,
		FindConditions: []kvf.FindCondition{
			{Fld: "st", Op: kvf.Matches, ValStr: "PA"},
			{Fld: "companyId", Op: kvf.EqualTo, ValInt: 2},
		},
		SortFlds: []kvf.SortKey{
			{Fld: "city", Dir: kvf.Desc},   // string desc
			{Fld: "address", Dir: kvf.Asc}, // string asc
		},
	}
	resp, err := run("qry", qryReq)
	if err != nil {
		panic("qry request failed")
	}
	log.Println("-- qry --")
	locRecs := make([]Location, len(resp.Recs))
	log.Println("response count", len(resp.Recs))
	for i, rec := range resp.Recs {
		json.Unmarshal(rec, &locRecs[i])
		log.Printf("%+v\n", locRecs[i])
	}
}

// format JSON in easy to view format
func fmtJSON(jsonContent []byte) string {
	var out bytes.Buffer
	json.Indent(&out, jsonContent, "", "  ")
	return out.String()
}

// Run method executes the api request using the provided payload.
func run(op string, payload interface{}) (*kvf.Response, error) {
	reqUrl := baseURL + op
	jsonContent, err := json.Marshal(&payload) // -> []byte

	//log.Println("--- client sending ---")
	//log.Println(fmtJSON(jsonContent))

	reqBody := bytes.NewReader(jsonContent) // -> io.Reader

	req, err := http.NewRequest("POST", reqUrl, reqBody)
	req.Header.Add("Content-Type", "application/json")

	resp, doErr := httpClient.Do(req)
	defer func() {
		if doErr == nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK || doErr != nil {
		log.Println("Request Failed, Status:", resp.StatusCode, " ", resp.Status, " - ", doErr, " --- XXX")
		return nil, doErr
	}
	result, err := io.ReadAll(resp.Body) // -> []byte
	if err != nil {
		log.Println("Read Http Response.Body Failed:", err)
	}

	//log.Println("--- client receiving ---")
	//log.Println(fmtJSON(result))

	kvfResp := new(kvf.Response)
	err = json.Unmarshal(result, kvfResp)

	return kvfResp, err
}
