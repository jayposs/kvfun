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

// GetOneRequest is used to get a specific record by Key.
type GetOneRequest struct {
	BktName string `json:"bktName"`
	Key     string `json:"key"` // key of record to be returned
}

// PutRequest is used to add or replace records.
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

type SortKey struct {
	Fld string `json:"fld"`
	Dir int    `json:"dir"`
}

// QryRequest is used to filter and sort records.
// Parms with nil or empty string value are ignored.
type QryRequest struct {
	BktName        string          `json:"bktName"`
	FindConditions []FindCondition `json:"findConditions"`
	SortFlds       []SortKey       `json:"sortFlds"` // see constants above
	StartKey       string          `json:"startKey"`
	EndKey         string          `json:"endKey"`
}
