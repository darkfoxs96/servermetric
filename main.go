package main

import (
	"fmt"
	"github.com/darkfoxs96/servermetric/alert"
	"github.com/darkfoxs96/servermetric/db"
	_ "github.com/darkfoxs96/servermetric/pusher/pushers/telegram"
	"gopkg.in/yaml.v2"
	"time"
)

type GlobalConfig struct {
	AlertConfig *alert.AlertConfig
	Pushers     map[string]map[string]interface{}
}

func main() {
	db.Run()

	valYaml := `
alertconfig:
 alerts:
  - if: "SELECT * FROM table_schema WHERE timestamp > {{sub .TimestampSec 360000}}"
    then: "ok {{.V1}} | {{.V2}} | {{.V3}}"
    else: "dont't ok"
  - if: "SELECT * FROM table_schema WHERE timestamp > {{sub .TimestampSec 360}}"
    then: "dont't ok"
    else: "ok"
  - if: "SELECT AVG(cpu_use) FROM table_schema HAVING AVG(cpu_use) > 0.6"
    then: "warning avg cpu: {{.V1}}"
    else: ""
 pushers: [ "telegram" ]
 checkeveryseconds: 0
 checkconnservereveryseconds: 5
pushers:
 telegram:
  token: "841563697:AAEDpNQBkNpFSUtae_ZgSRhxKzeJRntdrik"
  data: "/Users/peterkorotkiy/docker/ex/telegram-data.cfg"
`
	data := &GlobalConfig{}
	err := yaml.Unmarshal([]byte(valYaml), data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data.AlertConfig.Alerts)

	alert.Run(data.AlertConfig)

	time.Sleep(time.Hour)
	//
	//	err = pusher.Run(data)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	time.Sleep(time.Hour)
	//tx, err := db.DB.Begin()
	//if err != nil {
	//	fmt.Println("err tx", err)
	//	return
	//}
	//
	//stmt, err := tx.Prepare(`INSERT INTO table_schema VALUES (?, ?, ?);`)
	//if err != nil {
	//	fmt.Println("err in", err)
	//	return
	//}
	//defer stmt.Close()
	//
	//if _, err = stmt.Exec(
	//	uint64(time.Now().Unix()),
	//	0.5,
	//	"site2",
	//); err != nil {
	//	panic(err)
	//}

	//if err := tx.Commit(); err != nil {
	//	panic(err)
	//}

	//_, err := db.DB.Query(`SELECT timestamp, cpu_use, host_name FROM table_schema WHERE cpu_use > 0.8`)
	//if err != nil {
	//	panic(err)
	//	return
	//}
}
