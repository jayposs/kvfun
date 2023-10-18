// The util.go file contains funcs and types shared by other programs in the kvfun module.

package core

import (
	"bytes"
	"encoding/json"
	"kvfun/kvf"
	"log"
	"net/http"
)

var KeyFieldName = "id"

// Get func provides shorthand way of calling kvf.Run with Get or GetOne request.
// If only 1 key is included in parm list, then GetOneRequest is used.
func Get(httpClient *http.Client, bktName string, keys ...string) (*kvf.Response, error) {
	var req any
	var op string
	if len(keys) == 1 {
		op = "getone"
		req = kvf.GetOneRequest{
			BktName: bktName,
			Key:     keys[0],
		}
	} else {
		op = "get"
		req = kvf.GetRequest{
			BktName: bktName,
			Keys:    keys,
		}
	}
	resp, err := kvf.Run(httpClient, op, req)
	return resp, err
}

// Get func provides shorthand way of calling kvf.Run with Put request.
// Unlike PutOne, the recs parm must already be json.Marshaled, due to problems working with []any.
func Put(httpClient *http.Client, bktName string, recs [][]byte) (*kvf.Response, error) {
	req := kvf.PutRequest{
		BktName:  bktName,
		KeyField: KeyFieldName,
		Recs:     recs,
	}
	resp, err := kvf.Run(httpClient, "put", &req)
	if resp.PutCnt != len(recs) {
		log.Printf("Put Operation Did Not Process All Recs, Sent - %d, Put - %d", len(recs), resp.PutCnt)
	}
	return resp, err
}

// PutOne func provides shorthand way of calling kvf.Run with PutOne request.
// Unlike Put func, rec will be json.Marshalled.
func PutOne(httpClient *http.Client, bktName string, rec any) (*kvf.Response, error) {
	jsonRec, err := json.Marshal(rec)
	if err != nil {
		log.Println("PutOne json.Marshal rec failed", err)
		return nil, err
	}
	req := kvf.PutOneRequest{
		BktName:  bktName,
		KeyField: KeyFieldName,
		Rec:      jsonRec,
	}
	resp, err := kvf.Run(httpClient, "putone", &req)
	if resp.PutCnt != 1 {
		log.Println("PutOne Operation Did Not Process Rec")
	}
	return resp, err
}

// format JSON in easy to view format
func FmtJSON(jsonContent []byte) string {
	var out bytes.Buffer
	json.Indent(&out, jsonContent, "", "  ")
	return out.String()
}
