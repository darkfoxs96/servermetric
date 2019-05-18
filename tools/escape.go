package tools

import (
	"encoding/json"
)

func JsonEscape(s string) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return b
}
