package handlers

import (
	"fmt"
	"proc-net-tcp/pkg/tcp"
	"time"
)

const (
	//PROC_NET_TCP = "/proc/net/tcp"
	PROC_NET_TCP      = "/Users/smarin/workspaces/teleport/scratchpad/proc-net-tcp-test"
	READ_TCP_INTERVAL = 60
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
