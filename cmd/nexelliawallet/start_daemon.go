package main

import "github.com/Nexell-AI-Network/nexelliad/v2/cmd/nexelliawallet/daemon/server"

func startDaemon(conf *startDaemonConfig) error {
	return server.Start(conf.NetParams(), conf.Listen, conf.RPCServer, conf.KeysFile, conf.Profile, conf.Timeout)
}
