package hue

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var _SSDPData = []string {
	"M-SEARCH * HTTP/1.1",
	"HOST:239.255.255.250:1900",
	"MAN:\"ssdp:discover\"",
	"ST:ssdp:all",
	"MX:1",
}

type Bridge struct {
	ip net.IP
}

func _CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Discover Hue bridges via SSDP
func Discover() []net.IP {
	b := []net.IP{}

	rAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	_CheckError(err)
	lAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	_CheckError(err)
	conn, err := net.DialUDP("udp", lAddr, rAddr)
	_CheckError(err)

	defer conn.Close()

	// Write discovery packet to network
	bytes :=  []byte(strings.Join(_SSDPData, "\r\n"))
	_, err = conn.Write(bytes)
	_CheckError(err)

	// Read responses back for 5seconds
	timeoutDuration := 5 * time.Second
	buffer := make([]byte, 1024)
	for {
		fmt.Println("Entering loop")
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err);
			break
		}
		fmt.Println("Received from : ", addr)
		fmt.Println(string(buffer[:n]))
	}

	return b
}
