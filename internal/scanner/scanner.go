package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortInfo holds information about an open port.
type PortInfo struct {
	Protocol string
	Port     int
	Address  string
}

// String returns a human-readable representation of a PortInfo.
func (p PortInfo) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Scanner defines the interface for port scanning.
type Scanner interface {
	Scan() ([]PortInfo, error)
}

// TCPScanner scans for open TCP ports in a given range.
type TCPScanner struct {
	Host      string
	PortStart int
	PortEnd   int
}

// NewTCPScanner creates a new TCPScanner with the given parameters.
func NewTCPScanner(host string, portStart, portEnd int) *TCPScanner {
	return &TCPScanner{
		Host:      host,
		PortStart: portStart,
		PortEnd:   portEnd,
	}
}

// Scan checks each port in the range and returns those that are open.
func (s *TCPScanner) Scan() ([]PortInfo, error) {
	var openPorts []PortInfo

	for port := s.PortStart; port <= s.PortEnd; port++ {
		address := net.JoinHostPort(s.Host, strconv.Itoa(port))
		conn, err := net.Dial("tcp", address)
		if err != nil {
			if isConnectionRefused(err) {
				continue
			}
			continue
		}
		conn.Close()
		openPorts = append(openPorts, PortInfo{
			Protocol: "tcp",
			Port:     port,
			Address:  s.Host,
		})
	}

	return openPorts, nil
}

// isConnectionRefused checks if the error is a connection refused error.
func isConnectionRefused(err error) bool {
	return strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no connection could be made")
}
