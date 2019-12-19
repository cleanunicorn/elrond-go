package txpool

import (
	"sync"

	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/logger"
	"github.com/ElrondNetwork/elrond-go/storage"
	"github.com/ElrondNetwork/elrond-go/storage/storageUnit"
	"github.com/ElrondNetwork/elrond-go/storage/txcache"
)

var log = logger.GetOrCreate("dataretriever/txpool")

// shardedTxPool holds transaction caches organised by destination shard
type shardedTxPool struct {
	mutex             sync.RWMutex
	backingMap        map[string]*txPoolShard
	cacheConfig       storageUnit.CacheConfig
	mutexAddCallbacks sync.RWMutex
	onAddCallbacks    []func(key []byte)
}

type txPoolShard struct {
	CacheID string
	Cache   *txcache.TxCache
}

// NewShardedTxPool creates a new sharded tx pool
// Implements "dataRetriever.TxPool"
func NewShardedTxPool(config storageUnit.CacheConfig) dataRetriever.ShardedDataCacherNotifier {
	return &shardedTxPool{
		cacheConfig:       config,
		mutex:             sync.RWMutex{},
		backingMap:        make(map[string]*txPoolShard),
		mutexAddCallbacks: sync.RWMutex{},
		onAddCallbacks:    make([]func(key []byte), 0),
	}
}

// ShardDataStore is not implemented for this pool
func (txPool *shardedTxPool) ShardDataStore(cacheID string) storage.Cacher {
	cache := txPool.getTxCache(cacheID)
	return cache
}

// getTxCache returns the requested cache
func (txPool *shardedTxPool) getTxCache(cacheID string) *txcache.TxCache {
	shard := txPool.getOrCreateShard(cacheID)
	return shard.Cache
}

func (txPool *shardedTxPool) getOrCreateShard(cacheID string) *txPoolShard {
	txPool.mutex.RLock()
	shard, ok := txPool.backingMap[cacheID]
	txPool.mutex.RUnlock()

	if ok {
		return shard
	}

	// The cache not yet created, we'll create in a critical section
	txPool.mutex.Lock()
	defer txPool.mutex.Unlock()

	// We have to check again if not created (concurrency issue)
	shard, ok = txPool.backingMap[cacheID]
	if !ok {
		nChunksHint := txPool.cacheConfig.Shards
		cache := txcache.NewTxCache(nChunksHint)
		shard = &txPoolShard{
			CacheID: cacheID,
			Cache:   cache,
		}

		txPool.backingMap[cacheID] = shard
	}

	return shard
}

// AddData adds the transaction to the cache
func (txPool *shardedTxPool) AddData(key []byte, value interface{}, cacheID string) {
	txPool.addTx(key, value.(data.TransactionHandler), cacheID)
}

// addTx adds the transaction to the cache
func (txPool *shardedTxPool) addTx(txHash []byte, tx data.TransactionHandler, cacheID string) {
	shard := txPool.getOrCreateShard(cacheID)
	cache := shard.Cache
	_, added := cache.AddTx(txHash, tx)
	if added {
		txPool.onAdded(txHash)
	}
}

func (txPool *shardedTxPool) onAdded(txHash []byte) {
	txPool.mutexAddCallbacks.RLock()
	defer txPool.mutexAddCallbacks.RUnlock()

	for _, handler := range txPool.onAddCallbacks {
		go handler(txHash)
	}
}

// SearchFirstData is not implemented for this pool
func (txPool *shardedTxPool) SearchFirstData(key []byte) (interface{}, bool) {
	tx, ok := txPool.searchFirstTx(key)
	return tx, ok
}

// searchFirstTx searches the transaction against all shard data store, retrieving the first found
func (txPool *shardedTxPool) searchFirstTx(txHash []byte) (tx data.TransactionHandler, ok bool) {
	txPool.mutex.RLock()
	defer txPool.mutex.RUnlock()

	for key := range txPool.backingMap {
		shard := txPool.backingMap[key]
		tx, ok := shard.Cache.GetByTxHash(txHash)
		if ok {
			return tx, ok
		}
	}

	return nil, false
}

// RemoveData is not implemented for this pool
func (txPool *shardedTxPool) RemoveData(key []byte, cacheID string) {
	panic("shardedTxPool.RemoveData is not implemented")
}

// removeTx removes the transaction from the pool
func (txPool *shardedTxPool) removeTx(txHash []byte, cacheID string) {
	// shard := txPool.getOrCreateShard(cacheID)
	// _ = shard.Cache.RemoveTxByHash(txHash)
}

// RemoveSetOfDataFromPool is not implemented for this pool
func (txPool *shardedTxPool) RemoveSetOfDataFromPool(keys [][]byte, cacheID string) {
	//txPool.removeTxBulk(keys, cacheID)
}

// removeTxBulk removes a bunch of transactions from the pool
func (txPool *shardedTxPool) removeTxBulk(txHashes [][]byte, cacheID string) {
	// for _, key := range txHashes {
	// 	txPool.removeTx(key, cacheID)
	// }
}

// RemoveDataFromAllShards is not implemented for this pool
func (txPool *shardedTxPool) RemoveDataFromAllShards(key []byte) {
	panic("shardedTxPool.RemoveDataFromAllShards is not implemented")
}

// removeTxFromAllShards will remove the transaction from the pool (searches for it in all shards)
func (txPool *shardedTxPool) removeTxFromAllShards(txHash []byte) {
	// txPool.mutex.RLock()
	// defer txPool.mutex.RUnlock()

	// for _, shard := range txPool.backingMap {
	// 	cache := shard.Cache
	// 	_ = cache.RemoveTxByHash(txHash)
	// }
}

// MergeShardStores merges two shards of the pool
func (txPool *shardedTxPool) MergeShardStores(sourceCacheID, destCacheID string) {
	// sourceShard := txPool.getOrCreateShard(sourceCacheID)
	// sourceCache := sourceShard.Cache

	// sourceCache.ForEachTransaction(func(txHash []byte, tx data.TransactionHandler) {
	// 	txPool.addTx(txHash, tx, destCacheID)
	// })

	// txPool.mutex.Lock()
	// delete(txPool.backingMap, sourceCacheID)
	// txPool.mutex.Unlock()
}

// MoveData is not implemented for this pool
func (txPool *shardedTxPool) MoveData(sourceCacheID, destCacheID string, key [][]byte) {
	panic("shardedTxPool.MoveData is not implemented")
}

// MoveTxs moves the transactions between two caches
func (txPool *shardedTxPool) MoveTxs(sourceCacheID string, destCacheID string, txHashes [][]byte) {
	// sourceShard := txPool.getOrCreateShard(sourceCacheID)
	// sourceCache := sourceShard.Cache

	// for _, txHash := range txHashes {
	// 	tx, ok := sourceCache.GetByTxHash(txHash)
	// 	if ok {
	// 		txPool.addTx(txHash, tx, destCacheID)
	// 		txPool.removeTx(txHash, sourceCacheID)
	// 	}
	// }
}

// Clear clears everything in the pool
func (txPool *shardedTxPool) Clear() {
	// txPool.mutex.Lock()
	// for key := range txPool.backingMap {
	// 	delete(txPool.backingMap, key)
	// }
	// txPool.mutex.Unlock()
}

// ClearShardStore clears a specific cache
func (txPool *shardedTxPool) ClearShardStore(cacheID string) {
	// shard := txPool.getOrCreateShard(cacheID)
	// shard.Cache.Clear()
}

// CreateShardStore is a ShardedData method that is responsible for creating
//  a new shardStore with cacheId index in the shardedDataStore map
func (txPool *shardedTxPool) CreateShardStore(cacheID string) {
	panic("shardedTxPool.CreateShardStore() is not implemented (not needed; shard creation is managed internally)")
}

// RegisterHandler registers a new handler to be called when a new transaction is added
func (txPool *shardedTxPool) RegisterHandler(handler func(key []byte)) {
	// if handler == nil {
	// 	log.Error("attempt to register a nil handler")
	// 	return
	// }

	// txPool.mutexAddCallbacks.Lock()
	// txPool.onAddCallbacks = append(txPool.onAddCallbacks, handler)
	// txPool.mutexAddCallbacks.Unlock()
}

// IsInterfaceNil returns true if there is no value under the interface
func (txPool *shardedTxPool) IsInterfaceNil() bool {
	return txPool == nil
}
