package headersselectedtipstore

import (
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
)

type headersSelectedTipStagingShard struct {
	store          *headerSelectedTipStore
	newSelectedTip *externalapi.DomainHash
}

func (hsts *headerSelectedTipStore) stagingShard(stagingArea *model.StagingArea) *headersSelectedTipStagingShard {
	return stagingArea.GetOrCreateShard(hsts.shardID, func() model.StagingShard {
		return &headersSelectedTipStagingShard{
			store:          hsts,
			newSelectedTip: nil,
		}
	}).(*headersSelectedTipStagingShard)
}

func (hstss *headersSelectedTipStagingShard) Commit(dbTx model.DBTransaction) error {
	if hstss.newSelectedTip == nil {
		return nil
	}

	selectedTipBytes, err := hstss.store.serializeHeadersSelectedTip(hstss.newSelectedTip)
	if err != nil {
		return err
	}
	err = dbTx.Put(hstss.store.key, selectedTipBytes)
	if err != nil {
		return err
	}
	hstss.store.cache = hstss.newSelectedTip

	return nil
}

func (hstss *headersSelectedTipStagingShard) isStaged() bool {
	return hstss.newSelectedTip != nil
}
