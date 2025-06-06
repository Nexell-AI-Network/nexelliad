package model

import "github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"

// BlockParentBuilder exposes a method to build super-block parents for
// a given set of direct parents
type BlockParentBuilder interface {
	BuildParents(stagingArea *StagingArea,
		daaScore uint64,
		directParentHashes []*externalapi.DomainHash) ([]externalapi.BlockLevelParents, error)
}
