package mock

import (
	"github.com/ElrondNetwork/elrond-go/data/block"
	"github.com/ElrondNetwork/elrond-go/data/state"
)

// EpochValidatorInfoCreatorStub -
type EpochValidatorInfoCreatorStub struct {
	CreateValidatorInfoMiniBlocksCalled func(validatorInfos map[uint32][]*state.ValidatorInfo) (block.MiniBlockSlice, error)
	VerifyValidatorInfoMiniBlocksCalled func(miniblocks []*block.MiniBlock, validatorInfos map[uint32][]*state.ValidatorInfo) error
	CreateMarshalizedDataCalled         func(body block.Body) map[string][][]byte
	SaveTxBlockToStorageCalled          func(metaBlock *block.MetaBlock, body block.Body)
	DeleteTxsFromStorageCalled          func(metaBlock *block.MetaBlock)
}

// CreateValidatorInfoMiniBlocks -
func (e *EpochValidatorInfoCreatorStub) CreateValidatorInfoMiniBlocks(validatorInfo map[uint32][]*state.ValidatorInfo) (block.MiniBlockSlice, error) {
	if e.CreateValidatorInfoMiniBlocksCalled != nil {
		return e.CreateValidatorInfoMiniBlocksCalled(validatorInfo)
	}
	return make(block.MiniBlockSlice, 0), nil
}

// VerifyValidatorInfoMiniBlocks -
func (e *EpochValidatorInfoCreatorStub) VerifyValidatorInfoMiniBlocks(miniblocks []*block.MiniBlock, validatorInfos map[uint32][]*state.ValidatorInfo) error {
	if e.VerifyValidatorInfoMiniBlocksCalled != nil {
		return e.VerifyValidatorInfoMiniBlocksCalled(miniblocks, validatorInfos)
	}
	return nil
}

// SaveValidatorInfoBlocksToStorage -
func (e *EpochValidatorInfoCreatorStub) SaveValidatorInfoBlocksToStorage(metaBlock *block.MetaBlock, body block.Body) {
	if e.SaveTxBlockToStorageCalled != nil {
		e.SaveTxBlockToStorageCalled(metaBlock, body)
	}
}

// DeleteValidatorInfoBlocksFromStorage -
func (e *EpochValidatorInfoCreatorStub) DeleteValidatorInfoBlocksFromStorage(metaBlock *block.MetaBlock) {
	if e.DeleteTxsFromStorageCalled != nil {
		e.DeleteTxsFromStorageCalled(metaBlock)
	}
}

// IsInterfaceNil -
func (e *EpochValidatorInfoCreatorStub) IsInterfaceNil() bool {
	return e == nil
}