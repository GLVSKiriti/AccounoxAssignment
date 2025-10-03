package main

import (
	"fmt"
	"net"
	"time"
)

// Simple standard code that sends tcp packets i.e, connections to give ports
func sendPackets(addr string) {
	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			// Connection might fail if server isn't listening; retry after delay
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Fprintf(conn, "hello from myprocess\n")
		conn.Close()
		time.Sleep(500 * time.Millisecond) // control rate
	}
}

func main() {
	go sendPackets("127.0.0.1:4040")
	go sendPackets("127.0.0.1:8080")

	// keep process running
	select {}
}
