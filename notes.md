## Notes 
See client1/client1.go for examples of using most features.  
See loader/loader.go for example of bulk data loader.  
See server/server.go for db server program (~130 lines).  

References to "rec/record" mean the Value ( []byte ) used for Gets and Puts.  

**Sorting Options Used in Qry Request**   
Code located in kvf/handlers.go 
```
// constants used by sort code in Qry func
const (
	AscStr int = iota
	DescStr
	AscInt
	DescInt
)

type SortKey struct {
	Fld string `json:"fld"` // name (json) of field in record containing sort value
	Dir int    `json:"dir"` // direction (Asc/Desc) and field type (Str/Int)
}
```  
**Find Ops/Conditions Used in Qry Request**   
Code located in kvf/rec.go   
```
// FindCondition Ops
const (
	Contains int = iota
	Matches
	StartsWith
	LessThanStr
	GreaterThanStr
	LessThan    // int
	GreaterThan // int
	EqualTo     // int
)

// NOTE - The Op code determines if ValStr or ValInt is used for comparison
type FindCondition struct {
	Fld    string // name (json) of field in record containing compare value
	Op     int    // see constants above
	ValStr string // for Ops: Matches, StartsWith, Contains, LessThanStr, GreaterThanStr
	ValInt int    // for Ops: EqualTo, LessThan, GreaterThan
}
```  

## Steps To Add Feature  
* Add Request Type to kvf/reqtypes.go
* Add Handler Func to kvf/handlers.go
* Add Http Request handler logic to server/server.go

**Performance**  
Reading records is very fast. For large result sets, most of the time will be spent json.Marshalling the Response. Putting (add/replace) is pretty slow for large updates. Breaking large updates into smaller batches speeds things a lot. See loader/loader.go for example using goroutines to send multiple updates simultaneously. This method was much faster than sending a single large update. My test system (Intel® Core™ i3-8109U CPU @ 3.00GHz × 4 released 2018) took aprox 8 secs to load 85,000 records in single batch and less than 2 seconds to load same records in batches of 1000 using goroutines. Reading the same 85,000 records from a single bucket takes aprox 10ms.  

Note - I suspect the reason using goroutines speeds up the loader pgm is because the json.UnMarshal of the Put Request can run simultaneously with other requests. Only 1 Bolt Update transaction can run at a time.

Optimization is possible using [fastjson.Parser](https://pkg.go.dev/github.com/valyala/fastjson#Parser) when using multiple Qry FindConditions and SortKeys. Currently each field value is retrieved using fastjson.GetString or fastjson.GetInt. This change may require significant rewrite of code in kvf/rec.go.

## General Comments  

As stated in readme.md, this repo is a sample system rather than an includable package. The parent folder is kvfun. It contains the go.mod file for the module and git files. The "kvf" folder contains the types and funcs used by other programs to interact with the database. The "core" folder contains util.go and datatypes.go which have types and funcs shared by client programs. The "loader" and "client1" folders contain client programs. The "server" folder contains the server program.

Working with records stored as []bytes may cause some confusion. When putting records into the db, the individual records are json.Marshalled (converting typed record to []byte) and then the entire Request object is json.Marshalled. When getting records from the db, the entire Response object is json.Unmarshalled and then each Response.Rec is json.Unmarshalled into instance of specific record type. See client1/client1.go for examples.  

BoltDB only allows 1 program to open the database file, but allows for multiple read transactions to execute simultaneously. Client programs send requests to the server program which interacts with the database. The server uses http.ListenAndServe so multiple requests can run at the same time.

Currently the Response object created by handler funcs is json.Marshalled by the http handler func in the server pgm. References to db values cannot be in the Response.Recs after the db transaction has ended. Therefore handler funcs make a copy of each db val and save the copy in the Response. If the copy step is causing delayed responses, the logic could be changed to include the Response marshal step in the handler funcs. This change does cause some complications, but nothing major.

There are a number of "good to have" features that could be added, but would make it more complex, such as:  
* Sorting currently only allows for string and int types. Add other types.
* Secondary Indexing scheme (for large semi static data sets)  
* Finding with both and/or capability (currently only provides "and")
* Result paging
* Relational feature (I have designed a workable scheme)
* Nesting buckets (supported directly by Bolt)
* Replace "Put" with "Add", "Update", "Replace" functionality
* Depending on workload, using GOB encoding rather than JSON may be better