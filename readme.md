# KVFUN - Making it fun to work with key-value databases 
Written for bolt.db, but should be fairly easy to modify for other kv databases.   
This module contains the following folders/files:  
* kvf - types and funcs used by server and client programs
    * rec.go - funcs to extract vals from and search individual record  
    * handlers.go - func for each op (get, put, qry, ...)
    * reqtypes.go - struct type for each request
* server 
    * server.go - interacts with the db and accepts requests from client pgms     
* loader 
    * loader.go - example client pgm that bulk loads data from csv file 
* client1
    * client1.go - example client pgm that runs assortment of requests   

## Example Request    


change find and sort constants xxxxxxxxxxxxxxxxx

```
req := kvf.QryRequest{
	BktName: "location",
	FindConditions: []kvf.FindCondition{
		{Fld: "st", Op: kvf.Matches, ValStr: "PA"},
		{Fld: "companyId", Op: kvf.EqualTo, ValInt: 2},
	},
	SortFlds: []kvf.SortKey{
		{Fld: "city", Dir: kvf.Desc},   // string desc
		{Fld: "address", Dir: kvf.Asc}, // string asc
	},
}
resp, err := run("qry", req)
```  
A single server program interacts with the database. The server accepts http requests from any number of client programs. See server.go, client1.go and loader.go. 

## Purpose 
Take advantage of simplicity and speed of kev-value databases while adding useful functionality. This module has the following features:  
* Put and Get records
* Qry records using optional find conditions and sort keys.  

The code size is very small. It is intended to be a good starting point that is easy to extend as needed.  

My vision is to use this tool for the primary transactional database, serving web apps and batch applications. Data from this db would be uploaded to other tools such as Snowflake or Airtable to provide more robust data analysis capabilities.

## Setup  

I have a single folder called kvfun. The go.mod file is at this level. The other folders are inside kvfun. Third party pkgs are: 
* go.etcd.io/bbolt
* github.com/valyala/fastjson  

**NOTE - Requires Go version 1.21 or greater.**  
The reason for this requirement was to take advantage of the new "slices" pkg. It offers a more simplified way of sorting.  

## Speed  
Bolt DB (based on lmdb) is badass fast when it comes to reading records.

## Special Thanks  
To the creators of bolt/bbolt, LMDB, and fastjson !!
