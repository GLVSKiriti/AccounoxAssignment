package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

// Gets process pid using given process name
func getPID(procName string) (int, error) {
	cmd := exec.Command("pidof", procName)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	pids := strings.Fields(string(bytes.TrimSpace(out)))
	if len(pids) == 0 {
		return 0, fmt.Errorf("process %s not running", procName)
	}
	pid, err := strconv.Atoi(pids[0])
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// Move the process into a cgroup
func moveToCgroup(cgroupPath string, pid int) error {
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo %d > %s/cgroup.procs", pid, cgroupPath))
	return cmd.Run()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sudo ./main.go <process_name> <port>")
		return
	}

	procName := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	cgroupPath := "/sys/fs/cgroup/" + procName

	if err := rlimit.RemoveMemlock(); err != nil {
		panic(err)
	}

	// Find PID
	pid, err := getPID(procName)
	if err != nil {
		panic(err)
	}

	// Move process into cgroup
	if err := moveToCgroup(cgroupPath, pid); err != nil {
		panic(err)
	}

	fmt.Println("Process", procName, "with PID", pid, "moved to cgroup", cgroupPath)

	// Load eBPF program
	spec, err := ebpf.LoadCollectionSpec("drop_tcp_pack_of_proc.o")
	if err != nil {
		panic(err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		panic(err)
	}
	defer coll.Close()

	prog := coll.Programs["drop_tcp_pack_of_proc"]
	cfg := coll.Maps["config_map"]

	// Set allowed port
	key := uint32(0)
	allowedPort := uint16(port)
	if err := cfg.Put(key, allowedPort); err != nil {
		panic(err)
	}

	// Attach to cgroup egress
	lnk, err := link.AttachCgroup(link.CgroupOptions{
		Path:    cgroupPath,
		Attach:  ebpf.AttachCGroupInetEgress,
		Program: prog,
	})
	if err != nil {
		panic(err)
	}
	defer lnk.Close()

	fmt.Println("eBPF egress loaded on cgroup:", cgroupPath, "allowed port:", port)

	// Keep program alive
	for {
		time.Sleep(time.Hour)
	}
}
