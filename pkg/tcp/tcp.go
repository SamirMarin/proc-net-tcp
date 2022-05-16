package tcp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"proc-net-tcp/pkg/iptables"
	"proc-net-tcp/pkg/prometheus"
	"strings"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02 15:04:05"
)

type Tcp struct {
	TimeStamp   time.Time
	Connections map[string]ProcNetTcp
	PortScans   map[string]PortScan
	PromClient  *prometheus.PromClient
}

type PortScan struct {
	LocalIp  string
	SourceIP string
	Ports    map[string]Port
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
func (t *Tcp) NewTcp(filepath string) (*Tcp, string, string, error) {
	connectionsRead, err := readProcNetTcpFile(filepath)
	if err != nil {
		return nil, "", "", err
	}
	timeStamp := time.Now()

	//set up promClient if not set
	if t.PromClient == nil {
		t.PromClient = prometheus.NewPromClient(
			"proc_net_tcp_new_connections_total",
			"The total number of new connections",
		)
	}

	connections := map[string]ProcNetTcp{}
	portScans := map[string]PortScan{}
	if len(t.PortScans) > 0 {
		portScans = t.PortScans
	}

	for _, connection := range connectionsRead {
		ipLoc, portLoc := hexToStringIPPort(connection[1])
		ipRem, portRem := hexToStringIPPort(connection[2])
		connectionState := hexToDec(connection[3])
		key := fmt.Sprintf("%s:%s-%s:%s", ipLoc, portLoc, ipRem, portRem)

		//considering only established connections
		if connectionState == "1" {
			connections[key] = ProcNetTcp{
				LocalAdress:  ipLoc,
				LocalPort:    portLoc,
				RemoteAdress: ipRem,
				RemotePort:   portRem,
			}

			portScanKey := fmt.Sprintf("%s-%s", ipLoc, ipRem)
			portScan, ok := portScans[portScanKey]

			if !ok {
				portScans[portScanKey] = PortScan{
					LocalIp:  ipLoc,
					SourceIP: ipRem,
					Ports: map[string]Port{
						portLoc: Port{
							TimeStamp: timeStamp,
							Port:      portLoc,
						},
					},
				}
			} else {
				portScan.Ports[portLoc] = Port{
					TimeStamp: timeStamp,
					Port:      portLoc,
				}
				portScans[portScanKey] = portScan
			}
		}
	}

	newConnections := t.NewConnections(connections, timeStamp)

	portScans = cleanOldPorts(60*time.Second, timeStamp, portScans)
	currentPortScans := findPortScans(portScans, timeStamp)

	return &Tcp{
		TimeStamp:   timeStamp,
		Connections: connections,
		PortScans:   portScans,
		PromClient:  t.PromClient,
	}, newConnections, currentPortScans, nil
}

// Finds new connections by comparing previous read values with new values
func (t *Tcp) NewConnections(connections map[string]ProcNetTcp, timeStamp time.Time) string {
	var conns bytes.Buffer
	for key, newConn := range connections {
		_, ok := t.Connections[key]
		if !ok {
			conn := fmt.Sprintf(
				"%v: New connection: %s:%s -> %s:%s\n",
				timeStamp.Format(DATE_FORMAT),
				newConn.RemoteAdress,
				newConn.RemotePort,
				newConn.LocalAdress,
				newConn.LocalPort,
			)
			conns.WriteString(conn)
			t.PromClient.IncToCounter()
		}
	}
	return conns.String()
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

// single source IP connects to more than 3 host ports
func findPortScans(portScans map[string]PortScan, time time.Time) string {
	var scans bytes.Buffer
	for _, portScan := range portScans {
		if len(portScan.Ports) >= 3 {
			preFix := fmt.Sprintf("%s: Port scan detected: %s -> %s on ports ", time.Format(DATE_FORMAT), portScan.SourceIP, portScan.LocalIp)
			for _, port := range portScan.Ports {
				portStr := fmt.Sprintf("%s%s", preFix, port.Port)
				scans.WriteString(portStr)
				preFix = ","
			}
			scans.WriteString("\n")
			blockIp(portScan.SourceIP)
		}
	}
	return scans.String()
}

func blockIp(ip string) {
	ipt, err := iptables.NewIptables()
	if err != nil {
		errorMsg := fmt.Sprintf("Unable to create iptables client, error: %v", err)
		fmt.Println(errorMsg)
	} else {
		err = ipt.BlockSourceIp(ip)
		if err != nil {
			errorMsg := fmt.Sprintf("Unable to block IP: %s, error: %v", ip, err)
			fmt.Println(errorMsg)
		}
		fmt.Println(fmt.Sprintf("IP: %s Blocked", ip))
	}
}

// Cleans out port connections that happen over a minute ago
func cleanOldPorts(timeTresh time.Duration, time time.Time, portScans map[string]PortScan) map[string]PortScan {
	for _, portScan := range portScans {
		for key, port := range portScan.Ports {
			timeExpire := time.Add(-timeTresh)
			if port.TimeStamp.Before(timeExpire) {
				delete(portScan.Ports, key)
			}
		}
	}
	return portScans
}
