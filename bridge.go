package hue

import (
	"bytes"
	"net"
	"strings"
	"time"
)

// Hue docs say to use "IpBridge" over "hue-bridgeid"
const _SSDPIdentifier = "IpBridge"

var _SSDPData = []string{
	"M-SEARCH * HTTP/1.1",
	"HOST:239.255.255.250:1900",
	"MAN:\"ssdp:discover\"",
	"ST:ssdp:all",
	"MX:1",
}

type Bridge struct {
	ip net.IP
}

// Discover Hue bridges via SSDP.
// Returns a map of IP.String() to empty struct.
func Discover() (map[string]struct{}, error) {
	bridgeSet := make(map[string]struct{})

	rAddr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, rAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	lAddr, err := net.ListenMulticastUDP("udp4", nil, rAddr)
	if err != nil {
		return nil, err
	}
	defer lAddr.Close()

	// Write discovery packet to network
	if _, err = conn.Write([]byte(strings.Join(_SSDPData, "\r\n"))); err != nil {
		return nil, err
	}

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
		if bytes.Contains(buffer[:n], []byte(_SSDPIdentifier)) {
			bridgeSet[addr.IP.String()] = struct{}{}
		}
	}

	return bridgeSet, nil
}
