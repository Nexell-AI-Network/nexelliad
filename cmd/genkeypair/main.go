package main

import (
	"fmt"

	"github.com/Nexell-AI-Network/nexelliad/v2/cmd/nexelliawallet/libkaspawallet"
	"github.com/Nexell-AI-Network/nexelliad/v2/util"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		panic(err)
	}

	privateKey, publicKey, err := libkaspawallet.CreateKeyPair(false)
	if err != nil {
		panic(err)
	}

	addr, err := util.NewAddressPublicKey(publicKey, cfg.NetParams().Prefix)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Private key: %x\n", privateKey)
	fmt.Printf("Address: %s\n", addr)
}
