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

var (
	address string
	timeout time.Duration
)

func setup(ctx context.Context, cancel context.CancelFunc) (TelnetClient, error) {
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

func teardown(client TelnetClient, cancel context.CancelFunc) error {
	defer cancel()

	return client.Close()
}

func parseFlags() error {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout on connection")
	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		return errors.New("host must be first command-line argument")
	}
	port := flag.Arg(1)
	if port == "" {
		return errors.New("port must be second command-line argument")
	}
	address = net.JoinHostPort(host, port)

	return nil
}

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	client, err := setup(ctx, cancel)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := teardown(client, cancel); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
