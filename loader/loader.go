/*
 */
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"kvfun/kvf"
	//kvf "github.com/jayposs/kvfun"
)

var baseURL string = "http://localhost:8000/"

var httpClient *http.Client

var response kvf.Response

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

	loadData()

	var err error
	var resp *kvf.Response

	httpClient = new(http.Client)

	// DELETE / CREATE BUCKET ---------------------------------------
	bktReq := kvf.BktRequest{BktName: "location", Operation: "delete"}
	resp, _ = run("bkt", bktReq)

	bktReq.Operation = "create"
	resp, err = run("bkt", bktReq)
	if err != nil {
		panic("bkt create failed")
	}

	// PUT RECORDS INTO BUCKET -------------------------------------
	putReq := kvf.PutRequest{
		BktName:  "location",
		KeyField: "id",
		Recs:     make([][]byte, len(locationData)),
	}
	for i, rec := range locationData {
		putReq.Recs[i], _ = json.Marshal(rec) // convert each record to []byte
	}

	resp, err = run("put", putReq)
	if err != nil {
		panic("put req failed")
	}
	log.Println(resp)
}

func loadData() {
	var filePath = "/home/jay/data/properties.csv"

	file, err := os.Open(filePath)
	if err != nil {
		panic("open file failed")
	}
	reader := csv.NewReader(file)
	data, err := reader.ReadAll() // closes file, returns [][]string, error

	locationData = make([]Location, 0, len(data))
	for i := 1; i < len(data); i++ {
		rec := Location{
			Id:      data[i][0],
			Address: data[i][1],
			City:    data[i][2],
			St:      data[i][3],
			Zip:     data[i][4],
		}
		if rec.St == "TX" {
			rec.CompanyId = 1
		} else {
			rec.CompanyId = 2
		}
		locationData = append(locationData, rec)
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
