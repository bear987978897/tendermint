package mempool

import (
	"encoding/binary"
	"testing"

	"github.com/tendermint/tendermint/abci/example/kvstore"
	"github.com/tendermint/tendermint/proxy"
)

func BenchmarkReap(b *testing.B) {
	app := kvstore.NewKVStoreApplication()
	cc := proxy.NewLocalClientCreator(app)
	mempool, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	size := 10000
	for i := 0; i < size; i++ {
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		mempool.CheckTx(tx, nil)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mempool.ReapMaxBytesMaxGas(100000000, 10000000)
	}
}

func BenchmarkCheckTx(b *testing.B) {
	app := kvstore.NewKVStoreApplication()
	cc := proxy.NewLocalClientCreator(app)
	mempool := newMempoolWithApp(cc)

	for i := 0; i < b.N; i++ {
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		mempool.CheckTx(tx, nil)
	}
}

func BenchmarkCacheInsertTime(b *testing.B) {
	cache := newMapTxCache(b.N)
	txs := make([][]byte, b.N)
	txInfo := TxInfo{UnknownPeerID}
	for i := 0; i < b.N; i++ {
		txs[i] = make([]byte, 8)
		binary.BigEndian.PutUint64(txs[i], uint64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.PushTxWithInfo(txs[i], txInfo)
	}
}

// This benchmark is probably skewed, since we actually will be removing
// txs in parallel, which may cause some overhead due to mutex locking.
func BenchmarkCacheRemoveTime(b *testing.B) {
	cache := newMapTxCache(b.N)
	txs := make([][]byte, b.N)
	txInfo := TxInfo{UnknownPeerID}
	for i := 0; i < b.N; i++ {
		txs[i] = make([]byte, 8)
		binary.BigEndian.PutUint64(txs[i], uint64(i))
		cache.PushTxWithInfo(txs[i], txInfo)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Remove(txs[i])
	}
}
