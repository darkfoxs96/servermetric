package tools

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocraft/web"

	"ecomm/errorR"
)

func WriteJson(w web.ResponseWriter, err error, d interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(errorR.Error{Status: 400, Msg: err.Error()})
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(d)
	if err != nil {
		fmt.Println(err)
	}
}
