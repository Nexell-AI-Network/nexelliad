package consensusstatestore

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/utxolrucache"
	"github.com/Nexell-AI-Network/nexelliad/v2/util/staging"
)

var importingPruningPointUTXOSetKeyName = []byte("importing-pruning-point-utxo-set")

// consensusStateStore represents a store for the current consensus state
type consensusStateStore struct {
	shardID                         model.StagingShardID
	virtualUTXOSetCache             *utxolrucache.LRUCache
	tipsCache                       []*externalapi.DomainHash
	tipsKey                         model.DBKey
	utxoSetBucket                   model.DBBucket
	importingPruningPointUTXOSetKey model.DBKey
}

// New instantiates a new ConsensusStateStore
func New(prefixBucket model.DBBucket, utxoSetCacheSize int, preallocate bool) model.ConsensusStateStore {
	return &consensusStateStore{
		shardID:                         staging.GenerateShardingID(),
		virtualUTXOSetCache:             utxolrucache.New(utxoSetCacheSize, preallocate),
		tipsKey:                         prefixBucket.Key(tipsKeyName),
		importingPruningPointUTXOSetKey: prefixBucket.Key(importingPruningPointUTXOSetKeyName),
		utxoSetBucket:                   prefixBucket.Bucket(utxoSetBucketName),
	}
}

func (css *consensusStateStore) IsStaged(stagingArea *model.StagingArea) bool {
	return css.stagingShard(stagingArea).isStaged()
}
