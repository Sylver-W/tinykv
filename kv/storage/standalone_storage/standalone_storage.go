package standalone_storage

import (
	"path"

	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct {
	// Your Data Here (1).
	engine *engine_util.Engines
	conf   *config.Config
	//在engine基础上的一层封装，保存conf以便将来使用
}

func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	// Your Code Here (1).
	dbPath := conf.DBPath
	kvPath := path.Join(dbPath, "kv")
	raftPath := path.Join(dbPath, "raft")
	kvdb := engine_util.CreateDB(kvPath, false)
	raftdb := engine_util.CreateDB(raftPath, true)
	return &StandAloneStorage{
		engine: engine_util.NewEngines(kvdb, raftdb, kvPath, raftPath),
		conf:   conf,
	}
}

func (s *StandAloneStorage) Start() error {
	// Your Code Here (1).
	return nil
}

func (s *StandAloneStorage) Stop() error {
	// Your Code Here (1).
	return s.engine.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	// Your Code Here (1).
	txn := s.engine.Kv.NewTransaction(false)
	return NewStandAloneStorageReader(txn), nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	// Your Code Here (1).
	var err error
	for _, m := range batch {
		key, val, cf := m.Key(), m.Value(), m.Cf()
		if _, ok := m.Data.(storage.Put); ok {
			err = engine_util.PutCF(s.engine.Kv, cf, key, val)
		} else {
			err = engine_util.DeleteCF(s.engine.Kv, cf, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type StandAloneStorageReader struct {
	KvTxn *badger.Txn
}

func NewStandAloneStorageReader(txn *badger.Txn) *StandAloneStorageReader {
	return &StandAloneStorageReader{
		KvTxn: txn,
	}
}

func (sr *StandAloneStorageReader) GetCF(cf string, key []byte) ([]byte, error) {
	val, err := engine_util.GetCFFromTxn(sr.KvTxn, cf, key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return val, err
}

func (sr *StandAloneStorageReader) IterCF(cf string) engine_util.DBIterator {
	return engine_util.NewCFIterator(cf, sr.KvTxn)
}

func (sr *StandAloneStorageReader) Close() {
	sr.KvTxn.Discard()
}
