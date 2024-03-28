package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/config"
	"github.com/henrywhitaker3/adguard-exporter/internal/http"
)

func main() {
	global, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	clients := []*adguard.Client{}
	for _, conf := range global.Configs {
		clients = append(clients, adguard.NewClient(conf))
	}

	for _, client := range clients {
		out, err := client.GetStats(context.Background())
		if err != nil {
			panic(err)
		}
		b, err := json.Marshal(out)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	http := http.NewHttp()
	go http.Serve()

	<-sigs
	if err := http.Stop(ctx); err != nil {
		panic(err)
	}
}
