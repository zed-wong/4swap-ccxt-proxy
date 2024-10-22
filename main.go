package main

import (
	"context"
	fswap "github.com/fox-one/4swap-sdk-go/v2"
)

const (
	MTG_GROUP = "group"
	FSWAP_CLIENT = "client"
	HOST_KEY = "host"
	PORT_KEY = "port"

	HOST = "0.0.0.0"
	PORT = 80
)

func main() {
	client := fswap.New()
	ctx := context.Background()
	group, err := client.ReadGroup(ctx)
	if err != nil {
		panic(err)
	}

	ctx = context.WithValue(ctx, HOST_KEY, HOST)
	ctx = context.WithValue(ctx, PORT_KEY, PORT)
	ctx = context.WithValue(ctx, MTG_GROUP, group)
	ctx = context.WithValue(ctx, FSWAP_CLIENT, client)
	StartAPIServer(ctx)
}
