package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	closed  bool
}

func (c *client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect: %w", err)
	}
	log.Printf("...Connected to %s\n", c.address)

	return nil
}

func (c *client) Close() error {
	c.closed = true
	return c.conn.Close()
}

func (c *client) Send() error {
	r := bufio.NewReader(c.in)
	for {
		str, err := r.ReadString('\n')
		if errors.Is(err, io.EOF) {
			log.Println("...EOF")
			return nil
		}
		if err != nil {
			return c.formatSendError(err)
		}

		_, err = c.conn.Write([]byte(str))
		if err != nil {
			return c.formatSendError(err)
		}
	}
}

func (c *client) Receive() error {
	r := bufio.NewReader(c.conn)
	for {
		str, err := r.ReadString('\n')
		if errors.Is(err, io.EOF) {
			log.Println("...Connection was closed by peer")
			return nil
		}
		if err != nil {
			return c.formatReceiveError(err)
		}

		_, err = c.out.Write([]byte(str))
		if err != nil {
			return c.formatReceiveError(err)
		}
	}
}

func (c *client) formatSendError(err error) error {
	if c.closed {
		return nil
	}
	return fmt.Errorf("cannot send: %w", err)
}

func (c *client) formatReceiveError(err error) error {
	if c.closed {
		return nil
	}
	return fmt.Errorf("cannot receive: %w", err)
}
