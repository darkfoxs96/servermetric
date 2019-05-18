package alert

import (
	"encoding/json"
	"fmt"
	"github.com/darkfoxs96/servermetric/pusher"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/darkfoxs96/servermetric/alert/alertmodel"
	"github.com/darkfoxs96/servermetric/db"
)

type ForDB struct {
	TimestampSec     int64
	TimestampNanoSec int64
}

type AlertParams struct {
	IF   string
	THEN string
	ELSE string
}

type AlertConfig struct {
	Alerts                      []*AlertParams
	Pushers                     []string
	CheckEverySeconds           int64
	CheckConnServerEverySeconds int64
}

type ServerConnect struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	LastConnect int64  `json:"lastConnect,omitempty"`
}

// Errors
var (
	ErrNotFoundID = fmt.Errorf("Not found connection by id")
)

var (
	alertConfig *AlertConfig
	// alerts
	alertList      = []*alertmodel.Alert{}
	alertListMutex = sync.RWMutex{}
	// server connections
	serverConnMapPath        = os.Getenv("SERVERMETRICCONNMAP")
	serverConnMap            = map[int64]*ServerConnect{}
	serverConnMapMutex       = sync.RWMutex{}
	serverConnMapIncr  int64 = 0
)

func Run(config *AlertConfig) {
	alertConfig = config

	for _, v := range alertConfig.Alerts {
		if _, err := AppendAlert(v); err != nil {
			panic(err)
		}
	}

	go func() {
		for {
			if alertConfig.CheckEverySeconds <= 0 {
				time.Sleep(time.Second * 30)
				continue
			}

			time.Sleep(time.Second * time.Duration(alertConfig.CheckEverySeconds))

			if err := checkAlerts(); err != nil {
				panic(err)
			}
		}
	}()

	if serverConnMapPath != "" {
		b, err := ioutil.ReadFile(serverConnMapPath)
		if err != nil {
			panic(err)
		}

		serverConnMapMutex.Lock()
		if err = json.Unmarshal(b, &serverConnMap); err != nil {
			serverConnMapMutex.Unlock()
			panic(err)
		}

		for ID, _ := range serverConnMap {
			if ID > serverConnMapIncr {
				serverConnMapIncr = ID
			}
		}
		serverConnMapMutex.Unlock()
	}

	go func() {
		for {
			if alertConfig.CheckConnServerEverySeconds <= 0 {
				time.Sleep(time.Second * 30)
				continue
			}

			time.Sleep(time.Second * time.Duration(alertConfig.CheckConnServerEverySeconds))

			desConn, err := checkServerConnections()
			if err != nil {
				panic(err)
			}

			err = RemoveConnections(desConn...)
			if err != nil {
				panic(err)
			}
		}
	}()

}

// AppendConnection new alert
func AppendAlert(alertParams *AlertParams) (alert *alertmodel.Alert, err error) {
	alert, err = alertmodel.NewAlert(alertParams.IF, alertParams.THEN, alertParams.ELSE)
	if err != nil {
		return
	}

	alertListMutex.Lock()
	alertList = append(alertList, alert)
	alertListMutex.Unlock()

	return
}

// AppendConnection new conn
func AppendConnection(conn *ServerConnect) (err error, ID int64) {
	conn.LastConnect = time.Now().Unix()

	serverConnMapMutex.Lock()
	serverConnMapIncr++
	ID = serverConnMapIncr
	serverConnMap[ID] = conn
	serverConnMapMutex.Unlock()

	err = saveConnections()
	return
}

// UpdateConnections by id
func UpdateConnections(IDs ...int64) (err error) {
	t := time.Now().Unix()

	serverConnMapMutex.Lock()
	for _, ID := range IDs {
		conn := serverConnMap[ID]
		if conn == nil {
			err = ErrNotFoundID
			return
		}

		conn.LastConnect = t
	}
	serverConnMapMutex.Unlock()

	return saveConnections()
}

// RemoveConnections by id
func RemoveConnections(IDs ...int64) (err error) {
	serverConnMapMutex.Lock()
	for _, ID := range IDs {
		delete(serverConnMap, ID)
	}
	serverConnMapMutex.Unlock()

	return saveConnections()
}

func saveConnections() (err error) {
	if serverConnMapPath == "" {
		return
	}

	defer serverConnMapMutex.RUnlock()
	serverConnMapMutex.RLock()
	b, err := json.Marshal(serverConnMap)
	if err != nil {
		return
	}

	return ioutil.WriteFile(serverConnMapPath, b, 0644)
}

func checkAlerts() (err error) {
	defer alertListMutex.RUnlock()
	alertListMutex.RLock()

	t := time.Now()
	data := &ForDB{
		TimestampSec:     t.Unix(),
		TimestampNanoSec: t.UnixNano(),
	}

	for _, alert := range alertList {
		msgTHEN, msgELSE, err2 := alert.Check(db.DB, data)
		if err2 != nil {
			return err2
		}

		for _, msg := range msgTHEN {
			for _, pusherName := range alertConfig.Pushers {
				err = pusher.Push(pusherName, msg)
				if err != nil {
					return
				}
			}
		}

		for _, msg := range msgELSE {
			for _, pusherName := range alertConfig.Pushers {
				err = pusher.Push(pusherName, msg)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func checkServerConnections() (desConn []int64, err error) {
	defer serverConnMapMutex.RUnlock()
	serverConnMapMutex.RLock()
	t := time.Now().Unix()
	desConn = make([]int64, 0)

	for id, conn := range serverConnMap {
		if conn.LastConnect+alertConfig.CheckConnServerEverySeconds < t {
			for _, pusherName := range alertConfig.Pushers {
				err = pusher.Push(pusherName, "Server is not responding. Host: "+conn.Host+" Name: "+conn.Name)
				if err != nil {
					return
				}
			}

			desConn = append(desConn, id)
		}
	}

	return
}
