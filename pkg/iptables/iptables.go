package iptables

import (
	"github.com/coreos/go-iptables/iptables"
)

type IptablesClient struct {
	TableName string
	Chain     string
	Iptables  *iptables.IPTables
}

func NewIptables(params ...string) (*IptablesClient, error) {
	ipt, err := iptables.New()
	if err != nil {
		return nil, err
	}
	table := "filter"
	chain := "INPUT"
	if len(params) == 1 {
		table = params[0]
	} else if len(params) >= 2 {
		table = params[0]
		chain = params[1]
	}

	return &IptablesClient{
		TableName: table,
		Chain:     chain,
		Iptables:  ipt,
	}, nil
}

func (ic *IptablesClient) BlockSourceIp(ip string) error {
	err := ic.Iptables.Append(ic.TableName, ic.Chain, "-S", ip, "-j", "REJECT")
	return err
}
