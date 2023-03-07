package server

import (
	"context"
	"errors"
	"fmt"
	httpCli "github.com/akhmettolegen/proxy/internal/clients/http"
	"github.com/akhmettolegen/proxy/internal/database"
	"github.com/akhmettolegen/proxy/internal/database/drivers"
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

	// Database
	ds, err := setupDatabase(opts)
	if err != nil {
		log.Println(err)
		return
	}
	defer ds.Close()

	client := httpCli.NewClient(appCtx)

	proxyManager := proxy.NewManager(appCtx, client, ds.Task())

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

func setupDatabase(opts *configs.Config) (drivers.DataStore, error) {
	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataBaseName:  opts.DSDB,
		DataStoreName: opts.DSName,
	})
	if err != nil {
		return nil, err
	}

	if err := ds.Connect(); err != nil {
		errText := fmt.Sprintf("[ERROR] cannot connect to datastore %s: %v", opts.DSName, err)
		return nil, errors.New(errText)
	}

	fmt.Println("Connected to", ds.Name())

	return ds, nil
}
