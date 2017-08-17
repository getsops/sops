package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mozilla.org/sops/keyservice"

	"google.golang.org/grpc"
)

func main() {
	var network, addr string
	flag.StringVar(&network, "net", "tcp", "Network to listen on, eg tcp or unix")
	flag.StringVar(&addr, "addr", "127.0.0.1:5000", "Address to listen on, eg 127.0.0.1:5000 or /tmp/sops.sock")
	flag.Parse()
	lis, err := net.Listen(network, addr)
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	keyservice.RegisterKeyServiceServer(grpcServer, keyservice.Server{})
	log.Printf("Listening on %s://%s", network, addr)

	// Close socket if we get killed
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		lis.Close()
		os.Exit(0)
	}(sigc)
	grpcServer.Serve(lis)
}
