//	Program loader.go puts test data into "location" bucket. File size is aprox 85,000 records and is used for all testing.

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

var locationData []core.Location

var httpClient *http.Client

var wg sync.WaitGroup

func main() {

	// kvf.Debug = true

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
	putReq := kvf.PutRequest{
		BktName:  "location",
		KeyField: "id",
	}
	// upload records to db in batches of batchSize records, using goroutines, pause between runs
	batchSize := 1000
	var cnt int
	jsonRecs := make([][]byte, 0, batchSize)
	for _, rec := range locationData {
		jsonRec, _ := json.Marshal(rec) // convert each record to []byte
		jsonRecs = append(jsonRecs, jsonRec)
		cnt++
		if cnt == batchSize {
			cnt = 0
			putReq.Recs = make([][]byte, batchSize)
			copy(putReq.Recs, jsonRecs)
			wg.Add(1)
			go run(&putReq)
			jsonRecs = make([][]byte, 0, batchSize)
			time.Sleep(10 * time.Millisecond)
		}
	}
	if cnt > 0 {
		putReq.Recs = make([][]byte, cnt)
		copy(putReq.Recs, jsonRecs)
		wg.Add(1)
		go run(&putReq)
	}

	wg.Wait() // wait for all runs to finish before ending program
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
		}
		x++
		if x > 400 {
			x = 0
		}
		locationData = append(locationData, locRec)

		// log.Println(locRec)
	}
}

func run(req *kvf.PutRequest) {
	defer wg.Done()
	resp, err := kvf.Run(httpClient, "put", req)
	if err != nil {
		panic("put req failed")
	}
	log.Println(resp)
}
