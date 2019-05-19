package web

import (
	"net/http"
	"strconv"

	"github.com/darkfoxs96/servermetric/alert"
	"github.com/darkfoxs96/servermetric/tools"
)

func DisconnectController(w http.ResponseWriter, r *http.Request) {
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

	strID := r.URL.Query().Get("id")
	ID, err := strconv.Atoi(strID)
	if err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}

	err = alert.RemoveConnections(int64(ID))
	if err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success","id":` + strconv.Itoa(int(ID)) + `}`))
}
