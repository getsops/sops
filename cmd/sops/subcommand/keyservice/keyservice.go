package keyservice

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mozilla.org/sops/keyservice"
	"go.mozilla.org/sops/logging"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("KEYSERVICE")
}

// Opts are the options the key service server can take
type Opts struct {
	Network string
	Address string
	Prompt  bool
}

// Run runs a SOPS key service server
func Run(opts Opts) error {
	lis, err := net.Listen(opts.Network, opts.Address)
	if err != nil {
		return err
	}
	defer lis.Close()
	grpcServer := grpc.NewServer()
	keyservice.RegisterKeyServiceServer(grpcServer, keyservice.Server{
		Prompt: opts.Prompt,
	})
	log.Infof("Listening on %s://%s", opts.Network, opts.Address)

	// Close socket if we get killed
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-c
		log.Infof("Caught signal %s: shutting down.", sig)
		lis.Close()
		os.Exit(0)
	}(sigc)
	return grpcServer.Serve(lis)
}
