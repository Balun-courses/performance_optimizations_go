package config

import (
	"fmt"
	"net/netip"
)

type Config interface {
	GetApiAddress() (netip.AddrPort, error)
}

var _ Config = (*ConstantConfig)(nil)

type ConstantConfig struct {
}

func NewConstantConfig() Config {
	return &ConstantConfig{}
}

func (c *ConstantConfig) GetApiAddress() (netip.AddrPort, error) {
	address, err := netip.ParseAddrPort("0.0.0.0:8801")

	if err != nil {
		return address, fmt.Errorf("can not parse constant address, error - %w", err)
	}

	return address, nil
}
