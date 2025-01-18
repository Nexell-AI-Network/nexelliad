package serialization

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
)

// DomainBlockStatusToDbBlockStatus converts model.BlockStatus to DbBlockStatus
func DomainBlockStatusToDbBlockStatus(domainBlockStatus externalapi.BlockStatus) *DbBlockStatus {
	return &DbBlockStatus{
		Status: uint32(domainBlockStatus),
	}
}

// DbBlockStatusToDomainBlockStatus converts DbBlockStatus to model.BlockStatus
func DbBlockStatusToDomainBlockStatus(dbBlockStatus *DbBlockStatus) externalapi.BlockStatus {
	return externalapi.BlockStatus(dbBlockStatus.Status)
}
