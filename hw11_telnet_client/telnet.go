package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClientWithTimeout struct {
	timeout time.Duration
	address string
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClientWithTimeout) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	t.conn = conn
	return err
}

func (t *telnetClientWithTimeout) Close() error {
	return t.conn.Close()
}

func (t *telnetClientWithTimeout) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telnetClientWithTimeout) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	cl := telnetClientWithTimeout{timeout, address, in, out, nil}
	return &cl
}
