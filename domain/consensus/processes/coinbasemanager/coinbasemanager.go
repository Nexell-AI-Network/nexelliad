package coinbasemanager

import (
	"math"

	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/model/externalapi"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/constants"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/hashset"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/subnetworks"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/transactionhelper"
	"github.com/Nexell-AI-Network/nexelliad/v2/domain/consensus/utils/txscript"
	"github.com/Nexell-AI-Network/nexelliad/v2/infrastructure/db/database"
	"github.com/Nexell-AI-Network/nexelliad/v2/util"
	"github.com/pkg/errors"
)

type coinbaseManager struct {
	subsidyGenesisReward                    uint64
	preHalvingPhaseBaseSubsidy              uint64
	coinbasePayloadScriptPublicKeyMaxLength uint8
	genesisHash                             *externalapi.DomainHash
	halvingPhaseDaaScore                    uint64
	halvingPhaseBaseSubsidy                 uint64

	databaseContext     model.DBReader
	dagTraversalManager model.DAGTraversalManager
	ghostdagDataStore   model.GHOSTDAGDataStore
	acceptanceDataStore model.AcceptanceDataStore
	daaBlocksStore      model.DAABlocksStore
	blockStore          model.BlockStore
	pruningStore        model.PruningStore
	blockHeaderStore    model.BlockHeaderStore
}

func (c *coinbaseManager) ExpectedCoinbaseTransaction(stagingArea *model.StagingArea, blockHash *externalapi.DomainHash,
	coinbaseData *externalapi.DomainCoinbaseData) (expectedTransaction *externalapi.DomainTransaction, hasRedReward bool, err error) {

	ghostdagData, err := c.ghostdagDataStore.Get(c.databaseContext, stagingArea, blockHash, true)
	if !database.IsNotFoundError(err) && err != nil {
		return nil, false, err
	}

	// If there's ghostdag data with trusted data we prefer it because we need the original merge set non-pruned merge set.
	if database.IsNotFoundError(err) {
		ghostdagData, err = c.ghostdagDataStore.Get(c.databaseContext, stagingArea, blockHash, false)
		if err != nil {
			return nil, false, err
		}
	}

	acceptanceData, err := c.acceptanceDataStore.Get(c.databaseContext, stagingArea, blockHash)
	if err != nil {
		return nil, false, err
	}

	daaAddedBlocksSet, err := c.daaAddedBlocksSet(stagingArea, blockHash)
	if err != nil {
		return nil, false, err
	}

	txOuts := make([]*externalapi.DomainTransactionOutput, 0, len(ghostdagData.MergeSetBlues()))
	acceptanceDataMap := acceptanceDataFromArrayToMap(acceptanceData)
	if constants.BlockVersion == 1 {
		for _, blue := range ghostdagData.MergeSetBlues() {
			txOut, hasReward, err := c.coinbaseOutputForBlueBlockV1(stagingArea, blue, acceptanceDataMap[*blue], daaAddedBlocksSet)
			if err != nil {
				return nil, false, err
			}

			if hasReward {
				txOuts = append(txOuts, txOut)
			}
		}

		txOut, hasRedReward, err := c.coinbaseOutputForRewardFromRedBlocksV1(
			stagingArea, ghostdagData, acceptanceData, daaAddedBlocksSet, coinbaseData)
		if err != nil {
			return nil, false, err
		}

		if hasRedReward {
			txOuts = append(txOuts, txOut)
		}
	} else if constants.BlockVersion == 2 {
		for _, blue := range ghostdagData.MergeSetBlues() {
			txOut, devTx, hasReward, err := c.coinbaseOutputForBlueBlockV2(stagingArea, blue, acceptanceDataMap[*blue], daaAddedBlocksSet)
			if err != nil {
				return nil, false, err
			}

			if hasReward {
				txOuts = append(txOuts, txOut)
				txOuts = append(txOuts, devTx)
			}
		}

		txOut, devTx, hasRedReward, err := c.coinbaseOutputForRewardFromRedBlocksV2(
			stagingArea, ghostdagData, acceptanceData, daaAddedBlocksSet, coinbaseData)
		if err != nil {
			return nil, false, err
		}

		if hasRedReward {
			txOuts = append(txOuts, txOut)
			txOuts = append(txOuts, devTx)
		}
	}

	subsidy, err := c.CalcBlockSubsidy(stagingArea, blockHash)
	if err != nil {
		return nil, false, err
	}

	payload, err := c.serializeCoinbasePayload(ghostdagData.BlueScore(), coinbaseData, subsidy)
	if err != nil {
		return nil, false, err
	}

	domainTransaction := &externalapi.DomainTransaction{
		Version:      constants.MaxTransactionVersion,
		Inputs:       []*externalapi.DomainTransactionInput{},
		Outputs:      txOuts,
		LockTime:     0,
		SubnetworkID: subnetworks.SubnetworkIDCoinbase,
		Gas:          0,
		Payload:      payload,
	}
	return domainTransaction, hasRedReward, nil
}

func (c *coinbaseManager) daaAddedBlocksSet(stagingArea *model.StagingArea, blockHash *externalapi.DomainHash) (
	hashset.HashSet, error) {

	daaAddedBlocks, err := c.daaBlocksStore.DAAAddedBlocks(c.databaseContext, stagingArea, blockHash)
	if err != nil {
		return nil, err
	}

	return hashset.NewFromSlice(daaAddedBlocks...), nil
}

// coinbaseOutputForBlueBlock calculates the output that should go into the coinbase transaction of blueBlock
// If blueBlock gets no fee - returns nil for txOut
func (c *coinbaseManager) coinbaseOutputForBlueBlockV2(stagingArea *model.StagingArea,
	blueBlock *externalapi.DomainHash, blockAcceptanceData *externalapi.BlockAcceptanceData,
	mergingBlockDAAAddedBlocksSet hashset.HashSet) (*externalapi.DomainTransactionOutput, *externalapi.DomainTransactionOutput, bool, error) {

	blockReward, err := c.calcMergedBlockReward(stagingArea, blueBlock, blockAcceptanceData, mergingBlockDAAAddedBlocksSet)
	if err != nil {
		return nil, nil, false, err
	}

	devFeeDecodedAddress, err := util.DecodeAddress(constants.DevFeeAddress, util.Bech32PrefixKaspa)
	if err != nil {
		return nil, nil, false, err
	}
	devFeeScriptPublicKey, err := txscript.PayToAddrScript(devFeeDecodedAddress)
	if err != nil {
		return nil, nil, false, err
	}
	devFeeQuantity := uint64(float64(constants.DevFee) / 100 * float64(blockReward))
	blockReward -= devFeeQuantity
	if blockReward <= 0 {
		return nil, nil, false, nil
	}

	// the ScriptPublicKey for the coinbase is parsed from the coinbase payload
	_, coinbaseData, _, err := c.ExtractCoinbaseDataBlueScoreAndSubsidy(blockAcceptanceData.TransactionAcceptanceData[0].Transaction)
	if err != nil {
		return nil, nil, false, err
	}

	txOut := &externalapi.DomainTransactionOutput{
		Value:           blockReward,
		ScriptPublicKey: coinbaseData.ScriptPublicKey,
	}

	devTx := &externalapi.DomainTransactionOutput{
		Value:           devFeeQuantity,
		ScriptPublicKey: devFeeScriptPublicKey,
	}

	return txOut, devTx, true, nil
}

func (c *coinbaseManager) coinbaseOutputForBlueBlockV1(stagingArea *model.StagingArea,
	blueBlock *externalapi.DomainHash, blockAcceptanceData *externalapi.BlockAcceptanceData,
	mergingBlockDAAAddedBlocksSet hashset.HashSet) (*externalapi.DomainTransactionOutput, bool, error) {

	blockReward, err := c.calcMergedBlockReward(stagingArea, blueBlock, blockAcceptanceData, mergingBlockDAAAddedBlocksSet)
	if err != nil {
		return nil, false, err
	}

	if blockReward <= 0 {
		return nil, false, nil
	}

	// the ScriptPublicKey for the coinbase is parsed from the coinbase payload
	_, coinbaseData, _, err := c.ExtractCoinbaseDataBlueScoreAndSubsidy(blockAcceptanceData.TransactionAcceptanceData[0].Transaction)
	if err != nil {
		return nil, false, err
	}

	txOut := &externalapi.DomainTransactionOutput{
		Value:           blockReward,
		ScriptPublicKey: coinbaseData.ScriptPublicKey,
	}

	return txOut, true, nil
}

func (c *coinbaseManager) coinbaseOutputForRewardFromRedBlocksV2(stagingArea *model.StagingArea,
	ghostdagData *externalapi.BlockGHOSTDAGData, acceptanceData externalapi.AcceptanceData, daaAddedBlocksSet hashset.HashSet,
	coinbaseData *externalapi.DomainCoinbaseData) (*externalapi.DomainTransactionOutput, *externalapi.DomainTransactionOutput, bool, error) {

	acceptanceDataMap := acceptanceDataFromArrayToMap(acceptanceData)
	totalReward := uint64(0)
	for _, red := range ghostdagData.MergeSetReds() {
		reward, err := c.calcMergedBlockReward(stagingArea, red, acceptanceDataMap[*red], daaAddedBlocksSet)
		if err != nil {
			return nil, nil, false, err
		}
		totalReward += reward
	}

	devFeeDecodedAddress, err := util.DecodeAddress(constants.DevFeeAddress, util.Bech32PrefixKaspa)
	if err != nil {
		return nil, nil, false, err
	}
	devFeeScriptPublicKey, err := txscript.PayToAddrScript(devFeeDecodedAddress)
	if err != nil {
		return nil, nil, false, err
	}
	devFeeQuantity := uint64(float64(constants.DevFee) / 100 * float64(totalReward))
	totalReward -= devFeeQuantity
	if totalReward <= 0 {
		return nil, nil, false, nil
	}

	txOut := &externalapi.DomainTransactionOutput{
		Value:           totalReward,
		ScriptPublicKey: coinbaseData.ScriptPublicKey,
	}

	devTx := &externalapi.DomainTransactionOutput{
		Value:           devFeeQuantity,
		ScriptPublicKey: devFeeScriptPublicKey,
	}

	return txOut, devTx, true, nil
}

func (c *coinbaseManager) coinbaseOutputForRewardFromRedBlocksV1(stagingArea *model.StagingArea,
	ghostdagData *externalapi.BlockGHOSTDAGData, acceptanceData externalapi.AcceptanceData, daaAddedBlocksSet hashset.HashSet,
	coinbaseData *externalapi.DomainCoinbaseData) (*externalapi.DomainTransactionOutput, bool, error) {

	acceptanceDataMap := acceptanceDataFromArrayToMap(acceptanceData)
	totalReward := uint64(0)
	for _, red := range ghostdagData.MergeSetReds() {
		reward, err := c.calcMergedBlockReward(stagingArea, red, acceptanceDataMap[*red], daaAddedBlocksSet)
		if err != nil {
			return nil, false, err
		}
		totalReward += reward
	}
	if totalReward <= 0 {
		return nil, false, nil
	}

	txOut := &externalapi.DomainTransactionOutput{
		Value:           totalReward,
		ScriptPublicKey: coinbaseData.ScriptPublicKey,
	}

	return txOut, true, nil
}

func acceptanceDataFromArrayToMap(acceptanceData externalapi.AcceptanceData) map[externalapi.DomainHash]*externalapi.BlockAcceptanceData {
	acceptanceDataMap := make(map[externalapi.DomainHash]*externalapi.BlockAcceptanceData, len(acceptanceData))
	for _, blockAcceptanceData := range acceptanceData {
		acceptanceDataMap[*blockAcceptanceData.BlockHash] = blockAcceptanceData
	}
	return acceptanceDataMap
}

// CalcBlockSubsidy returns the subsidy amount a block at the provided blue score
// should have. This is mainly used for determining how much the coinbase for
// newly generated blocks awards as well as validating the coinbase for blocks
// has the expected value.
func (c *coinbaseManager) CalcBlockSubsidy(stagingArea *model.StagingArea, blockHash *externalapi.DomainHash) (uint64, error) {
	if blockHash.Equal(c.genesisHash) {
		return c.subsidyGenesisReward, nil
	}
	blockDaaScore, err := c.daaBlocksStore.DAAScore(c.databaseContext, stagingArea, blockHash)
	if err != nil {
		return 0, err
	}
	if blockDaaScore < c.halvingPhaseDaaScore {
		return c.preHalvingPhaseBaseSubsidy, nil
	}

	blockSubsidy := c.calcHalvingPeriodBlockSubsidy(blockDaaScore)
	return blockSubsidy, nil
}

func (c *coinbaseManager) calcHalvingPeriodBlockSubsidy(blockDaaScore uint64) uint64 {
	// We define a year as 365.25 days and a year as 365.25 / 12 = 30.4375
	// secondsPerMonth = 30.4375 * 24 * 60 * 60 = 2629800
	// secondsPerYear = 2629800 * 12
	const secondsPerYear = 31557600
	// Note that this calculation implicitly assumes that block per second = 1 (by assuming daa score diff is in second units).
	yearsSinceHalvingPhaseStarted := (blockDaaScore - c.halvingPhaseDaaScore) / secondsPerYear
	// yearsSinceHalvingHalvingPhaseStarted := (blockDaaScore - c.halvingPhaseDaaScore) / secondsPerYear
	// Return the pre-calculated value from subsidy-per-year table
	return c.getHalvingPeriodBlockSubsidyFromTable(yearsSinceHalvingPhaseStarted)
}

/*
This table was pre-calculated by calling `calcHalvingPeriodBlockSubsidyFloatCalc` for all years until reaching 0 subsidy.
To regenerate this table, run `TestBuildSubsidyTable` in coinbasemanager_test.go (note the `halvingPhaseBaseSubsidy` therein)
*/
var subsidyByHalvingYearTable = []uint64{
	600000000, 300000000, 150000000, 75000000, 37500000, 18750000, 
	9375000, 4687500, 2343750, 1171875, 585937, 292968, 146484, 
	73242, 36621, 18310, 9155, 4577, 2288, 1144, 572, 
	286, 143, 71, 35, 17, 8, 4, 2, 1, 0,
}

func (c *coinbaseManager) getHalvingPeriodBlockSubsidyFromTable(year uint64) uint64 {
	if year >= uint64(len(subsidyByHalvingYearTable)) {
		year = uint64(len(subsidyByHalvingYearTable) - 1)
	}
	return subsidyByHalvingYearTable[year]
}

func (c *coinbaseManager) calcHalvingPeriodBlockSubsidyFloatCalc(year uint64) uint64 {
	baseSubsidy := c.halvingPhaseBaseSubsidy
	subsidy := float64(baseSubsidy) / math.Pow(2, float64(year))
	return uint64(subsidy)
}

func (c *coinbaseManager) calcMergedBlockReward(stagingArea *model.StagingArea, blockHash *externalapi.DomainHash,
	blockAcceptanceData *externalapi.BlockAcceptanceData, mergingBlockDAAAddedBlocksSet hashset.HashSet) (uint64, error) {

	if !blockHash.Equal(blockAcceptanceData.BlockHash) {
		return 0, errors.Errorf("blockAcceptanceData.BlockHash is expected to be %s but got %s",
			blockHash, blockAcceptanceData.BlockHash)
	}

	if !mergingBlockDAAAddedBlocksSet.Contains(blockHash) {
		return 0, nil
	}

	totalFees := uint64(0)
	for _, txAcceptanceData := range blockAcceptanceData.TransactionAcceptanceData {
		if txAcceptanceData.IsAccepted {
			totalFees += txAcceptanceData.Fee
		}
	}

	block, err := c.blockStore.Block(c.databaseContext, stagingArea, blockHash)
	if err != nil {
		return 0, err
	}

	_, _, subsidy, err := c.ExtractCoinbaseDataBlueScoreAndSubsidy(block.Transactions[transactionhelper.CoinbaseTransactionIndex])
	if err != nil {
		return 0, err
	}

	return subsidy + totalFees, nil
}

// New instantiates a new CoinbaseManager
func New(
	databaseContext model.DBReader,

	subsidyGenesisReward uint64,
	preHalvingPhaseBaseSubsidy uint64,
	coinbasePayloadScriptPublicKeyMaxLength uint8,
	genesisHash *externalapi.DomainHash,
	halvingPhaseDaaScore uint64,
	halvingPhaseBaseSubsidy uint64,

	dagTraversalManager model.DAGTraversalManager,
	ghostdagDataStore model.GHOSTDAGDataStore,
	acceptanceDataStore model.AcceptanceDataStore,
	daaBlocksStore model.DAABlocksStore,
	blockStore model.BlockStore,
	pruningStore model.PruningStore,
	blockHeaderStore model.BlockHeaderStore) model.CoinbaseManager {

	return &coinbaseManager{
		databaseContext: databaseContext,

		subsidyGenesisReward:                    subsidyGenesisReward,
		preHalvingPhaseBaseSubsidy:              preHalvingPhaseBaseSubsidy,
		coinbasePayloadScriptPublicKeyMaxLength: coinbasePayloadScriptPublicKeyMaxLength,
		genesisHash:                             genesisHash,
		halvingPhaseDaaScore:                    halvingPhaseDaaScore,
		halvingPhaseBaseSubsidy:                 halvingPhaseBaseSubsidy,

		dagTraversalManager: dagTraversalManager,
		ghostdagDataStore:   ghostdagDataStore,
		acceptanceDataStore: acceptanceDataStore,
		daaBlocksStore:      daaBlocksStore,
		blockStore:          blockStore,
		pruningStore:        pruningStore,
		blockHeaderStore:    blockHeaderStore,
	}
}
