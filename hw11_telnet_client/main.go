package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setup(ctx context.Context, cancel context.CancelFunc, address string, timeout time.Duration) (TelnetClient, error) {
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	// goroutine for sending messages to server
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Send(); err != nil {
					if errors.Is(err, io.EOF) {
						fmt.Fprintln(os.Stderr, "...EOF")
						cancel()
					}
					return
				}
			}
		}
	}()

	// goroutine for receiving messages from server
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Receive(); err != nil {
					if errors.Is(err, io.EOF) {
						fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
						cancel()
					}
					return
				}
			}
		}
	}()

	return client, nil
}

func teardown(client TelnetClient) error {
	if err := client.Close(); err != nil {
		return err
	}

	return nil
}

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout on connection")
	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		log.Fatal("host must be first command-line argument")
	}

	port := flag.Arg(1)
	if port == "" {
		log.Fatal("port must be second command-line argument")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client, err := setup(ctx, cancel, net.JoinHostPort(host, port), timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer teardown(client)

	<-ctx.Done()
}
