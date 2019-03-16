package main

import (
    "bufio"
    "encoding/hex"
    "log"
    "net"
    "os"
    "strings"
    "strconv"
)

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
        input_split := strings.Split(input, " ")
        port, _ := strconv.Atoi(input_split[2])
        go forward(input_split[0], input_split[1], port)
    }
}
