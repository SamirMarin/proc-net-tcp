package handlers

import (
	"fmt"
	"proc-net-tcp/pkg/tcp"
	"time"
)

const (
	//PROC_NET_TCP = "/proc/net/tcp"
	PROC_NET_TCP = "/Users/smarin/workspaces/teleport/scratchpad/proc-net-tcp"
)

func Tcp() {

	for {
		fmt.Println("Read proc/net/tcp")
		_, err := tcp.Connections(PROC_NET_TCP)
		if err != nil {
			fmt.Println(err)

		}
		time.Sleep(30 * time.Second)
	}

}
