package cache

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nutsdb/nutsdb"
)

// NutsDB wraps the NutsDB instance and provides shared access.
type NutsDB struct {
	DB *nutsdb.DB
}

// Open opens or creates a NutsDB database at the given directory.
// Auto-creates the directory if it doesn't exist.
func Open(dir string) (*NutsDB, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(dir),
		nutsdb.WithSegmentSize(16*1024*1024),
		nutsdb.WithRWMode(nutsdb.MMap),
		nutsdb.WithEntryIdxMode(nutsdb.HintKeyValAndRAMIdxMode),
	)
	if err != nil {
		return nil, err
	}

	ndb := &NutsDB{DB: db}

	// Ensure buckets exist (ignore "already exists" error on subsequent runs)
	for _, bucket := range []string{"sessions", "users"} {
		_ = db.Update(func(tx *nutsdb.Tx) error {
			return tx.NewBucket(nutsdb.DataStructureBTree, bucket)
		})
	}

	log.Printf("[nutsdb] opened at %s", dir)
	return ndb, nil
}

// Close closes the NutsDB database.
func (n *NutsDB) Close() error {
	if n.DB != nil {
		return n.DB.Close()
	}
	return nil
}

// Path returns the default NutsDB data path relative to a base directory.
func Path(baseDir string) string {
	return filepath.Join(baseDir, "cache")
}

// int64Key encodes an int64 as big-endian bytes for NutsDB key.
func int64Key(id int64) []byte {
	return []byte{
		byte(id >> 56), byte(id >> 48), byte(id >> 40), byte(id >> 32),
		byte(id >> 24), byte(id >> 16), byte(id >> 8), byte(id),
	}
}

// ttlToUint32 converts time.Duration to NutsDB TTL in seconds (uint32).
// Returns 0 for 0 or negative durations (no TTL).
func ttlToUint32(d time.Duration) uint32 {
	if d <= 0 {
		return 0
	}
	secs := int64(d.Seconds())
	if secs <= 0 {
		return 0
	}
	if secs > 1<<32-1 {
		return 1<<32 - 1 // max uint32
	}
	return uint32(secs)
}
