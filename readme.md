# KVFUN - Making it fun to work with key-value databases 
Written for BoltDB, but should be fairly easy to modify for other kv databases.  

**Requires Go version 1.21 or higher**  

The primary functional improvements over straight BoltDB are:
* Search capability - result set only includes records meeting Find Conditions (thanks fastjson) 
* Sorting - by mutiple fields by name and ascending/descending by field
* Server - separate program that processes requests from multiple client programs  

See notes.md for more detailed information.

KVFUN module is an example system with multiple programs and packages demonstrating a more functional use of a KV database. 

It is not an installable package. The folder "kvf" could be used as a stand alone package. The code in the 4 kvf Go files is short and simple (IMO), so pretty easy to modify for your needs and tastes. 
   
A single running server program interacts with the database. The server accepts http requests from any number of client programs. See server.go, client1.go, and loader.go. 
 
**See client1/client1.go for examples of how to use the kvf features and interact with server program.**

This module contains the following folders/files:  
* kvf - types and funcs used by server and client programs
    * rec.go - funcs to extract vals from and search individual record  
    * handlers.go - func for each op (get, put, qry, ...)
    * reqtypes.go - struct type for each request (get, put, qry, ...)
	* run.go - func used by client pgms to send request to server pgm
* server 
    * server.go - interacts with the db and accepts requests from client pgms     
* loader 
    * loader.go - example client pgm that bulk loads data from csv file 
* client1
    * client1.go - example client pgm that demonstrates use of all request types    

## Example Request    

```
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

All sub folders/files are in a single top level folder (kvfun for this repo). The go.mod and git files are at this level.  
Third party pkgs are: 
* go.etcd.io/bbolt
* github.com/valyala/fastjson  

**Requires Go version 1.21 or greater.**  
The reason for this requirement is to take advantage of new features that offer a simplified way of sorting.  

### Speed  
Bolt DB (based on lmdb) is badass fast when it comes to reading records. See notes.md for more info. 
  
### Testing  
Testing was done with a single bucket holding 85,000 records loaded from a .csv file. 
Programs server.go, loader.go, and client1.go worked as intended.
 
Review the code and write your own tests to ensure reliability for your needs. Code is small and simple :).  
BoltDB has a good track record for reliability and kvfun adds a very thin layer on top.

## Special Thanks  
To the creators of bolt/bbolt, LMDB, and fastjson !!
