package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/network/netadapter/server/grpcserver/protowire"
)

var commandTypes = []reflect.Type{
	reflect.TypeOf(protowire.NexelliadMessage_AddPeerRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetConnectedPeerInfoRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetPeerAddressesRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetCurrentNetworkRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetInfoRequest{}),

	reflect.TypeOf(protowire.NexelliadMessage_GetBlockRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetBlocksRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetHeadersRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetBlockCountRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetBlockDagInfoRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetSelectedTipHashRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetVirtualSelectedParentBlueScoreRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetVirtualSelectedParentChainFromBlockRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_ResolveFinalityConflictRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_EstimateNetworkHashesPerSecondRequest{}),

	reflect.TypeOf(protowire.NexelliadMessage_GetBlockTemplateRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_SubmitBlockRequest{}),

	reflect.TypeOf(protowire.NexelliadMessage_GetMempoolEntryRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetMempoolEntriesRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetMempoolEntriesByAddressesRequest{}),

	reflect.TypeOf(protowire.NexelliadMessage_SubmitTransactionRequest{}),

	reflect.TypeOf(protowire.NexelliadMessage_GetUtxosByAddressesRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetBalanceByAddressRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_GetCoinSupplyRequest{}),

	reflect.TypeOf(protowire.NexelliadMessage_BanRequest{}),
	reflect.TypeOf(protowire.NexelliadMessage_UnbanRequest{}),
}

type commandDescription struct {
	name       string
	parameters []*parameterDescription
	typeof     reflect.Type
}

type parameterDescription struct {
	name   string
	typeof reflect.Type
}

func commandDescriptions() []*commandDescription {
	commandDescriptions := make([]*commandDescription, len(commandTypes))

	for i, commandTypeWrapped := range commandTypes {
		commandType := unwrapCommandType(commandTypeWrapped)

		name := strings.TrimSuffix(commandType.Name(), "RequestMessage")
		numFields := commandType.NumField()

		var parameters []*parameterDescription
		for i := 0; i < numFields; i++ {
			field := commandType.Field(i)

			if !isFieldExported(field) {
				continue
			}

			parameters = append(parameters, &parameterDescription{
				name:   field.Name,
				typeof: field.Type,
			})
		}
		commandDescriptions[i] = &commandDescription{
			name:       name,
			parameters: parameters,
			typeof:     commandTypeWrapped,
		}
	}

	return commandDescriptions
}

func (cd *commandDescription) help() string {
	sb := &strings.Builder{}
	sb.WriteString(cd.name)
	for _, parameter := range cd.parameters {
		_, _ = fmt.Fprintf(sb, " [%s]", parameter.name)
	}
	return sb.String()
}
