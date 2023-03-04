package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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

	address := net.JoinHostPort(args[0], args[1])
	tClient := NewTelnetClient(address, ts, os.Stdin, os.Stdout)
	if err := tClient.Connect(); err != nil {
		log.Fatalf("Error connecting %s: %s", address, err.Error())
	}
	defer func() {
		_ = tClient.Close()
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		defer wg.Done()
		for {
			if err := tClient.Receive(); err != nil {
				if !errors.Is(err, io.EOF) {
					_, _ = fmt.Fprint(os.Stderr, err)
				}
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			if err := tClient.Send(); err != nil {
				if !errors.Is(err, io.EOF) {
					_, _ = fmt.Fprint(os.Stderr, err)
				}
				return
			}
		}
	}()

	go func() {
		for {
			sig := <-sigChan
			if sig == syscall.SIGINT {
				fmt.Println("CTRL+C Pressed. Exit...")
				_ = tClient.Close()
				os.Exit(0)
			}
		}
	}()

	wg.Wait()
}
