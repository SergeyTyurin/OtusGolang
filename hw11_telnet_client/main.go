package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "file to read from")
}

func main() {
	flag.Parse()
	withTimeout := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "timeout" {
			withTimeout = true
		}
	})
	addr, err := parseCommandLine(os.Args, withTimeout)
	if err != nil {
		log.Fatal(err)
	}

	client := NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Fprintln(os.Stderr, "...connected to", addr)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		defer cancel()
		err = client.Send()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer cancel()
		err = client.Receive()
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
}

func parseCommandLine(args []string, withTimeout bool) (string, error) {
	count := len(args)
	if (withTimeout && count != 4) || (!withTimeout && count != 3) {
		tmpWithTime := "go-telnet --timeout=10s host port" // #nosec G101
		tmpWithoutTime := "go-telnet host port"            // #nosec G101
		errFormatStr := `args error: require "%v" or "%v"`
		return "", fmt.Errorf(errFormatStr, tmpWithTime, tmpWithoutTime)
	}

	host := args[count-2]
	port := args[count-1]

	return net.JoinHostPort(host, port), nil
}
