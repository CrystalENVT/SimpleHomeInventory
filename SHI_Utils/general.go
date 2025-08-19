package shiutils

import (
	"encoding/json"
)

// pure GoLang Pretty Print - Source: https://stackoverflow.com/a/51270134
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
