// Program server.go accepts http requests from client programs and interacts with the bolt db.
// All requests use the Post method.
// All responses are instances of *kvf.Response.
// The dbHandler func calls appropriate request handler in handlers.go.

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"kvfun/kvf"

	bolt "go.etcd.io/bbolt"
)

var dbPath = "/home/jay/data/kvftest.db"
var db *bolt.DB

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.GetRequest
		dbHandler("get", &request, w, r)
	})
	http.HandleFunc("/getone", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.GetOneRequest
		dbHandler("getone", &request, w, r)
	})
	http.HandleFunc("/getall", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.GetAllRequest
		dbHandler("getall", &request, w, r)
	})
	http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.PutRequest
		dbHandler("put", &request, w, r)
	})
	http.HandleFunc("/putone", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.PutOneRequest
		dbHandler("putone", &request, w, r)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.DeleteRequest
		dbHandler("delete", &request, w, r)
	})
	http.HandleFunc("/qry", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.QryRequest
		dbHandler("qry", &request, w, r)
	})
	http.HandleFunc("/bkt", func(w http.ResponseWriter, r *http.Request) {
		var request kvf.BktRequest
		dbHandler("bkt", &request, w, r)
	})
	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	})

	log.Println(http.ListenAndServe(":8000", nil))
}

func dbHandler(op string, request any, w http.ResponseWriter, r *http.Request) {
	log.Println("request started")
	jsonContent, err := io.ReadAll(r.Body) // -> []byte
	if err != nil {
		log.Println("readall of request body failed", op, err)
		return
	}
	err = json.Unmarshal(jsonContent, request)
	if err != nil {
		log.Println("json.Unmarshal failed", op, err)
		log.Println(string(jsonContent))
		return
	}
	var response *kvf.Response
	switch op {
	case "get":
		db.View(func(tx *bolt.Tx) error {
			response = kvf.Get(tx, request.(*kvf.GetRequest))
			return nil
		})
	case "getone":
		db.View(func(tx *bolt.Tx) error {
			response = kvf.GetOne(tx, request.(*kvf.GetOneRequest))
			return nil
		})
	case "getall":
		db.View(func(tx *bolt.Tx) error {
			response = kvf.GetAll(tx, request.(*kvf.GetAllRequest))
			return nil
		})
	case "put":
		db.Update(func(tx *bolt.Tx) error {
			response = kvf.Put(tx, request.(*kvf.PutRequest))
			return nil
		})
	case "delete":
		db.Update(func(tx *bolt.Tx) error {
			response = kvf.Delete(tx, request.(*kvf.DeleteRequest))
			return nil
		})
	case "putone":
		db.Update(func(tx *bolt.Tx) error {
			response = kvf.PutOne(tx, request.(*kvf.PutOneRequest))
			return nil
		})
	case "qry":
		db.View(func(tx *bolt.Tx) error {
			response = kvf.Qry(tx, request.(*kvf.QryRequest))
			return nil
		})
	case "bkt":
		db.Update(func(tx *bolt.Tx) error {
			response = kvf.Bkt(tx, request.(*kvf.BktRequest))
			return nil
		})
	}
	jsonData, err := json.Marshal(response) // if sending response to remote requester, then compression is probably a good idea
	if err != nil {
		log.Println("json.Marshal response failed", err)
		log.Println(response)
		return
	}
	w.Write(jsonData)
	log.Println("request done")
}
