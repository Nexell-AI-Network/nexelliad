package protowire

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/pkg/errors"
)

func (x *NexelliadMessage_Verack) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_Verack is nil")
	}
	return &appmessage.MsgVerAck{}, nil
}

func (x *NexelliadMessage_Verack) fromAppMessage(_ *appmessage.MsgVerAck) error {
	return nil
}
