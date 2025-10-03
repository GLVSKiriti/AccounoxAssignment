//go:build ignore_c_file
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>
#include <linux/sched.h>
#include <linux/pkt_cls.h>

// using a bpf map array for which stores allowed port of myprocess
struct bpf_map_def SEC("maps") config_map = {
    .type = BPF_MAP_TYPE_ARRAY,
    .key_size = sizeof(int),
    .value_size = sizeof(__u16),
    .max_entries = 1,
};

// cgroup_skb egress program
SEC("cgroup_skb/egress")
int drop_tcp_pack_of_proc(struct __sk_buff *skb)
{
    // Load IP header safely
    __u8 ip_header[20];
    if (bpf_skb_load_bytes(skb, 0, ip_header, sizeof(ip_header)) < 0)
        return 0;

    struct iphdr *ip = (struct iphdr *)ip_header;

    if (ip->protocol != IPPROTO_TCP)
        return 1; // allow non-TCP

    // Load TCP header safely
    __u8 tcp_header[20];
    if (bpf_skb_load_bytes(skb, ip->ihl * 4, tcp_header, sizeof(tcp_header)) < 0)
        return 0;

    struct tcphdr *tcp = (struct tcphdr *)tcp_header;

    // Load allowed port from map
    int key = 0;
    __u16 *allowed_port = bpf_map_lookup_elem(&config_map, &key);
    if (!allowed_port)
        return 0;

    // Drop traffic other than allowed port
    if ((__builtin_bswap16(tcp->dest)) != *allowed_port)
        return 0;

    return 1;
}

char _license[] SEC("license") = "GPL";