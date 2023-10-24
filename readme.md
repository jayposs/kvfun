## KVFun - *Fun*damentally Better Way to Use Key-Value Databases 

Developed with A.I. (American Ingenuity)  **You'll be thinking WTF (Wow That's Fast)**  

### Special Thanks - To the creators of bolt/bbolt, LMDB, fastjson, and Go !!  

Written for BoltDB, but should be fairly easy to modify for other kv databases.  

The primary functional improvements over straight BoltDB are:
* Search capability - result set only includes records meeting Find Conditions (thanks fastjson) 
* Sorting - by mutiple fields by name and ascending/descending by field
* Server - separate program that processes requests from multiple client programs  

See [notes.md](notes.md) for more detailed information.  

**Other Options** with similar goals are [BoltHold](https://github.com/timshannon/bolthold) and [Storm](https://github.com/asdine/storm) (both have way more features and way more code). 

*Project Targets* - Small, Simple, Useful  

*Code Flow* - 
Client pgm sends http request to running server pgm. Server pgm calls func in handlers.go which processes request. 

*Review these 3 code files to see if/how targets were achieved:*    
* server/server.go - uses http.ListenAndServe to process client requests by calling funcs in handlers.go
* kvf/handlers.go - 1 func for each request type (get, get1, put, put1, qry, ...)
* kvf/rec.go - handler support funcs
  
---  

A single running server program interacts with the database. The server accepts http requests from any number of client programs. See server.go, client1.go, and loader.go. 
 
**See client1/client1.go for examples of how to use the kvf features and interact with server program.**

This module contains the following folders/files:  
* kvf - types and funcs used by server and client programs  
    * rec.go - funcs to extract vals from and search individual record  
    * handlers.go - func for each op (get, put, qry, ...)
    * kvftypes.go - types and constants, primarily struct types for requests (get, put, qry, ...)
	* run.go - func used by client pgms to send request to server pgm
* server 
    * server.go - interacts with the db and accepts requests from client pgms     
* loader 
    * loader.go - example client pgm that bulk loads data from csv file 
* client1
    * client1.go - example client pgm that demonstrates use of all request types  
* core
	* util.go - request builder shortcut funcs and other util funcs shared by clients
	* datatypes.go - struct types for each db record type	  

## Example Request    

```
// see client1/client1.go for shortcuts to creating requests  
req := kvf.QryRequest{
	BktName: "location",
	FindConditions: []kvf.FindCondition{
		{Fld: "st", Op: kvf.Matches, ValStr: "PA"},
		{Fld: "locationType", Op: kvf.GreaterThan, ValInt: 1},
	},
	SortFlds: []kvf.SortKey{
		{Fld: "locationType", Dir: kvf.AscInt}, // ascending int
		{Fld: "city", Dir: kvf.DescStr},        // descending string
	},
}
resp, err := kvf.Run(httpClient, "qry", req)
```  

## Purpose 
Take advantage of simplicity and speed of kev-value databases while adding useful functionality. 

The code size is very small. It is intended to be a good starting point that is easy to extend/modify as needed.  

My vision is to use a kv db as a transaction database, serving web apps and batch applications. Data from this db would be synced to other tools such as Snowflake or Airtable to provide more robust data analysis/user interaction capabilities.

### Setup  
KVFUN module is an example system with multiple programs and packages demonstrating a more FUNctional use of a KV database. The idea is for you to create a Go Module containing the kvf, server, and custom core packages plus client programs. Or pick and choose elements for an existing project.

It is not an installable package. The folder "kvf" could be used as a stand alone package. The code in the 4 kvf Go files is short and simple (IMO:), so pretty easy to modify for your needs and tastes. 
   
All sub folders/files are in a single top level folder (kvfun for this repo). The go.mod and git files are at this level.  
Third party pkgs are: 
* go.etcd.io/bbolt
* github.com/valyala/fastjson  

See directions for [cloning a github repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository).

**Requires Go version 1.21 or greater.**  
The reason for this requirement is to take advantage of new features that offer a simplified way of sorting.  

### Speed  
Bolt DB (based on lmdb) is badass fast when it comes to reading records. See notes.md for more info. 
  
### Testing  
Testing was done with a single bucket holding 85,000 records loaded from a .csv file. 
Programs server.go, loader.go, and client1.go work as intended.
 
Review the code and write your own tests to ensure reliability for your needs. 
BoltDB has a good track record for reliability and kvfun adds a very thin layer on top.
