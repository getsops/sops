package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

type providerConfig struct {
	DelayMilliseconds  int    `json:"delay_ms"`
	IgnoreCancellation bool   `json:"ignore_cancellation"`
	StartedFile        string `json:"started_file"`
}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	if err := wait(ctx, request.Config); err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: append([]byte(nil), request.DataKey...),
		Metadata:         map[string]interface{}{"process_id": strconv.Itoa(os.Getpid())},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	processID, ok := request.Metadata["process_id"].(string)
	if !ok || processID != strconv.Itoa(os.Getpid()) {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("plugin process changed between wrap and unwrap")
	}
	if err := wait(ctx, request.Config); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: append([]byte(nil), request.EncryptedDataKey...)}, nil
}

func wait(ctx context.Context, config map[string]interface{}) error {
	var settings providerConfig
	if err := sopsplugin.DecodeConfig(config, &settings); err != nil {
		return err
	}
	if settings.DelayMilliseconds == 0 {
		return nil
	}
	if settings.StartedFile != "" {
		if err := os.WriteFile(settings.StartedFile, []byte("started"), 0600); err != nil {
			return fmt.Errorf("write start marker: %w", err)
		}
	}
	timer := time.NewTimer(time.Duration(settings.DelayMilliseconds) * time.Millisecond)
	defer timer.Stop()
	if settings.IgnoreCancellation {
		<-timer.C
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func main() {
	if os.Getenv("SOPS_TEST_PLUGIN_START_DELAY") == "1" {
		time.Sleep(10 * time.Second)
	}
	sopsplugin.Serve(provider{})
}
