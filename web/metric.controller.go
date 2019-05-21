package web

import (
	"fmt"
	"github.com/darkfoxs96/servermetric/alert"
	"github.com/darkfoxs96/servermetric/db"
	"github.com/darkfoxs96/servermetric/tools"
	"net/http"
	"strings"
)

type MetricData struct {
	Fields string         `json:"fields"`
	Types  []string       `json:"types"`
	Data   *MetricDataArr `json:"data"`
}

type Metric struct {
	ServerID int64                  `json:"serverId"`
	Name     string                 `json:"name"`
	Metrics  map[string]*MetricData `json:"metrics"`
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

	for _, v := range metric.Metrics {
		err = v.Data.UnmarshalBuf(v.Types)
		if err != nil {
			tools.WriteJson(w, r, err, nil)
			return
		}
	}

	err = alert.UpdateConnections(metric.ServerID)
	if err != nil {
		tools.WriteJson(w, r, err, nil)
		return
	}

	for name, metricData := range metric.Metrics {
		if len(metricData.Data.Data) == 0 {
			continue
		}

		tx, err := db.DB.Begin()
		if err != nil {
			tools.WriteJson(w, r, err, nil)
			return
		}

		pLen := len(metricData.Data.Data[0])
		values := ""
		for i := 0; i < pLen; i++ {
			values += "?,"
		}
		values = values[:len(values)-1]

		sqlReq := `INSERT INTO ` + name + ` (` + metricData.Fields + `) VALUES (` + values + `);`
		stmt, err := tx.Prepare(sqlReq)
		if err != nil {
			if stmt != nil {
				_ = stmt.Close()
			}
			_ = tx.Rollback()
			tools.WriteJson(w, r, err, nil)
			return
		}

		for _, params := range metricData.Data.Data {
			for i, val := range params {
				if v, ok := val.(float64); ok {
					if !strings.Contains(fmt.Sprint(v), ".") {
						params[i] = int(v)
					}
				}
			}

			if _, err = stmt.Exec(params...); err != nil {
				_ = stmt.Close()
				_ = tx.Rollback()
				tools.WriteJson(w, r, err, nil)
				return
			}
		}

		_ = stmt.Close()
		if err = tx.Commit(); err != nil {
			tools.WriteJson(w, r, err, nil)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"success"}`))
}
