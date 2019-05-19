package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var ErrBadJson = fmt.Errorf("bad json data")

func ParseJson(r *http.Request, obj interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(obj)
	_ = r.Body.Close()
	if err != nil {
		return ErrBadJson
	}

	return
}
