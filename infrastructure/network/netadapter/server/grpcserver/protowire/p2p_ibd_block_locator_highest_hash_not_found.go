package protowire

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/pkg/errors"
)

func (x *NexelliadMessage_IbdBlockLocatorHighestHashNotFound) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_IbdBlockLocatorHighestHashNotFound is nil")
	}
	return &appmessage.MsgIBDBlockLocatorHighestHashNotFound{}, nil
}

func (x *NexelliadMessage_IbdBlockLocatorHighestHashNotFound) fromAppMessage(message *appmessage.MsgIBDBlockLocatorHighestHashNotFound) error {
	x.IbdBlockLocatorHighestHashNotFound = &IbdBlockLocatorHighestHashNotFoundMessage{}
	return nil
}
