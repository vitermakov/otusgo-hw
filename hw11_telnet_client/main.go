package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	var ts time.Duration
	flag.DurationVar(&ts, "timeout", time.Second*10, "wrong timeout param")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("example: go-telnet --timeout=10s host port")
	}

	port, err := strconv.Atoi(args[1])
	if err != nil || port <= 0 {
		log.Fatalf("invalid port %s", args[1])
	}

	address := fmt.Sprintf("%s:%d", args[0], port)
	tClient := NewTelnetClient(address, ts, os.Stdin, os.Stdout)
	if err := tClient.Connect(); err != nil {
		log.Fatalf("Error connecting %s: %s", address, err.Error())
	}
	defer func() {
		_ = tClient.Close()
	}()

	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer cancel()

	go manageProcess(ctx, cancel, tClient)

	<-ctx.Done()
}

func manageProcess(ctx context.Context, cancel context.CancelFunc, client TelnetClient) {
	go func() {
		<-ctx.Done()
		client.Close()
		fmt.Fprint(os.Stderr, "client closed")
	}()

	go func() {
		for {
			if err := client.Receive(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Fprint(os.Stderr, err)
				}
				cancel()
				return
			}
		}
	}()

	go func() {
		for {
			if err := client.Send(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Fprint(os.Stderr, err)
				}
				cancel()
				return
			}
		}
	}()
}
