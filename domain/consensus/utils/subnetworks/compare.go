package subnetworks

import (
	"bytes"

	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
)

// Less returns true iff id a is less than id b
func Less(a, b externalapi.DomainSubnetworkID) bool {
	return bytes.Compare(a[:], b[:]) < 0
}
