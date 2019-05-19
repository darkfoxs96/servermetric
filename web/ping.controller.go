package web

import (
	"net/http"
	"strconv"
	"time"
)

func PingController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error":"method not allowed"}`))
		return
	}

	key := r.URL.Query().Get("key")
	if key != config.Key {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"forbidden"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"pong":` + strconv.Itoa(int(time.Now().UnixNano())) + `}`))
}
