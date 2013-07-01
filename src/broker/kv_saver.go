package broker

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	SAVER_SAVE = "INSERT IGNORE INTO `kv_saver` (`key`, `value`, `touched_at`) VALUES (?, ?, NOW())"
	SAVER_GET  = "SELECT `value` FROM `kv_saver` WHERE `key` IN (%s) limit 1"
)

type KVSaver struct {
	db *sql.DB
}

func NewKVSaver(db *sql.DB) *KVSaver {
	return &KVSaver{
		db: db,
	}
}

func (s *KVSaver) Save(keys []string, value string) error {
	for _, key := range keys {
		_, err := s.db.Exec(SAVER_SAVE, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *KVSaver) Check(keys []string) (string, bool, error) {
	keysSql := make([]string, len(keys))
	for i, key := range keys {
		keysSql[i] = fmt.Sprintf("\"%s\"", key)
	}
	query := fmt.Sprintf(SAVER_GET, strings.Join(keysSql, ","))
	var row *sql.Rows
	row, err := s.db.Query(query)
	if err != nil {
		return "", false, err
	}
	defer row.Close()

	value, exist := "", false
	for row.Next() {
		row.Scan(&value)
		exist = true
	}
	return value, exist, nil
}
