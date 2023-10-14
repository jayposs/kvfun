// File handlers.go contains a func to process each request type.
// These funcs are called by the dbHandler func in the server.go program.

package kvf

import (
	"cmp"
	"log"
	"slices"

	bolt "go.etcd.io/bbolt"
)

var DefaultQryRespSize = 300 // response slice initial allocation for this size

// constants used by Qry func sort logic
const (
	AscStr int = iota
	DescStr
	AscInt
	DescInt
)

// SortKey used by Qry func sort logic
type SortKey struct {
	Fld string `json:"fld"` // name of field
	Dir int    `json:"dir"` // direction (asc/desc) and field type (Str/Int)
}

// Response used for all requests
type Response struct {
	Status    string   `json:"status"` // ok or fail
	Msg       string   `json:"msg"`
	TimeStamp string   `json:"timeStamp"`
	Recs      [][]byte `json:"recs"` // for request responses with potentially more than 1 record
	Rec       []byte   `json:"rec"`  // for requests that only return 1 record
}

// Get returns recs with keys matching requested keys.
func Get(tx *bolt.Tx, req *GetRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	resp.Recs = make([][]byte, 0, 20)

	for _, key := range req.Keys {
		v := bkt.Get([]byte(key))
		if v == nil {
			log.Println("key not found", key)
			resp.Msg = "Requested Record(s) Not Found"
			continue // NOTE - THIS BEHAVIOUR MAY NOT BE APPROPRIATE FOR ALL SITUATIONS
		}
		vcopy := make([]byte, len(v)) // ref to v are invalid outside tx, so copy
		copy(vcopy, v)
		resp.Recs = append(resp.Recs, vcopy)
	}
	resp.Status = "ok"
	return resp
}

// GetOne returns a rec where key matches requested key.
func GetOne(tx *bolt.Tx, req *GetOneRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	v := bkt.Get([]byte(req.Key))
	if v == nil {
		log.Println("key not found", req.Key)
		resp.Status = "fail"
		resp.Msg = "Requested Record Not Found - " + req.Key
		return resp
	}
	resp.Rec = make([]byte, len(v)) // ref to v are invalid outside tx, so copy
	copy(resp.Rec, v)

	resp.Status = "ok"
	return resp
}

// GetAll returns all records in specified bucket.
// Optionally, Start and End keys can be included in the request.
// If StartKey != "", then result begins at 1st key >= Start key.
// If EndKey != "", then result ends at last key <= End key.
func GetAll(tx *bolt.Tx, req *GetAllRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	csr := bkt.Cursor()

	result := make([][]byte, 0, DefaultQryRespSize)

	var k, v []byte
	if req.StartKey == "" {
		k, v = csr.First()
	} else {
		k, v = csr.Seek([]byte(req.StartKey))
	}

	log.Println("GetAll read loop start")
	for k != nil {
		key := string(k)
		if req.EndKey != "" && key > req.EndKey {
			break
		}
		result = append(result, v)
		k, v = csr.Next()
	}
	log.Println("GetAll read loop done")

	resp.Recs = make([][]byte, 0, len(result))
	for _, v := range result {
		vcopy := make([]byte, len(v))
		copy(vcopy, v) // ref to v are invalid outside tx, so copy
		resp.Recs = append(resp.Recs, vcopy)
	}
	log.Println("GetAll copy vals done")

	resp.Status = "ok"
	return resp
}

// Put adds or replaces records, based on existence of key.
// The KeyField specified in the request is used as the key and this field must exist in all request.Recs.
func Put(tx *bolt.Tx, req *PutRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	for _, rec := range req.Recs { // req.Recs is [][]byte
		key := recGetStr(rec, req.KeyField)
		if key == "" {
			log.Println("key value not found in record for specified KeyField - ", req.KeyField)
			log.Println(string(rec))
			resp.Status = "fail"
			resp.Msg = "key value not found in record for specified KeyField - " + req.KeyField
			continue
		}
		err := bkt.Put([]byte(key), rec)
		if err != nil {
			log.Println("put failed", err)
			resp.Status = "fail"
			resp.Msg = "Put Request Failed - " + err.Error()
			return resp
		}
	}
	resp.Status = "ok"
	return resp
}

// PutOne adds or replaces a single record. Works same as Put.
func PutOne(tx *bolt.Tx, req *PutOneRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	key := recGetStr(req.Rec, req.KeyField)
	if key == "" {
		log.Println("key value not found in record", req.KeyField)
		resp.Status = "fail"
		resp.Msg = "key value not found in record - " + req.KeyField
		return resp
	}
	err := bkt.Put([]byte(key), req.Rec)
	if err != nil {
		log.Println("put failed", err)
		resp.Status = "fail"
		resp.Msg = "Put Request Failed - " + err.Error()
		return resp
	}
	resp.Status = "ok"
	return resp
}

// Delete deletes recs with keys matching specified keys.
func Delete(tx *bolt.Tx, req *DeleteRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	for _, key := range req.Keys {
		err := bkt.Delete([]byte(key))
		if err != nil { // key not found does not return error
			log.Println("delete error - ", key, err)
			resp.Status = "fail"
			resp.Msg = "delete error - " + key
			return resp
		}
	}
	resp.Status = "ok"
	return resp
}

// Qry returns records that meet request FindConditions and in specified sort order.
// See type SortKey and Op constants above.
func Qry(tx *bolt.Tx, req *QryRequest) *Response {

	resp := new(Response)
	bkt := openBkt(tx, resp, req.BktName)
	if bkt == nil {
		return resp
	}
	csr := bkt.Cursor()

	result := make(map[string][]byte, DefaultQryRespSize)

	var k, v []byte
	if req.StartKey == "" {
		k, v = csr.First()
	} else {
		k, v = csr.Seek([]byte(req.StartKey))
	}

	keys := make([]string, 0, DefaultQryRespSize)
	var keep bool

	log.Println("qry read loop start")
	for k != nil {
		key := string(k)
		if req.EndKey != "" && key > req.EndKey {
			break
		}
		if req.FindConditions == nil {
			keep = true
		} else {
			keep = recFind(v, req.FindConditions)
		}
		if keep {
			result[key] = v
			keys = append(keys, key)
		}
		k, v = csr.Next()
	}
	log.Println("qry read loop done")

	if req.SortFlds != nil {
		log.Println("sort start")
		slices.SortFunc(keys, func(a, b string) int {
			reca := result[a]
			recb := result[b]
			var n int
			for _, sortkey := range req.SortFlds {
				switch sortkey.Dir {
				case AscStr, DescStr: // compare string flds
					vala := recGetStr(reca, sortkey.Fld)
					valb := recGetStr(recb, sortkey.Fld)
					n = cmp.Compare(vala, valb)
				case AscInt, DescInt: // compare int flds
					vala := recGetInt(reca, sortkey.Fld)
					valb := recGetInt(recb, sortkey.Fld)
					n = cmp.Compare(vala, valb)
				}
				if n == 0 { // sort key values are equal
					continue
				}
				if sortkey.Dir == DescStr || sortkey.Dir == DescInt {
					n = n * -1
				}
				return n
			}
			return 0 // all sort key values are equal
		})
		log.Println("sort done")
	}

	// load response.Recs slice in order based on sorted order of keys
	resp.Recs = make([][]byte, 0, len(keys))
	for _, key := range keys {
		v := result[key]
		vcopy := make([]byte, len(v))
		copy(vcopy, v) // refs to returned vals are invalid outside tx, so response is loaded with copies
		resp.Recs = append(resp.Recs, vcopy)
	}
	resp.Status = "ok"
	return resp
}

// Bkt performs bucket requests such as "create" and "delete"
func Bkt(tx *bolt.Tx, req *BktRequest) *Response {

	resp := new(Response)
	var err error
	switch req.Operation {
	case "create":
		_, err = tx.CreateBucket([]byte(req.BktName))
	case "delete":
		err = tx.DeleteBucket([]byte(req.BktName))
	}
	if err != nil {
		log.Println("Bkt Operation Failed-" + req.Operation + "-" + req.BktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Operation Failed-" + req.Operation + "-" + req.BktName
		return resp
	}
	resp.Status = "ok"
	return resp
}

func openBkt(tx *bolt.Tx, resp *Response, bktName string) *bolt.Bucket {
	bkt := tx.Bucket([]byte(bktName))
	if bkt == nil {
		log.Println("Bkt Not Found - ", bktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Not Found - " + bktName
	}
	return bkt
}
