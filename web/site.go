package web

import (
	"fmt"
	"net/http"
)

type ServerConfig struct {
	Port string
	Key  string
}

var (
	config *ServerConfig
)

func Run(c *ServerConfig) {
	config = c

	http.HandleFunc("/api/ping", PingController)
	http.HandleFunc("/api/metric", MetricController)
	http.HandleFunc("/api/connect", ConnectController)
	http.HandleFunc("/api/disconnect", DisconnectController)

	//interrupt := make(chan os.Signal, 2)
	//signal.Notify(interrupt,
	//	syscall.SIGINT,
	//	syscall.SIGTERM,
	//	syscall.SIGHUP,
	//	syscall.SIGQUIT)

	srv := &http.Server{
		Addr: ":" + config.Port,
	}

	go func() {
		fmt.Println(srv.ListenAndServe())
	}()

	fmt.Println("The service is ready to listen and serve. PORT: " + config.Port)

	//sig := <-interrupt
	//fmt.Println("signal: ", sig)
	//fmt.Println("exiting")
	////
	//fmt.Println("The service is shutting down...")
}
