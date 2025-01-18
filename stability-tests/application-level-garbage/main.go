package main

import (
	"fmt"
	"os"

	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/config"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/netadapter/standalone"
	"github.com/Nexell-AI-Network/nexelliad/v2/stability-tests/common"
	"github.com/Nexell-AI-Network/nexelliad/v2/util/panics"
	"github.com/Nexell-AI-Network/nexelliad/v2/util/profiling"
)

func main() {
	defer panics.HandlePanic(log, "applicationLevelGarbage-main", nil)
	err := parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config: %+v", err)
		os.Exit(1)
	}
	defer backendLog.Close()
	common.UseLogger(backendLog, log.Level())
	cfg := activeConfig()
	if cfg.Profile != "" {
		profiling.Start(cfg.Profile, log)
	}

	nexelliadConfig := config.DefaultConfig()
	nexelliadConfig.NetworkFlags = cfg.NetworkFlags

	minimalNetAdapter, err := standalone.NewMinimalNetAdapter(nexelliadConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating minimalNetAdapter: %+v", err)
		backendLog.Close()
		os.Exit(1)
	}

	blocksChan, err := readBlocks()
	if err != nil {
		log.Errorf("Error reading blocks: %+v", err)
		backendLog.Close()
		os.Exit(1)
	}

	err = sendBlocks(cfg.NodeP2PAddress, minimalNetAdapter, blocksChan)
	if err != nil {
		log.Errorf("Error sending blocks: %+v", err)
		backendLog.Close()
		os.Exit(1)
	}
}
