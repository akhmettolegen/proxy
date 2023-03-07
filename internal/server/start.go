package server

import (
	"context"
	"errors"
	"fmt"
	httpCli "github.com/akhmettolegen/proxy/internal/clients/http"
	"github.com/akhmettolegen/proxy/internal/managers/proxy"
	"github.com/akhmettolegen/proxy/internal/server/configs"
	"github.com/akhmettolegen/proxy/internal/server/http"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	appCtx, appCtxCancel := context.WithCancel(context.Background())
	defer appCtxCancel()

	go catchForTermination(appCtxCancel, os.Interrupt, syscall.SIGTERM)

	opts := configs.ConfigWithParsedFlags()

	client := httpCli.NewClient(appCtx)

	proxyManager := proxy.NewManager(appCtx, client)

	servers, serversCtx := errgroup.WithContext(appCtx)

	httpSrv := http.NewAPIServer(
		serversCtx,
		opts,
		http.WithProxyManager(proxyManager),
	)

	servers.Go(func() error {
		if err := httpSrv.Run(); err != nil {
			return errors.New(fmt.Sprintf("HTTP server: %v", err))
		}

		httpSrv.Wait()
		return nil
	})

	if err := servers.Wait(); err != nil {
		log.Printf("[INFO] process terminated, %s", err)
		return
	}
}

func catchForTermination(cancel context.CancelFunc, signals ...os.Signal) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, signals...)
	<-stop
	log.Print("[WARN] interrupt signal")
	cancel()
}
