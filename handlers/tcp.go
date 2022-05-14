package handlers

import (
	"fmt"
	"proc-net-tcp/pkg/tcp"
	"time"
)

const (
	//PROC_NET_TCP = "/proc/net/tcp"
	PROC_NET_TCP = "/Users/smarin/workspaces/teleport/scratchpad/proc-net-tcp2"
)

func Tcp() {
	currentTcp := &tcp.Tcp{}
	for {
		fmt.Println("Finding new connections")

		tcp, connInfo, err := currentTcp.NewTcp(PROC_NET_TCP)
		currentTcp = tcp
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(connInfo)
		time.Sleep(60 * time.Second)
	}
}
