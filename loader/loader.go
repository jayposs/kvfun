// Program loader.go puts test data into "location" bucket.
// Provides example of bulk data loader.
// File is aprox 85,000 records and is used for all testing.

package main

import (
	"encoding/csv"
	"encoding/json"
	"kvfun/core"
	"kvfun/kvf"
	"log"
	"net/http"
	"os"

	"sync"
	"time"
)

var locationData []core.Location // loaded with data from csv file by loadData func below

var httpClient *http.Client

func main() {

	// kvf.Debug = true

	var wg sync.WaitGroup

	loadData() // loads data from .csv into locationData slice

	var err error
	var resp *kvf.Response

	httpClient = new(http.Client)

	// DELETE / CREATE BUCKET ---------------------------------------
	bktReq := kvf.BktRequest{BktName: "location", Operation: "delete"}
	resp, _ = kvf.Run(httpClient, "bkt", bktReq)

	bktReq.Operation = "create"
	resp, err = kvf.Run(httpClient, "bkt", bktReq)
	if err != nil {
		log.Panicln("bkt create failed", err, resp.Msg)
	}

	// PUT RECORDS INTO BUCKET -------------------------------------

	// upload records to db in batches of batchSize records, using goroutines
	batchSize := 1000
	putReq := newPutReq(batchSize)
	for _, rec := range locationData {
		jsonRec, err := json.Marshal(rec) // convert each record to []byte
		if err != nil {
			log.Fatalln("json.Marshal failed", err)
		}
		putReq.Recs = append(putReq.Recs, jsonRec)
		if len(putReq.Recs) == batchSize {
			wg.Add(1)
			go run(putReq, &wg)
			putReq = newPutReq(batchSize)
			time.Sleep(10 * time.Millisecond) // pause may be appropriate for large number of requests
		}
	}
	if len(putReq.Recs) > 0 {
		wg.Add(1)
		go run(putReq, &wg)
	}

	wg.Wait() // wait for all runs to finish before ending program
}

func newPutReq(batchSize int) *kvf.PutRequest {
	return &kvf.PutRequest{
		BktName:  "location",
		KeyField: "id",
		Recs:     make([][]byte, 0, batchSize),
	}
}

func loadData() {
	var filePath = "/home/jay/data/properties.csv"
	file, err := os.Open(filePath)
	if err != nil {
		panic("open csv file failed")
	}
	reader := csv.NewReader(file)
	csvRecs, err := reader.ReadAll()

	locationData = make([]core.Location, 0, len(csvRecs))

	var x int // used to provide random values
	for i, csvRec := range csvRecs {
		if i == 0 { // skip header
			continue
		}
		if csvRec[0] == "" {
			continue
		}
		locRec := core.Location{
			Id:      csvRec[0],
			Address: csvRec[1],
			City:    csvRec[2],
			St:      csvRec[3],
			Zip:     csvRec[4],
			Notes: []string{
				"Note #1",
				"Note #2",
			},
		}
		if x < 100 {
			locRec.LocationType = 1
			locRec.LastActionDt = "2021-03-22"
		} else if x < 200 {
			locRec.LocationType = 2
			locRec.LastActionDt = "2022-06-10"
		} else if x < 300 {
			locRec.LocationType = 3
			locRec.LastActionDt = "2023-09-01"
		} else {
			x = 0
		}
		x++
		locationData = append(locationData, locRec)
	}
}

func run(req *kvf.PutRequest, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := kvf.Run(httpClient, "put", req)
	if err != nil {
		panic("put req failed")
	}
	if resp.Status != kvf.Ok {
		log.Println("ERROR", kvf.StatusTxt[resp.Status], resp.Msg)
	}
}
