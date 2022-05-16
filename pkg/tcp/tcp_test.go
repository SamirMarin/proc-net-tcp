package tcp

import (
	"fmt"
	"proc-net-tcp/pkg/prometheus"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestReadProcNetFileNotFound(t *testing.T) {
	filePath := "testdata/proc-net-tcp-file-no-found"
	_, err := readProcNetTcpFile(filePath)
	expectedError := "open testdata/proc-net-tcp-file-no-found: no such file or directory"

	if expectedError != err.Error() {
		t.Errorf("\n expectedError: %s\n err: %s", expectedError, err.Error())
	}
}

func TestReadProcNetFile(t *testing.T) {
	filePath := "testdata/proc-net-tcp-2"
	fileArr, err := readProcNetTcpFile(filePath)
	if err != nil {
		t.Errorf("error reading file: %v", err)
	}
	expectedFileArr := [][]string{
		[]string{
			"0:",
			"00000000:2382",
			"00000000:0000",
			"0A",
			"00000000:00000000",
			"00:00000000",
			"00000000",
			"0",
			"0",
			"28020",
			"1",
			"ffff88800703e900",
			"100",
			"0",
			"0",
			"10",
			"0",
		},
	}
	if !reflect.DeepEqual(expectedFileArr, fileArr) {
		t.Errorf("\n expectedFileArr: %v\n fileArr: %v", expectedFileArr, fileArr)
	}
}

func TestNewConnectionsNone(t *testing.T) {
	prevTcp := Tcp{
		TimeStamp: time.Now().Add(-10 * time.Second),
		Connections: map[string]ProcNetTcp{
			"172.31.24.7:80-33.191.104.244:56989": ProcNetTcp{
				LocalAdress:  "172.31.24.7",
				LocalPort:    "80",
				RemoteAdress: "33.191.104.244",
				RemotePort:   "56989",
			},
			"172.31.24.7:9000-131.123.105.156:56808": ProcNetTcp{
				LocalAdress:  "172.31.24.7",
				LocalPort:    "9000",
				RemoteAdress: "131.123.105.156",
				RemotePort:   "56808",
			},
		},
	}
	connections := map[string]ProcNetTcp{
		"172.31.24.7:80-33.191.104.244:56989": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "80",
			RemoteAdress: "33.191.104.244",
			RemotePort:   "56989",
		},
	}
	newConnections := prevTcp.NewConnections(connections, time.Now())

	expectedNewConnections := ""
	if expectedNewConnections != newConnections {
		t.Errorf("\n expectedNewConnections: %s\n newConnections: %s", expectedNewConnections, newConnections)
	}
}

func TestNewConnectionsPrevGone(t *testing.T) {
	prevTcp := Tcp{
		TimeStamp: time.Now().Add(-10 * time.Second),
		Connections: map[string]ProcNetTcp{
			"172.31.24.7:80-33.191.104.244:56989": ProcNetTcp{
				LocalAdress:  "172.31.24.7",
				LocalPort:    "80",
				RemoteAdress: "33.191.104.244",
				RemotePort:   "56989",
			},
			"172.31.24.7:9000-131.123.105.156:56808": ProcNetTcp{
				LocalAdress:  "172.31.24.7",
				LocalPort:    "9000",
				RemoteAdress: "131.123.105.156",
				RemotePort:   "56808",
			},
		},
	}
	connections := map[string]ProcNetTcp{
		"172.31.24.7:80-34.191.104.244:56989": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "80",
			RemoteAdress: "34.191.104.244",
			RemotePort:   "56989",
		},
		"172.31.24.7:9000-132.123.105.156:56808": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "9000",
			RemoteAdress: "132.123.105.156",
			RemotePort:   "56808",
		},
	}
	timeNow := time.Now()
	newConnections := prevTcp.NewConnections(connections, timeNow)
	newConnectionsArr := strings.Split(newConnections, "")
	sort.Strings(newConnectionsArr)

	expectedNewConnections := fmt.Sprintf(
		"%v: New connection: 34.191.104.244:56989 -> 172.31.24.7:80\n%v: New connection: 132.123.105.156:56808 -> 172.31.24.7:9000\n",
		timeNow.Format(DATE_FORMAT),
		timeNow.Format(DATE_FORMAT),
	)
	expectedNewConnectionsArr := strings.Split(expectedNewConnections, "")
	sort.Strings(expectedNewConnectionsArr)

	if !reflect.DeepEqual(expectedNewConnectionsArr, newConnectionsArr) {
		t.Errorf("\n newConnections: %s\n expectedNewConnections: %s", expectedNewConnections, newConnections)
	}
}

func TestNewConnectionsNewAdded(t *testing.T) {
	prevTcp := Tcp{
		TimeStamp: time.Now().Add(-10 * time.Second),
		Connections: map[string]ProcNetTcp{
			"172.31.24.7:80-33.191.104.244:56989": ProcNetTcp{
				LocalAdress:  "172.31.24.7",
				LocalPort:    "80",
				RemoteAdress: "33.191.104.244",
				RemotePort:   "56989",
			},
			"172.31.24.7:9000-131.123.105.156:56808": ProcNetTcp{
				LocalAdress:  "172.31.24.7",
				LocalPort:    "9000",
				RemoteAdress: "131.123.105.156",
				RemotePort:   "56808",
			},
		},
	}
	connections := map[string]ProcNetTcp{
		"172.31.24.7:80-34.191.104.244:56989": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "80",
			RemoteAdress: "34.191.104.244",
			RemotePort:   "56989",
		},
		"172.31.24.7:9000-132.123.105.156:56808": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "9000",
			RemoteAdress: "132.123.105.156",
			RemotePort:   "56808",
		},
		"172.31.24.7:80-33.191.104.244:56989": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "80",
			RemoteAdress: "33.191.104.244",
			RemotePort:   "56989",
		},
		"172.31.24.7:9000-131.123.105.156:56808": ProcNetTcp{
			LocalAdress:  "172.31.24.7",
			LocalPort:    "9000",
			RemoteAdress: "131.123.105.156",
			RemotePort:   "56808",
		},
	}

	timeNow := time.Now()
	newConnections := prevTcp.NewConnections(connections, timeNow)
	newConnectionsArr := strings.Split(newConnections, "")
	sort.Strings(newConnectionsArr)

	expectedNewConnections := fmt.Sprintf(
		"%v: New connection: 34.191.104.244:56989 -> 172.31.24.7:80\n%v: New connection: 132.123.105.156:56808 -> 172.31.24.7:9000\n",
		timeNow.Format(DATE_FORMAT),
		timeNow.Format(DATE_FORMAT),
	)
	expectedNewConnectionsArr := strings.Split(expectedNewConnections, "")
	sort.Strings(expectedNewConnectionsArr)

	if !reflect.DeepEqual(expectedNewConnectionsArr, newConnectionsArr) {
		t.Errorf("\n expectedNewConnections: %s\n newConnections: %s", expectedNewConnections, newConnections)
	}
}

func TestNewTcpNoPortScan(t *testing.T) {
	prevTcp := Tcp{}
	_, _, portScanStr, _ := prevTcp.NewTcp("testdata/proc-net-tcp-2")
	expectedPortScanStr := ""
	if expectedPortScanStr != portScanStr {
		t.Errorf("\n expectedPortScanStr: %s\n portScanStr: %s", expectedPortScanStr, portScanStr)
	}
}

func TestNewTcpPortScan(t *testing.T) {
	prevTcp := Tcp{
		PromClient: prometheus.NewPromClient(
			"proc_net_tcp_new_connections_total_2",
			"The total number of new connections_2",
		),
		SkipIpBlock: true,
	}
	tcp, _, portScanStr, _ := prevTcp.NewTcp("testdata/proc-net-tcp")
	portScanArr := strings.Split(portScanStr, "")
	sort.Strings(portScanArr)

	expectedPortScanStr := fmt.Sprintf("%v: Port scan detected: 216.19.179.173 -> 172.31.24.7 on ports 23,24,22\n", tcp.TimeStamp.Format(DATE_FORMAT))
	expectedPortScanArr := strings.Split(expectedPortScanStr, "")
	sort.Strings(expectedPortScanArr)

	if !reflect.DeepEqual(expectedPortScanArr, portScanArr) {
		t.Errorf("\n expectedPortScanStr: %s\n portScanStr: %s", expectedPortScanStr, portScanStr)
	}
}

func TestNewTcpPortScanInPrevMin(t *testing.T) {
	prevTcp := Tcp{
		TimeStamp: time.Now().Add(-10 * time.Second),
		PortScans: map[string]PortScan{
			"172.31.24.7-216.19.179.180": PortScan{
				LocalIp:  "172.31.24.7",
				SourceIP: "216.19.179.180",
				Ports: map[string]Port{
					"7777": Port{
						TimeStamp: time.Now().Add(-30 * time.Second),
						Port:      "7777",
					},
					"7778": Port{
						TimeStamp: time.Now().Add(-40 * time.Second),
						Port:      "7778",
					},
					"7378": Port{
						TimeStamp: time.Now().Add(-20 * time.Second),
						Port:      "7378",
					},
				},
			},
		},
		PromClient: prometheus.NewPromClient(
			"proc_net_tcp_new_connections_total_3",
			"The total number of new connections_3",
		),
		SkipIpBlock: true,
	}
	tcp, _, portScanStr, _ := prevTcp.NewTcp("testdata/proc-net-tcp")
	portScanArr := strings.Split(portScanStr, "")
	sort.Strings(portScanArr)

	expectedPortScanStr := fmt.Sprintf("%v: Port scan detected: 216.19.179.180 -> 172.31.24.7 on ports 7777,7778,7378\n%v: Port scan detected: 216.19.179.173 -> 172.31.24.7 on ports 23,24,22\n", tcp.TimeStamp.Format(DATE_FORMAT), tcp.TimeStamp.Format(DATE_FORMAT))
	expectedPortScanArr := strings.Split(expectedPortScanStr, "")
	sort.Strings(expectedPortScanArr)

	if !reflect.DeepEqual(expectedPortScanArr, portScanArr) {
		t.Errorf("\n expectedPortScanStr: %s\n portScanStr: %s", expectedPortScanStr, portScanStr)
	}
}

func TestNewTcpPortScanRemoveLongerThanMin(t *testing.T) {
	prevTcp := Tcp{
		TimeStamp: time.Now().Add(-10 * time.Second),
		PortScans: map[string]PortScan{
			"172.31.24.7-216.19.179.180": PortScan{
				LocalIp:  "172.31.24.7",
				SourceIP: "216.19.179.180",
				Ports: map[string]Port{
					"7777": Port{
						TimeStamp: time.Now().Add(-30 * time.Second),
						Port:      "7777",
					},
					"7778": Port{
						TimeStamp: time.Now().Add(-70 * time.Second),
						Port:      "7778",
					},
					"7378": Port{
						TimeStamp: time.Now().Add(-20 * time.Second),
						Port:      "7378",
					},
				},
			},
		},
		PromClient: prometheus.NewPromClient(
			"proc_net_tcp_new_connections_total_4",
			"The total number of new connections_4",
		),
		SkipIpBlock: true,
	}
	tcp, _, portScanStr, _ := prevTcp.NewTcp("testdata/proc-net-tcp")
	portScanArr := strings.Split(portScanStr, "")
	sort.Strings(portScanArr)

	expectedPortScanStr := fmt.Sprintf("%v: Port scan detected: 216.19.179.173 -> 172.31.24.7 on ports 23,24,22\n", tcp.TimeStamp.Format(DATE_FORMAT))
	expectedPortScanArr := strings.Split(expectedPortScanStr, "")
	sort.Strings(expectedPortScanArr)

	if !reflect.DeepEqual(expectedPortScanArr, portScanArr) {
		t.Errorf("\n expectedPortScanStr: %s\n portScanStr: %s", expectedPortScanStr, portScanStr)
	}
}
