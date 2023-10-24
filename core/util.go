// The util.go file contains funcs and types shared by other programs in the kvfun module.

package core

import (
	"bytes"
	"encoding/json"
	"kvfun/kvf"
	"log"
	"net/http"
)

var KeyFieldName = "id" // json name of field name in record, containing key value for Put/PutOne requests

// FmtJSON formats json into easy to view format
func FmtJSON(jsonContent []byte) string {
	var out bytes.Buffer
	json.Indent(&out, jsonContent, "", "  ")
	return out.String()
}

// Get provides shorthand way of calling kvf.Run with Get or GetOne request.
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

// Put provides shorthand way of calling kvf.Run with Put request.
// KeyFieldName is loaded from var KeyFieldName above. This field must be defined in each record.
// Unlike PutOne, recs in input recs parm must already be json.Marshaled, due to problems working with []any.
func Put(httpClient *http.Client, bktName string, recs [][]byte) (*kvf.Response, error) {
	req := kvf.PutRequest{
		BktName:  bktName,
		KeyField: KeyFieldName, // var KeyFieldName is defined above
		Recs:     recs,
	}
	resp, err := kvf.Run(httpClient, "put", &req)
	if resp.PutCnt != len(recs) {
		log.Printf("Put Operation Did Not Process All Recs, Sent - %d, Put - %d", len(recs), resp.PutCnt)
	}
	return resp, err
}

// PutOne provides shorthand way of calling kvf.Run with PutOne request.
// KeyFieldName is loaded from var KeyFieldName above. This field must be defined in each record.
// Unlike Put func, input rec is struct type value and will be json.Marshalled.
func PutOne(httpClient *http.Client, bktName string, rec any) (*kvf.Response, error) {
	jsonRec, err := json.Marshal(rec)
	if err != nil {
		log.Println("PutOne json.Marshal rec failed", err)
		return nil, err
	}
	req := kvf.PutOneRequest{
		BktName:  bktName,
		KeyField: KeyFieldName, // var KeyFieldName is defined above
		Rec:      jsonRec,
	}
	resp, err := kvf.Run(httpClient, "putone", &req)
	if resp.PutCnt != 1 {
		log.Println("PutOne Operation Did Not Process Rec")
	}
	return resp, err
}

// Qry provides shorthand way of calling kvf.Run with Qry request.
// If either findConditions or sortFlds are not needed, call with nil value.
// If 1 startEndKey is passed, then StartKey is set in request.
// If 2 startEndKeys are passed, then both StartKey and EndKey are set in request.
func Qry(httpClient *http.Client, bktName string, findConditions []kvf.FindCondition, sortFlds []kvf.SortKey, startEndKeys ...string) (*kvf.Response, error) {
	req := kvf.QryRequest{
		BktName:        bktName,
		FindConditions: findConditions,
		SortFlds:       sortFlds,
	}
	if len(startEndKeys) > 0 {
		req.StartKey = startEndKeys[0]
	}
	if len(startEndKeys) > 1 {
		req.EndKey = startEndKeys[1]
	}
	resp, err := kvf.Run(httpClient, "qry", &req)
	return resp, err
}

// FindStr creates []kvf.FindCondition with 1 str condition loaded
func FindStr(fld string, op int, val string) []kvf.FindCondition {
	findConditions := make([]kvf.FindCondition, 0, 5)
	findConditions = append(findConditions, kvf.FindCondition{Fld: fld, Op: op, ValStr: val})
	return findConditions
}

// FindInt creates []kvf.FindCondition with 1 int condition loaded
func FindInt(fld string, op int, val int) []kvf.FindCondition {
	findConditions := make([]kvf.FindCondition, 0, 5)
	findConditions = append(findConditions, kvf.FindCondition{Fld: fld, Op: op, ValInt: val})
	return findConditions
}

// SortBy creates []kvf.SortKey with 1 Fld/Dir loaded
func SortBy(fld string, dir int) []kvf.SortKey {
	sortFlds := make([]kvf.SortKey, 0, 5)
	sortFlds = append(sortFlds, kvf.SortKey{Fld: fld, Dir: dir})
	return sortFlds
}
