package protowire

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/app/appmessage"
	"github.com/pkg/errors"
)

func (x *NexelliadMessage_BlockHeaders) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "NexelliadMessage_BlockHeaders is nil")
	}
	blockHeaders, err := x.BlockHeaders.toAppMessage()
	if err != nil {
		return nil, err
	}
	return &appmessage.BlockHeadersMessage{
		BlockHeaders: blockHeaders,
	}, nil
}

func (x *BlockHeadersMessage) toAppMessage() ([]*appmessage.MsgBlockHeader, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "BlockHeadersMessage is nil")
	}
	blockHeaders := make([]*appmessage.MsgBlockHeader, len(x.BlockHeaders))
	for i, blockHeader := range x.BlockHeaders {
		var err error
		blockHeaders[i], err = blockHeader.toAppMessage()
		if err != nil {
			return nil, err
		}
	}

	return blockHeaders, nil
}

func (x *NexelliadMessage_BlockHeaders) fromAppMessage(blockHeadersMessage *appmessage.BlockHeadersMessage) error {
	blockHeaders := make([]*BlockHeader, len(blockHeadersMessage.BlockHeaders))
	for i, blockHeader := range blockHeadersMessage.BlockHeaders {
		blockHeaders[i] = &BlockHeader{}
		err := blockHeaders[i].fromAppMessage(blockHeader)
		if err != nil {
			return err
		}
	}

	x.BlockHeaders = &BlockHeadersMessage{
		BlockHeaders: blockHeaders,
	}
	return nil
}
