package kvf

import (
	"bytes"
	"encoding/json"
)

// format JSON in easy to view format
func FmtJSON(jsonContent []byte) string {
	var out bytes.Buffer
	json.Indent(&out, jsonContent, "", "  ")
	return out.String()
}
