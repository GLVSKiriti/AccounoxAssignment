//go:build ignore_c_file
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/in.h> 
#include <linux/ip.h>
#include <linux/tcp.h>

// This is a bpf map which is an array having one entry i.e, port number
// required to configure the port from userspace
struct bpf_map_def SEC("maps") target_port_map = {
    .type = BPF_MAP_TYPE_ARRAY,
    .key_size = sizeof(int),
    .value_size = sizeof(int),
    .max_entries = 1,
};

// A TCP packet generally has etherenet header -> ip header -> tcp header -> payload etc..
SEC("xdp")
int drop_tcp_port(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    // check wether ipv4 packet or not else dont drop it
    struct ethhdr *eth = data; 
    if (eth + 1 > data_end || eth->h_proto != __constant_htons(ETH_P_IP)) return XDP_PASS;

    // if not tcp packet dont drop it
    struct iphdr *ip = data + sizeof(*eth);
    if (ip + 1 > data_end || ip->protocol != IPPROTO_TCP) return XDP_PASS;

    struct tcphdr *tcp = (void *)ip + ip->ihl*4;
    if (tcp + 1 > data_end) return XDP_PASS;

    // look up the user configured port from bpf map
    int key = 0;
    int *target_port = bpf_map_lookup_elem(&target_port_map, &key);
    if (!target_port) return XDP_PASS;

    if (tcp->dest == __constant_htons(*target_port))
        return XDP_DROP; // drop the packet

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
