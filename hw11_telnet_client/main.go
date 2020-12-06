package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(0)

	timeout := flag.Duration("timeout", 10*time.Second, "server connect timeout")
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("Please define address and port")
	}
	address := flag.Arg(0) + ":" + flag.Arg(1)

	ctx, cancel := context.WithCancel(context.Background())

	go watchSignals(cancel)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	go send(client, cancel)
	go receive(client, cancel)

	<-ctx.Done()
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

func send(client TelnetClient, cancel context.CancelFunc) {
	err := client.Send()
	if err != nil {
		log.Println(err)
	}
	cancel()
}

func receive(client TelnetClient, cancel context.CancelFunc) {
	err := client.Receive()
	if err != nil {
		log.Println(err)
	}
	cancel()
}
