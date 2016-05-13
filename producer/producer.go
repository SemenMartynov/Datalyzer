package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/crackcomm/nsqueue/producer"
)

var (
	nsqdAddr = flag.String("nsqd", "172.18.0.3:4150", "nsqd tcp address")
	amount   = flag.Int("amount", 20, "Amount of messages to produce every 100 ms")
)

func main() {
	flag.Parse()
	producer.Connect(*nsqdAddr)

	for _ = range time.Tick(100 * time.Millisecond) {
		fmt.Println("Ping...")
		for i := 0; i < *amount; i++ {
			body, _ := time.Now().MarshalBinary()
			producer.PublishAsync("latency-test", body, nil)
		}
	}
}
