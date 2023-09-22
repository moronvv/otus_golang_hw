package main

import (
	"io"
	"net"
	"time"
)

const (
	buffSize = 1024
	network  = "tcp"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	conn net.Conn

	in      io.ReadCloser
	out     io.Writer
	address string
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *client) Connect() error {
	var err error

	c.conn, err = net.DialTimeout(network, c.address, c.timeout)
	return err
}

func (c *client) Close() error {
	if err := c.in.Close(); err != nil {
		return err
	}

	return c.conn.Close()
}

func (c *client) Send() error {
	buffer := make([]byte, buffSize)

	n, err := c.in.Read(buffer)
	if err != nil {
		return err
	}

	_, err = c.conn.Write(buffer[:n])
	return err
}

func (c *client) Receive() error {
	buffer := make([]byte, buffSize)

	n, err := c.conn.Read(buffer)
	if err != nil {
		return err
	}

	_, err = c.out.Write(buffer[:n])
	return err
}
