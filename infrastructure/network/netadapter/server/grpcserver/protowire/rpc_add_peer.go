package protowire

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/pkg/errors"
)

func (x *NexelliadMessage_AddPeerRequest) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_AddPeerRequest is nil")
	}
	return x.AddPeerRequest.toAppMessage()
}

func (x *NexelliadMessage_AddPeerRequest) fromAppMessage(message *appmessage.AddPeerRequestMessage) error {
	x.AddPeerRequest = &AddPeerRequestMessage{
		Address:     message.Address,
		IsPermanent: message.IsPermanent,
	}
	return nil
}

func (x *AddPeerRequestMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "AddPeerRequestMessage is nil")
	}
	return &appmessage.AddPeerRequestMessage{
		Address:     x.Address,
		IsPermanent: x.IsPermanent,
	}, nil
}

func (x *NexelliadMessage_AddPeerResponse) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_AddPeerResponse is nil")
	}
	return x.AddPeerResponse.toAppMessage()
}

func (x *NexelliadMessage_AddPeerResponse) fromAppMessage(message *appmessage.AddPeerResponseMessage) error {
	var err *RPCError
	if message.Error != nil {
		err = &RPCError{Message: message.Error.Message}
	}
	x.AddPeerResponse = &AddPeerResponseMessage{
		Error: err,
	}
	return nil
}

func (x *AddPeerResponseMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "AddPeerResponseMessage is nil")
	}
	rpcErr, err := x.Error.toAppMessage()
	// Error is an optional field
	if err != nil && !errors.Is(err, errorNil) {
		return nil, err
	}
	return &appmessage.AddPeerResponseMessage{
		Error: rpcErr,
	}, nil
}
