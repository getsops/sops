package keyservice

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/logging"

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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	grpcServer := grpc.NewServer()
	server := keyservice.NewServer(ctx)
	server.Prompt = opts.Prompt
	defer server.Close()
	keyservice.RegisterKeyServiceServer(grpcServer, server)
	log.Infof("Listening on %s://%s", opts.Network, opts.Address)

	serveDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			log.Info("Shutdown requested; stopping key service.")
			grpcServer.GracefulStop()
		case <-serveDone:
		}
	}()
	err = grpcServer.Serve(lis)
	close(serveDone)
	if err == grpc.ErrServerStopped {
		return nil
	}
	return err
}
