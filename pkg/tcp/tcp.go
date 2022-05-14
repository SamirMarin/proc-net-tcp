package tcp

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"time"
)

type Tcp struct {
	TimeStamp   time.Time
	Connections map[string]ProcNetTcp
	PortScan    PortScan
}

type PortScan struct {
	sourceIP string
	port     []Port
}

type Port struct {
	TimeStamp time.Time
	Port      string
}

type ProcNetTcp struct {
	LocalAdress  string
	LocalPort    string
	RemoteAdress string
	RemotePort   string
}

// returns a new tcp connection
func Connections(filepath string) (*Tcp, error) {
	connectionsRead, err := readProcNetTcpFile(filepath)
	if err != nil {
		return nil, err
	}

	connections := map[string]ProcNetTcp{}
	for _, connection := range connectionsRead {
		ipLoc, portLoc := hexToStringIPPort(connection[1])
		ipRem, portRem := hexToStringIPPort(connection[2])
		key := fmt.Sprintf("%s:%s-%s:%s", ipLoc, portLoc, ipRem, portRem)

		connections[key] = ProcNetTcp{
			LocalAdress:  ipLoc,
			LocalPort:    portLoc,
			RemoteAdress: ipRem,
			RemotePort:   portRem,
		}
	}
	fmt.Println(connections)

	return &Tcp{
		Connections: connections,
	}, nil
}

// Reads /proc/net/tcp file
func readProcNetTcpFile(path string) ([][]string, error) {
	tcpConncetionsFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	entries := strings.Split(string(tcpConncetionsFile), "\n")[1:]

	entriesClean := [][]string{}
	for _, entry := range entries {
		entryTrimed := strings.TrimSpace(entry)
		entryArr := strings.Split(entryTrimed, " ")

		entryClean := []string{}
		for _, entryVal := range entryArr {
			if entryVal != "" {
				entryClean = append(entryClean, entryVal)
			}
		}
		if len(entryClean) > 0 {
			entriesClean = append(entriesClean, entryClean)
		}
	}
	return entriesClean, nil
}

// Get ip, port from hex: ip:port
func hexToStringIPPort(ipPort string) (string, string) {
	ipPortArr := strings.Split(ipPort, ":")
	ip := ipHexDecStr(ipPortArr[0])
	port := hexToDec(ipPortArr[1])

	return ip, port
}

// Converts hex to dec
func hexToDec(hex string) string {
	n := new(big.Int)
	n.SetString(hex, 16)

	return n.String()
}

// Get IP in decimal from hex given in little endian
func ipHexDecStr(ip string) string {
	ipStr := ""
	for i, _ := range ip {
		if i%2 == 1 {
			ipStr = hexToDec(ip[i-1:i+1]) + ipStr
		} else if i != 0 {
			ipStr = "." + ipStr
		}
	}
	return ipStr
}
