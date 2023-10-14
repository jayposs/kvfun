// Used by all client programs to make requests to a running server program.

package kvf

import (
	"bytes"
	"encoding/json"
	"io"
	"kvfun/core"
	"log"
	"net/http"
)

var BaseURL string = "http://localhost:8000/" // client pgm can override default if needed
var Debug bool                                // set by client to turn on debugging

// Run func executes the api request using the provided payload.
func Run(httpClient *http.Client, op string, payload interface{}) (*Response, error) {
	reqUrl := BaseURL + op
	jsonContent, err := json.Marshal(&payload) // -> []byte

	if Debug {
		log.Println("--- client sending ---")
		log.Println(core.FmtJSON(jsonContent))
	}

	reqBody := bytes.NewReader(jsonContent) // -> io.Reader

	req, err := http.NewRequest("POST", reqUrl, reqBody)
	req.Header.Add("Content-Type", "application/json")

	resp, doErr := httpClient.Do(req)
	defer func() {
		if doErr == nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != http.StatusOK || doErr != nil {
		log.Println("Request Failed, Status:", resp.StatusCode, " ", resp.Status, " - ", doErr, " --- XXX")
		return nil, doErr
	}
	result, err := io.ReadAll(resp.Body) // -> []byte
	if err != nil {
		log.Println("Read Http Response.Body Failed:", err)
	}

	if Debug {
		log.Println("--- client receiving ---")
		log.Println(core.FmtJSON(result))
	}

	kvfResp := new(Response)
	err = json.Unmarshal(result, kvfResp)

	return kvfResp, err
}
