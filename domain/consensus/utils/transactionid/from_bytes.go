package transactionid

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
)

// FromBytes creates a DomainTransactionID from the given byte slice
func FromBytes(transactionIDBytes []byte) (*externalapi.DomainTransactionID, error) {
	hash, err := externalapi.NewDomainHashFromByteSlice(transactionIDBytes)
	if err != nil {
		return nil, err
	}
	return (*externalapi.DomainTransactionID)(hash), nil
}