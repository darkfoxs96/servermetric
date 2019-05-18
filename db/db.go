package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/kshvakov/clickhouse"
)

var (
	DB       *sql.DB
	username = "clickhouse"
	password = "clickhouse"
	database = "default"
	host     = "localhost"
	port     = "9001"
	iterConn = 0
)

func Run() {
	connectDB(username, password, database, host, port)
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

	db, err := sql.Open("clickhouse", "tcp://"+host+":"+port+"?debug=false&username="+login+"&password="+pass+"&database="+dbName+"&") // "user="+login+" password="+pass+" dbname="+dbName+" host="+h+" port="+p
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
