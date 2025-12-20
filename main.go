// Package main is for starting the server
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/domahidizoltan/zhero/server"
)

func main() {
	srv := server.New()
	srv.Start()
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	<-quitCh
	srv.Stop()
}
