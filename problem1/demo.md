1. compile drop_tcp_port.c

clang -O2 -target bpf -c drop_tcp_port.c -o drop_tcp_port.o

2. settup testing environment

    1. setup vethpair
        sudo ip link add veth0 type veth peer name veth1
        sudo ip link set veth0 up
        sudo ip link set veth1 up
        sudo ip addr add 10.0.0.1/24 dev veth0
        sudo ip addr add 10.0.0.2/24 dev veth1

    2. Create network namespace
        sudo ip netns add ns1
        sudo ip link set veth0 netns ns1
        sudo ip netns exec ns1 ip link set veth0 up
        sudo ip netns exec ns1 ip addr add 10.0.0.1/24 dev veth0

3. Now veth0 is inside ns1 namespace and veth1 on host 

4. Now run go code
    sudo go run main.go veth1 4040
    which can now drop all tcp packets at port 4040 through veth1 interface

5. Send Test TCP packets from veth0
    sudo ip netns exec ns1 hping3 -S 10.0.0.2 -p 4040 -c 5

6. cleanup 
    sudo ip netns delete ns1
    sudo ip link delete veth1 type veth

