package rpchandlers

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/Nexell-AI-Network/nexelliad/v2/app/rpc/rpccontext"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/netadapter/router"
	"github.com/pkg/errors"
)

// HandleGetBalancesByAddresses handles the respectively named RPC command
func HandleGetBalancesByAddresses(context *rpccontext.Context, _ *router.Router, request appmessage.Message) (appmessage.Message, error) {
	if !context.Config.UTXOIndex {
		errorMessage := &appmessage.GetBalancesByAddressesResponseMessage{}
		errorMessage.Error = appmessage.RPCErrorf("Method unavailable when nexelliad is run without --utxoindex")
		return errorMessage, nil
	}

	getBalancesByAddressesRequest := request.(*appmessage.GetBalancesByAddressesRequestMessage)

	allEntries := make([]*appmessage.BalancesByAddressesEntry, len(getBalancesByAddressesRequest.Addresses))
	for i, address := range getBalancesByAddressesRequest.Addresses {
		balance, err := getBalanceByAddress(context, address)

		if err != nil {
			rpcError := &appmessage.RPCError{}
			if !errors.As(err, &rpcError) {
				return nil, err
			}
			errorMessage := &appmessage.GetUTXOsByAddressesResponseMessage{}
			errorMessage.Error = rpcError
			return errorMessage, nil
		}
		allEntries[i] = &appmessage.BalancesByAddressesEntry{
			Address: address,
			Balance: balance,
		}
	}

	response := appmessage.NewGetBalancesByAddressesResponse(allEntries)
	return response, nil
}
