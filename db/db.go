package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/kshvakov/clickhouse"
)

type DBConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
}

var (
	DB       *sql.DB
	params   *DBConfig
	iterConn = 0
)

func Run(p *DBConfig) {
	params = p
	connectDB(params.Username, params.Password, params.Database, params.Host, params.Port)
}

func connectDB(login, pass, dbName, h, p string) {
	time.Sleep(5 * time.Second)
	var err error

	defer func() {
		if pan := recover(); pan != nil || err != nil {
			if iterConn < 0 {
				iterConn++
				connectDB(login, pass, dbName, h, p)
			} else {
				panic(fmt.Errorf("ERROR CONNECT TO DB: %v %v", pan, err))
			}
		}
	}()

	db, err := sql.Open("clickhouse", "tcp://"+h+":"+p+"?debug=false&compress=false&username="+login+"&password="+pass+"&database="+dbName+"&")
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	fmt.Println("DB Connected.")
	DB = db
}

func AddMetrics(metrics map[string]map[string]string) (err error) {
	var sqlStr string

	for name, params := range metrics {
		sqlStr = "CREATE TABLE IF NOT EXISTS " + name + " ("
		mode := ""

		for fieldName, typeField := range params {
			if fieldName == "ENGINE" {
				mode = typeField
				continue
			}

			sqlStr += ` "` + fieldName + `" ` + typeField + `,`
		}

		sqlStr = sqlStr[:len(sqlStr)-1]
		sqlStr += ") ENGINE = " + mode + ";"

		_, err = DB.Exec(sqlStr)
		if err != nil {
			return
		}
	}

	return
}
