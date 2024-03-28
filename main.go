package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/config"
)

func main() {
	configs, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	clients := []*adguard.Client{}
	for _, conf := range configs {
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
}
