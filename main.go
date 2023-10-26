package main

import (
	"blackbox-v2/api"
	"blackbox-v2/pkg/mongoc"
	"context"
	"os"
	"sync"
)

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	//  testing git commit on mac
	args := os.Args[1:]
	if len(args) == 0 {
		panic("No arguments provided")
	}
	if args[0] == "setup" {
		setup()
	} else if args[0] == "server" {
		api.RunServer()
		// HandleDisconnections(ctx)
	} else {
		panic("Invalid argument")
	}
}

func HandleDisconnections(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		mongoc.HandleMongoConn(ctx)
		wg.Done()
	}()
	wg.Wait()
}
