package tools

import (
	"encoding/json"
	"fmt"

	"github.com/gocraft/web"
)

var ErrBadJson = fmt.Errorf("bad json data")

func ParseJson(r *web.Request, obj interface{}) (err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(obj)
	_ = r.Body.Close()
	if err != nil {
		return ErrBadJson
	}

	return
}
