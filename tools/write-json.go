package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func WriteJson(w http.ResponseWriter, r *http.Request, err error, d interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(Error{Status: 400, Msg: err.Error()})
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
