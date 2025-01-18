package protowire

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/pkg/errors"
)

func (x *NexelliadMessage_InvRelayBlock) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_InvRelayBlock is nil")
	}
	return x.InvRelayBlock.toAppMessage()
}

func (x *InvRelayBlockMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "InvRelayBlockMessage is nil")
	}
	hash, err := x.Hash.toDomain()
	if err != nil {
		return nil, err
	}

	return &appmessage.MsgInvRelayBlock{Hash: hash}, nil

}

func (x *NexelliadMessage_InvRelayBlock) fromAppMessage(msgInvRelayBlock *appmessage.MsgInvRelayBlock) error {
	x.InvRelayBlock = &InvRelayBlockMessage{
		Hash: domainHashToProto(msgInvRelayBlock.Hash),
	}
	return nil
}
