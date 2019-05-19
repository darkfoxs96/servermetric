package web

import (
	"github.com/darkfoxs96/servermetric/alert"
	"github.com/darkfoxs96/servermetric/db"
	"net/http"

	"github.com/darkfoxs96/servermetric/tools"
)

type Metric struct {
	ServerID int64                      `json:"serverId"`
	Name     string                     `json:"name"`
	Metrics  map[string][][]interface{} `json:"metrics"`
}

func MetricController(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
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

	metric := &Metric{}
	err := tools.ParseJson(r, metric)
	if err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}

	err = alert.UpdateConnections(metric.ServerID)
	if err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}
	defer tx.Rollback()

	for name, allParams := range metric.Metrics {
		if len(allParams) == 0 {
			continue
		}

		pLen := len(allParams[0])
		values := ""
		for i := 0; i < pLen; i++ {
			values += "?,"
		}
		values = values[:len(values)-1]

		stmt, err := tx.Prepare(`INSERT INTO ` + name + ` VALUES (` + values + `);`)
		if err != nil {
			tools.WriteJson(w, r, err, nil)
			return
		}

		for _, params := range allParams {
			if _, err = stmt.Exec(params); err != nil {
				_ = stmt.Close()
				tools.WriteJson(w, r, err, nil)
				return
			}
		}

		_ = stmt.Close()
	}

	if err = tx.Commit(); err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))
}
