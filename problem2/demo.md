## 1. First compile drop_tcp_pack_of_proc.c file
```
    clang -O2 -target bpf -c drop_tcp_pack_of_proc.c -o drop_tcp_pack_of_proc.o
```

## 2. Create a C group
```
    sudo mkdir -p /sys/fs/cgroup/myprocess
```

## 3. Start our dummy process (myprocess)
```
    ./myProcess/myprocess
```

## 4. Run go program
```
    sudo go run main.go myprocess 4040
```

## 5. Start two listeners at port 4040 and 8080 in sperate terminals
```
    # Allowed port
    nc -l -k 4040   # messages from myprocess pass

    # Blocked port
    nc -l -k 8080   # messages are dropped
```

## Resources Used
* cilium ebpf this [example](https://github.com/cilium/ebpf/blob/main/examples/cgroup_skb/cgroup_skb.c)
* ebpf documentation about [cgroup_skb/egress](https://docs.ebpf.io/linux/program-type/BPF_PROG_TYPE_CGROUP_SKB/)