package main

import (
    "bufio"
    "encoding/hex"
    "log"
    "net"
    "os"
    "regexp"
    "strings"
    "strconv"
)

func validate_ip(ip string) (found bool) {
    // regex from https://stackoverflow.com/questions/13145397/regex-for-multicast-ip-address
    ip_regex := regexp.MustCompile(`^2(?:2[4-9]|3\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]\d?|0)){3}:[0-9]+$`)
    return ip_regex.MatchString(ip)
}

func validate_input(input string) {
    input = strings.TrimSuffix(input, "\n")
    input_split := strings.Split(input, " ")
    if (!validate_ip(input_split[0])) {
        log.Fatal("Ip %s didn't pass regex", input_split[0])
        return
    }
    if (!validate_ip(input_split[1])) {
        log.Fatal("Ip %s didn't pass regex", input_split[1])
        return
    }
    port, err := strconv.Atoi(input_split[2])
    if err != nil {
        log.Fatal(err)
        return
    }
    go forward(input_split[0], input_split[1], port)
}

func forward(addr_to_listen string, addr_to_send string, datagram_size int) {
    addr_listen, err := net.ResolveUDPAddr("udp", addr_to_listen)
	if err != nil {
		log.Fatal(err)
    }
    udp, err := net.ListenMulticastUDP("udp", nil, addr_listen)
	if err != nil {
		log.Fatal(err)
    }
    addr_send, err := net.ResolveUDPAddr("udp", addr_to_send)
    if err != nil {
        log.Fatal(err)
    }
    sender, err := net.DialUDP("udp", nil, addr_send)
    if err != nil {
        log.Fatal(err)
    }

    udp.SetReadBuffer(datagram_size)
    for {
        data_arr := make([]byte, datagram_size)
        data_size, src, err := udp.ReadFromUDP(data_arr)
        if err != nil {
            log.Fatal("Failed:", err)
        }
        log.Println("Received from:", src)
        log.Println(hex.Dump(data_arr[:data_size]))
        sender.Write(data_arr[:data_size])
    }
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        input, _ := reader.ReadString('\n')
        validate_input(input)
    }
}
