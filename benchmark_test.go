package benchmarks_test

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/pingcap/go-ycsb/pkg/generator"
	"github.com/tidwall/buntdb"
	bolt "go.etcd.io/bbolt"
)

const (
	boltDbPath   = "bolt.db"
	badgerDbPath = "badger.db"
	buntDbPath   = "bunt.db"
	numItems     = 100000
	keyPrefix    = "key"
	valueSize    = 100
	bucketName   = "benchmark"
)

func generateData(count int) ([]string, [][]byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	keys := make([]string, count)
	values := make([][]byte, count)

	zipfian := generator.NewScrambledZipfian(0, int64(count)-1, generator.ZipfianConstant)

	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("%s%d", keyPrefix, zipfian.Next(r))
		values[i] = make([]byte, valueSize)
		rand.Read(values[i])
	}

	return keys, values
}

func benchmarkBuntDB(b *testing.B, keys []string, values [][]byte) {
	db, err := buntdb.Open(buntDbPath)
	if err != nil {
		b.Fatal(err)
	}
	if err := db.SetConfig(buntdb.Config{
		SyncPolicy: buntdb.Always,
	}); err != nil {

	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll(buntDbPath)
	}()

	var lastIdx int
	b.Run("BuntDB-Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			err := db.Update(func(tx *buntdb.Tx) error {
				_, _, err := tx.Set(keys[idx], string(values[idx]), nil)
				return err
			})
			if err != nil {
				b.Fatal(err)
			}
			lastIdx = idx
		}
	})

	b.Run("BuntDB-Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % lastIdx
			err := db.View(func(tx *buntdb.Tx) error {
				_, err := tx.Get(keys[idx])
				return err
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkBadgerDB(b *testing.B, keys []string, values [][]byte) {
	err := os.Mkdir(badgerDbPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("Directory '%s' already exists\n", badgerDbPath)
		} else {
			fmt.Printf("Error creating directory: %v\n", err)
		}
	}

	opts := badger.DefaultOptions(badgerDbPath)
	opts.SyncWrites = true
	opts.Logger = newBadgerLogger(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})))

	db, err := badger.Open(opts)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll(badgerDbPath)
	}()

	var lastIdx int
	b.Run("BadgerDB-Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			err := db.Update(func(txn *badger.Txn) error {
				return txn.Set([]byte(keys[idx]), values[idx])
			})
			if err != nil {
				b.Fatal(err)
			}
			lastIdx = idx
		}
	})

	b.Run("BadgerDB-Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % lastIdx
			err := db.View(func(txn *badger.Txn) error {
				_, err := txn.Get([]byte(keys[idx]))
				return err
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkBBolt(b *testing.B, keys []string, values [][]byte) {
	err := os.Mkdir(boltDbPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("Directory '%s' already exists\n", boltDbPath)
		} else {
			fmt.Printf("Error creating directory: %v\n", err)
		}
	}

	db, err := bolt.Open(fmt.Sprintf(boltDbPath+"/db"), 0600, &bolt.Options{
		FreelistType: bolt.FreelistMapType,
		Timeout:      1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = db.Close()
		_ = os.RemoveAll(boltDbPath)
	}()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		b.Fatal(err)
	}

	var lastIdx int
	b.Run("BBolt-Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % len(keys)
			err := db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(bucketName))
				return bucket.Put([]byte(keys[idx]), values[idx])
			})
			if err != nil {
				b.Fatal(err)
			}
			lastIdx = idx
		}
	})

	b.Run("BBolt-Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			idx := i % lastIdx
			err := db.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(bucketName))
				_ = bucket.Get([]byte(keys[idx]))
				return nil
			})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDatabases(b *testing.B) {
	keys, values := generateData(numItems)

	b.Run("BuntDB", func(b *testing.B) {
		benchmarkBuntDB(b, keys, values)
	})

	b.Run("BadgerDB", func(b *testing.B) {
		benchmarkBadgerDB(b, keys, values)
	})

	b.Run("BBolt", func(b *testing.B) {
		benchmarkBBolt(b, keys, values)
	})
}
