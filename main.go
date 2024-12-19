package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/config"
	"github.com/henrywhitaker3/adguard-exporter/internal/http"
	"github.com/henrywhitaker3/adguard-exporter/internal/metrics"
	"github.com/henrywhitaker3/adguard-exporter/internal/worker"
)

func main() {
	metrics.Init()
	global, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	clients := []*adguard.Client{}
	for _, conf := range global.Configs {
		clients = append(clients, adguard.NewClient(conf))
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	http := http.NewHttp(global.Server.Debug)
	go http.Serve()
	go worker.Work(ctx, global.Server.Interval, clients)
	http.Ready(true)
	http.Healthy(true)

	<-sigs
	if err := http.Stop(ctx); err != nil {
		panic(err)
	}
}
