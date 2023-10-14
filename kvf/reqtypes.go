// File reqtypes.go contains all the request types used by client programs and processed by a server program.

package kvf

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
