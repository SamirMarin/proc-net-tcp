package tcp

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Client struct {
	Connections []ProcNetTcp
}

type ProcNetTcp struct {
	LocalAdress  string
	LocalPort    string
	RemoteAdress string
	RemotePort   string
	Inode        string
}

func Connections(filepath string) (*Client, error) {
	connectionsRead, err := readFile(filepath)
	if err != nil {
		return nil, err
	}
	fmt.Println(connectionsRead)

	for _, connection := range connectionsRead {
		fmt.Println(connection)
		for _, entry := range connection {
			fmt.Println(entry)
		}
	}

	connections := []ProcNetTcp{}

	return &Client{
		Connections: connections,
	}, nil
}

func readFile(path string) ([][]string, error) {
	tcpConncetionsFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	entries := strings.Split(string(tcpConncetionsFile), "\n")

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
		entriesClean = append(entriesClean, entryClean)
	}
	return entriesClean, nil
}
