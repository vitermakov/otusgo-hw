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
	"sync"
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

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// с помощью sigChan оповестим контекст в manageProcess, что надо закрыться.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// с помощью done оповестим main, что в manageProcess все завершилось (возможно с ошибкой).
	done := make(chan struct{})

	go manageProcess(ctx, cancel, done, tClient)

	select {
	case <-sigChan:
		cancel()
		<-done
	case <-done:
		close(sigChan)
	}
}

func manageProcess(ctx context.Context, cancel context.CancelFunc, done chan struct{}, client TelnetClient) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("context receive done")
				return
			default:
				log.Println("before receive")
				if err := client.Receive(); err != nil {
					if !errors.Is(err, io.EOF) {
						log.Printf("error receive: %s\n", err.Error())
					} else {
						log.Println("connection closed")
					}
					cancel()
					// Каким образом внешним сигналом прибить
					os.Stdin.WriteString("\n")
				}
				log.Println("after receive")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Println("context send done")
				return
			default:
				log.Println("before send")
				if err := client.Send(); err != nil {
					if !errors.Is(err, io.EOF) {
						log.Printf("error send: %s\n", err.Error())
					} else {
						log.Println("stdin closed")
					}
					cancel()
				}
				log.Println("after send")
			}
		}
	}()

	wg.Wait()

	done <- struct{}{}
}
