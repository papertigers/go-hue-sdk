package hue

import (
	"log"
	"net"
	"strings"
	"time"
)

// Hue docs say to use "IpBridge" over "hue-bridgeid"
const _SSDPIdentifier =  "IpBridge"

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

// Discover Hue bridges via SSDP.
// Returns a map of IP.String() to empty struct.
func Discover() map[string]struct{} {
	bridgeSet := make(map[string]struct{})

	rAddr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	_CheckError(err)
	conn, err := net.DialUDP("udp4", nil, rAddr)
	_CheckError(err)
	lAddr, err := net.ListenMulticastUDP("udp4", nil, rAddr)
	_CheckError(err)

	defer conn.Close()
	defer lAddr.Close()

	// Write discovery packet to network
	bytes :=  []byte(strings.Join(_SSDPData, "\r\n"))
	_, err = conn.Write(bytes)
	_CheckError(err)

	// Read responses back for short time period
	timeoutDuration := 30 * time.Second
	lAddr.SetReadDeadline(time.Now().Add(timeoutDuration))

	for {
		buffer := make([]byte, 256)
		n, addr, err := lAddr.ReadFromUDP(buffer)
		if err != nil {
			// TODO log error other than TimeOut
			break
		}
		if strings.Contains(string(buffer[:n]), _SSDPIdentifier) {
			bridgeSet[addr.IP.String()] = struct{}{}
		}
	}

	return bridgeSet
}
