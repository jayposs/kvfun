package kvf

import (
	//"encoding/json"
	"cmp"
	"log"
	"slices"
	"sort"

	bolt "go.etcd.io/bbolt"
)

const defaultQryRespSize = 300 // response slice initial allocation for this size

// constants used by Qry sort func
const (
	AscStr int = iota
	DescStr
	AscInt
	DescInt
)

type Response struct {
	Status    string   `json:"status"` // ok or fail
	Msg       string   `json:"msg"`
	TimeStamp string   `json:"timeStamp"`
	Recs      [][]byte `json:"recs"`
	Rec       []byte   `json:"rec"`
}

func Get(tx *bolt.Tx, req *GetRequest) *Response {

	resp := new(Response)
	bkt := tx.Bucket([]byte(req.BktName))
	if bkt == nil {
		log.Println("Bkt Not Found", req.BktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Not Found " + req.BktName
		return resp
	}
	resp.Recs = make([][]byte, 0, 20)

	for _, key := range req.Keys {
		v := bkt.Get([]byte(key))
		if v == nil {
			log.Println("key not found", key)
			resp.Msg = "Requested Record(s) Not Found"
			continue
		}
		val := make([]byte, len(v)) // ref to v are invalid outside tx, so copy
		copy(val, v)
		resp.Recs = append(resp.Recs, val)
	}
	resp.Status = "ok"
	return resp
}

func GetOne(tx *bolt.Tx, req *GetOneRequest) *Response {

	resp := new(Response)
	bkt := tx.Bucket([]byte(req.BktName))
	if bkt == nil {
		log.Println("Bkt Not Found", req.BktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Not Found " + req.BktName
		return resp
	}
	v := bkt.Get([]byte(req.Key))
	if v == nil {
		log.Println("key not found", req.Key)
		resp.Msg = "Requested Record Not Found"
		return resp
	}
	val := make([]byte, len(v)) // ref to v are invalid outside tx, so copy
	copy(val, v)
	resp.Rec = val

	resp.Status = "ok"
	return resp
}

func Put(tx *bolt.Tx, req *PutRequest) *Response {

	resp := new(Response)
	bkt := tx.Bucket([]byte(req.BktName))
	if bkt == nil {
		log.Println("Bkt Not Found", req.BktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Not Found " + req.BktName
		return resp
	}
	for _, rec := range req.Recs { // req.Recs is [][]byte
		key := recGetStr(rec, req.KeyField)
		if key == "" {
			log.Println("key value not found in record", req.KeyField)
			resp.Status = "fail"
			resp.Msg = "key value not found in record - " + req.KeyField
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

func PutOne(tx *bolt.Tx, req *PutOneRequest) *Response {

	resp := new(Response)
	bkt := tx.Bucket([]byte(req.BktName))
	if bkt == nil {
		log.Println("Bkt Not Found", req.BktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Not Found " + req.BktName
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

func Qry(tx *bolt.Tx, req *QryRequest) *Response {

	resp := new(Response)
	bkt := tx.Bucket([]byte(req.BktName))
	if bkt == nil {
		log.Println("Bkt Not Found", req.BktName)
		resp.Status = "fail"
		resp.Msg = "Bkt Not Found " + req.BktName
		return resp
	}
	csr := bkt.Cursor()

	result := make(map[string][]byte, defaultQryRespSize)

	var k, v []byte
	if req.StartKey == "" {
		k, v = csr.First()
	} else {
		k, v = csr.Seek([]byte(req.StartKey))
	}

	keys := make([]string, 0, defaultQryRespSize)
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

	if req.SortFlds == nil {
		sort.Strings(keys)
	} else {
		log.Println("sort start")
		slices.SortFunc(keys, func(a, b string) int {
			reca := result[a]
			recb := result[b]
			var n int
			for _, sortkey := range req.SortFlds {
				if sortkey.Dir == AscStr || sortkey.Dir == DescStr { // compare string flds
					vala := recGetStr(reca, sortkey.Fld)
					valb := recGetStr(recb, sortkey.Fld)
					n = cmp.Compare(vala, valb)
				} else if sortkey.Dir == AscInt || sortkey.Dir == DescInt { // compare int flds
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

	// refs to returned vals are invalid outside tx, so response is loaded with copies
	// NOTE - the copy step could be eliminated if marshal process is included in handler
	resp.Recs = make([][]byte, len(keys))
	for i, key := range keys {
		v := result[key]
		rec := make([]byte, len(v))
		copy(rec, v)
		resp.Recs[i] = rec
	}
	resp.Status = "ok"
	return resp
}

func Bkt(tx *bolt.Tx, req *BktRequest) *Response {

	resp := new(Response)
	var err error
	if req.Operation == "create" {
		_, err = tx.CreateBucket([]byte(req.BktName))
	} else if req.Operation == "delete" {
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
