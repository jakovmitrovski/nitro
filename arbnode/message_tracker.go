package arbnode

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/offchainlabs/nitro/arbos/arbostypes"
	"github.com/offchainlabs/nitro/arbos/l1pricing"
	"github.com/offchainlabs/nitro/arbos/l2pricing"
	"github.com/offchainlabs/nitro/arbutil"
)

var (
	MessageTrackingL1Prefix = []byte("msgtrack_l1")
	MessageTrackingL2Prefix = []byte("msgtrack_l2")
	LatestStateIndexPrefix  = []byte("msgtrack_lateststateindex")
	L1LatestStateIndexKey   = []byte("L1")
	L2LatestStateIndexKey   = []byte("L2")
)

type HasPrefix interface {
	TrackingPrefix() []byte
}

const ARBITRUM_ONE_GENESIS_BLOCK = 22207817

// L1 Tracking
type MessageTrackingL1Data struct {
	Message      arbostypes.L1IncomingMessage
	L1TxHash     common.Hash
	DataLocation batchDataLocation
}

// L2 Tracking
type MessageTrackingL2Data struct {
	L2BlockNumber          uint64
	L2BlockHash            common.Hash
	L1PricingState         l1pricing.L1PricingState
	L2PricingState         l2pricing.L2PricingState
	BrotliCompressionLevel uint64
}

type LatestStateIndex struct {
	StateIndex arbutil.MessageIndex
}

func messageTrackingKey(prefix []byte, msgNum arbutil.MessageIndex) []byte {
	var key [8]byte
	binary.BigEndian.PutUint64(key[:], uint64(msgNum))
	return append(prefix, key[:]...)
}

func (MessageTrackingL1Data) TrackingPrefix() []byte {
	return MessageTrackingL1Prefix
}

func (MessageTrackingL2Data) TrackingPrefix() []byte {
	return MessageTrackingL2Prefix
}

func AddTrackingData[T HasPrefix](db ethdb.Batch, msgNum arbutil.MessageIndex, data *T) error {
	if data == nil {
		return nil
	}
	prefix := (*data).TrackingPrefix()
	key := messageTrackingKey(prefix, msgNum)
	encoded, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}
	return db.Put(key, encoded)
}

func GetTrackingDataAt[T HasPrefix](db ethdb.Database, msgNum arbutil.MessageIndex) (*T, error) {
	var zero T
	prefix := zero.TrackingPrefix()
	key := messageTrackingKey(prefix, msgNum)
	data, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	var result T
	if err := rlp.DecodeBytes(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func SetLatestStateIndex(db ethdb.Batch, latestStateIndex LatestStateIndex, keySuffix []byte) error {
	encoded, err := rlp.EncodeToBytes(&latestStateIndex)
	if err != nil {
		return err
	}
	return db.Put(append(LatestStateIndexPrefix, keySuffix...), encoded)
}

func GetLatestStateIndex(db ethdb.Database, keySuffix []byte) (*LatestStateIndex, error) {
	data, err := db.Get(append(LatestStateIndexPrefix, keySuffix...))
	if err != nil {
		return nil, err
	}
	var result LatestStateIndex
	if err := rlp.DecodeBytes(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
