// File kvftypes.go contains constants, Response type and Request types used by client programs and processed by a server program.

package kvf

// Response Status Values
const (
	Ok int = iota
	Fail
	Warning
)

var StatusTxt = map[int]string{
	0: "Ok",
	1: "Fail",
	2: "Warning",
}

// Response used for all requests
type Response struct {
	Status int      `json:"status"` // see constants above Ok, Warning, Fail
	Msg    string   `json:"msg"`
	Recs   [][]byte `json:"recs"`   // for request responses with potentially more than 1 record
	Rec    []byte   `json:"rec"`    // for requests that only return 1 record
	PutCnt int      `json:"putCnt"` // number of records either added or replaced by Put operation
}

// Constants used in QryRequest.SortFlds and by handlers.Qry() sort logic
const (
	AscStr int = iota
	DescStr
	AscInt
	DescInt
)

// SortKey used in QryRequest.SortFlds and by handlers.Qry() sort logic
type SortKey struct {
	Fld string `json:"fld"` // name of field
	Dir int    `json:"dir"` // direction (asc/desc) and field type (Str/Int)
}

// -----  REMAINDER OF FILE IS REQUEST TYPES CREATED BY CLIENTS AND PROCESSED BY SERVER ------

// BktRequest is used to create or delete bkt.
type BktRequest struct {
	BktName   string `json:"bktName"`
	Operation string `json:"operation"` // "create", "delete"
}

// GetRequest is used to get specific records by Key.
type GetRequest struct {
	BktName string   `json:"bktName"`
	Keys    []string `json:"keys"` // keys of records to be returned
}

// GetAllRequest is used to get all records in bucket ordered by key.
// Use StartKey/EndKey to get all records in a range.
type GetAllRequest struct {
	BktName  string `json:"bktName"`
	StartKey string `json:"startKey"`
	EndKey   string `json:"endKey"`
}

// GetOneRequest is used to get a specific record by Key.
type GetOneRequest struct {
	BktName string `json:"bktName"`
	Key     string `json:"key"` // key of record to be returned
}

// PutRequest is used to add or replace records. If key exists, existing record is replaced.
type PutRequest struct {
	BktName  string   `json:"bktName"`
	KeyField string   `json:"keyField"` // field in Rec containing value to be used as key
	Recs     [][]byte `json:"recs"`     // records to be added or replaced in db
}

// PutOneRequest is used to add or replace a record.
type PutOneRequest struct {
	BktName  string `json:"bktName"`
	KeyField string `json:"keyField"` // field in Rec containing value to be used as key
	Rec      []byte `json:"rec"`      // record to be added or replaced in db
}

// DeleteRequest is used to delete specific records by Key.
type DeleteRequest struct {
	BktName string   `json:"bktName"`
	Keys    []string `json:"keys"` // keys of records to be deleted
}

// QryRequest is used to filter and sort records.
// Parms with nil or empty string value are ignored.
type QryRequest struct {
	BktName        string          `json:"bktName"`
	FindConditions []FindCondition `json:"findConditions"` // FindCondition type defined in rec.go (where find logic is located)
	SortFlds       []SortKey       `json:"sortFlds"`       // SortKey type defined in handlers.go, see Qry func
	StartKey       string          `json:"startKey"`
	EndKey         string          `json:"endKey"`
}
