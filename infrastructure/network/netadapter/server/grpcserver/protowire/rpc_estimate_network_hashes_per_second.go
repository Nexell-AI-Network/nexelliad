package protowire

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/pkg/errors"
)

func (x *NexelliadMessage_EstimateNetworkHashesPerSecondRequest) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_EstimateNetworkHashesPerSecondRequest is nil")
	}
	return x.EstimateNetworkHashesPerSecondRequest.toAppMessage()
}

func (x *NexelliadMessage_EstimateNetworkHashesPerSecondRequest) fromAppMessage(message *appmessage.EstimateNetworkHashesPerSecondRequestMessage) error {
	x.EstimateNetworkHashesPerSecondRequest = &EstimateNetworkHashesPerSecondRequestMessage{
		WindowSize: message.WindowSize,
		StartHash:  message.StartHash,
	}
	return nil
}

func (x *EstimateNetworkHashesPerSecondRequestMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "EstimateNetworkHashesPerSecondRequestMessage is nil")
	}
	return &appmessage.EstimateNetworkHashesPerSecondRequestMessage{
		WindowSize: x.WindowSize,
		StartHash:  x.StartHash,
	}, nil
}

func (x *NexelliadMessage_EstimateNetworkHashesPerSecondResponse) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_EstimateNetworkHashesPerSecondResponse is nil")
	}
	return x.EstimateNetworkHashesPerSecondResponse.toAppMessage()
}

func (x *NexelliadMessage_EstimateNetworkHashesPerSecondResponse) fromAppMessage(message *appmessage.EstimateNetworkHashesPerSecondResponseMessage) error {
	var err *RPCError
	if message.Error != nil {
		err = &RPCError{Message: message.Error.Message}
	}
	x.EstimateNetworkHashesPerSecondResponse = &EstimateNetworkHashesPerSecondResponseMessage{
		NetworkHashesPerSecond: message.NetworkHashesPerSecond,
		Error:                  err,
	}
	return nil
}

func (x *EstimateNetworkHashesPerSecondResponseMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "EstimateNetworkHashesPerSecondResponseMessage is nil")
	}
	rpcErr, err := x.Error.toAppMessage()
	// Error is an optional field
	if err != nil && !errors.Is(err, errorNil) {
		return nil, err
	}

	if rpcErr != nil && x.NetworkHashesPerSecond != 0 {
		return nil, errors.New("EstimateNetworkHashesPerSecondResponseMessage contains both an error and a response")
	}

	return &appmessage.EstimateNetworkHashesPerSecondResponseMessage{
		NetworkHashesPerSecond: x.NetworkHashesPerSecond,
		Error:                  rpcErr,
	}, nil
}
