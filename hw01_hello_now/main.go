package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const timeFormat = "2006-01-02 15:04:05 -0700 MST"

func main() {
	printCurrentTime()
	printExactTime()
}

func printCurrentTime() {
	now := time.Now()
	fmt.Printf("current time: %s\n", now.Format(timeFormat))
}

func printExactTime() {
	now, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalf("Error from ntp: %v", err)
	}
	fmt.Printf("exact time: %s\n", now.Format(timeFormat))
}
