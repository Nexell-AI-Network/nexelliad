package model

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
)

// DifficultyManager provides a method to resolve the
// difficulty value of a block
type DifficultyManager interface {
	StageDAADataAndReturnRequiredDifficulty(stagingArea *StagingArea, blockHash *externalapi.DomainHash, isBlockWithTrustedData bool) (uint32, error)
	RequiredDifficulty(stagingArea *StagingArea, blockHash *externalapi.DomainHash) (uint32, error)
	EstimateNetworkHashesPerSecond(startHash *externalapi.DomainHash, windowSize int) (uint64, error)
}
