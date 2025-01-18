package rpccontext

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/protocol"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/utxoindex"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/config"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/addressmanager"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/connmanager"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/netadapter"
)

// Context represents the RPC context
type Context struct {
	Config            *config.Config
	NetAdapter        *netadapter.NetAdapter
	Domain            domain.Domain
	ProtocolManager   *protocol.Manager
	ConnectionManager *connmanager.ConnectionManager
	AddressManager    *addressmanager.AddressManager
	UTXOIndex         *utxoindex.UTXOIndex
	ShutDownChan      chan<- struct{}

	NotificationManager *NotificationManager
}

// NewContext creates a new RPC context
func NewContext(cfg *config.Config,
	domain domain.Domain,
	netAdapter *netadapter.NetAdapter,
	protocolManager *protocol.Manager,
	connectionManager *connmanager.ConnectionManager,
	addressManager *addressmanager.AddressManager,
	utxoIndex *utxoindex.UTXOIndex,
	shutDownChan chan<- struct{}) *Context {

	context := &Context{
		Config:            cfg,
		NetAdapter:        netAdapter,
		Domain:            domain,
		ProtocolManager:   protocolManager,
		ConnectionManager: connectionManager,
		AddressManager:    addressManager,
		UTXOIndex:         utxoIndex,
		ShutDownChan:      shutDownChan,
	}
	context.NotificationManager = NewNotificationManager(cfg.ActiveNetParams)

	return context
}
