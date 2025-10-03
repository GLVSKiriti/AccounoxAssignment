package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <iface> <port>")
		return
	}

	iface := os.Args[1]
	port := os.Args[2]

	// Convert port to int
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("invalid port: %v", err)
	}

	// Load precompiled eBPF object
	spec, err := ebpf.LoadCollectionSpec("drop_tcp_port.o")
	if err != nil {
		log.Fatalf("failed to load obj: %v", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		log.Fatalf("failed to create collection: %v", err)
	}
	defer coll.Close()

	prog := coll.Programs["drop_tcp_port"]

	// Attach XDP to interface
	l, err := link.AttachXDP(link.XDPOptions{
		Program:   prog,
		Interface: ifaceIndex(iface),
	})
	if err != nil {
		log.Fatalf("failed to attach XDP: %v", err)
	}
	defer l.Close()

	// Set port in map
	m := coll.Maps["target_port_map"]
	var key uint32 = 0
	var po uint32 = uint32(p)
	if err := m.Put(key, po); err != nil {
		log.Fatalf("failed to set port: %v", err)
	}

	fmt.Printf("Dropping TCP packets on port %d on interface %s\n", p, iface)
	select {} // wait forever
}

func ifaceIndex(name string) int {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Fatalf("interface not found: %v", err)
	}
	return iface.Index
}
