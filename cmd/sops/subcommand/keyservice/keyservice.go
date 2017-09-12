package keyservice

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mozilla.org/sops/keyservice"

	"google.golang.org/grpc"
)

// Opts are the options the key service server can take
type Opts struct {
	Network string
	Address string
}

// Run runs a SOPS key service server
func Run(opts Opts) error {
	lis, err := net.Listen(opts.Network, opts.Address)
	if err != nil {
		return err
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	keyservice.RegisterKeyServiceServer(grpcServer, keyservice.Server{})
	log.Printf("Listening on %s://%s", opts.Network, opts.Address)

	// Close socket if we get killed
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		lis.Close()
		os.Exit(0)
	}(sigc)
	return grpcServer.Serve(lis)
}
