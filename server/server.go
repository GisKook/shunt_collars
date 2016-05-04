package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/giskook/gotcp"
	"github.com/giskook/shunt_collars"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// read configuration
	shunt_collars.Config, _ = shunt_collars.ReadConfig("./conf.json")

	port := shunt_collars.Config.ServerConfig.BindPort

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+fmt.Sprint(port))
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &shunt_collars.Callback{}, &shunt_collars.TrackerProtocol{})

	// starts service
	srv.Start(listener, time.Second)
	log.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
