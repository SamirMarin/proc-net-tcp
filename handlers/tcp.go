package handlers

import (
	"fmt"
	"time"

	"github.com/SamirMarin/proc-net-tcp/pkg/tcp"
)

const (
	PROC_NET_TCP      = "/proc/net/tcp"
	READ_TCP_INTERVAL = 10
)

func Tcp() {
	currentTcp := &tcp.Tcp{}
	for {
		fmt.Println("Finding new connections")

		tcp, connInfo, portScanInfo, err := currentTcp.NewTcp(PROC_NET_TCP)
		if err != nil {
			errMsg := fmt.Sprintf("error found getting new TCP connections: %x", err.Error())
			fmt.Println(errMsg)
		}
		currentTcp = tcp
		fmt.Print(connInfo)
		fmt.Print(portScanInfo)
		time.Sleep(READ_TCP_INTERVAL * time.Second)
	}
}
