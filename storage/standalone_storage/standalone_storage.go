package standalone_storage

import (
	"path/filepath"

	"github.com/pingcap/badger"

	"github.com/OpenKikCoc/raftkv/config"
	"github.com/OpenKikCoc/raftkv/storage"
	"github.com/OpenKikCoc/raftkv/util/engine_util"
	"github.com/pingcap/kvproto/pkg/kvrpcpb"
	//"github.com/OpenKikCoc/raftkv/proto"
)

type StandAloneReader struct {
	kvTxn   *badger.Txn
	raftTxn *badger.Txn
}

func (reader *StandAloneReader) getTxn(cf string) *badger.Txn {
	if cf == "raft" {
		return reader.raftTxn
	}
	return reader.kvTxn
}

func (reader *StandAloneReader) GetCF(cf string, key []byte) ([]byte, error) {
	txn := reader.getTxn(cf)
	val, err := engine_util.GetCFFromTxn(txn, cf, key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return val, err
}

func (reader *StandAloneReader) IterCF(cf string) engine_util.DBIterator {
	txn := reader.getTxn(cf)
	return engine_util.NewCFIterator(cf, txn)
}

// ---------------------------------------------------------------------------

var (
	kvSubpath   = "kv"
	raftSubpath = "raft"
)

// StandAloneStorage is an implementation of `Storage` for a single-node instance.
type StandAloneStorage struct {
	engine *engine_util.Engines
	config *config.Config
}

func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	dbPath := conf.DBConfig.DBPath
	kvPath := filepath.Join(dbPath, kvSubpath)
	raftPath := filepath.Join(dbPath, raftSubpath)

	kvDB := engine_util.CreateDB(kvPath, false)
	raftDB := engine_util.CreateDB(raftPath, true)
	return &StandAloneStorage{
		engine: engine_util.NewEngines(kvDB, raftDB, kvPath, raftPath),
		config: conf,
	}
}

func (s *StandAloneStorage) Start() error {
	return nil
}

func (s *StandAloneStorage) Stop() error {
	return s.engine.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.Reader, error) {
	var (
		kvTxn   = s.engine.Kv.NewTransaction(false)
		raftTxn = s.engine.Kv.NewTransaction(false)
	)
	return &StandAloneReader{
		kvTxn:   kvTxn,
		raftTxn: raftTxn,
	}, nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	for _, m := range batch {
		switch m.Data.(type) {
		case storage.Put:
			put := m.Data.(storage.Put)
			var txn *badger.Txn
			if put.Cf == "raft" {
				txn = s.engine.Raft.NewTransaction(true)
			} else {
				txn = s.engine.Kv.NewTransaction(true)
			}
			err := txn.Set(engine_util.KeyWithCF(put.Cf, put.Key), put.Value)
			if err != nil {
				return err
			}
			err = txn.Commit(nil)
			if err != nil {
				return err
			}
		case storage.Delete:
			delete := m.Data.(storage.Delete)
			var txn *badger.Txn
			if delete.Cf == "raft" {
				txn = s.engine.Raft.NewTransaction(true)
			} else {
				txn = s.engine.Kv.NewTransaction(true)
			}
			err := txn.Delete(engine_util.KeyWithCF(delete.Cf, delete.Key))
			if err != nil {
				return err
			}
			err = txn.Commit(nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
