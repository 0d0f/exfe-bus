package main

import (
	"broker"
	"database/sql"
	"fmt"
	"github.com/googollee/go-multiplexer"
	"strings"
)

const (
	SAVER_SAVE = "INSERT IGNORE INTO `cross_saver` (`reference_id`, `cross_id`, `touched_at`) VALUES (?, ?, NOW())"
	SAVER_GET  = "SELECT cross_id FROM cross_saver WHERE reference_id IN (%s) limit 1"
)

type CrossSaver struct {
	db *broker.DBMultiplexer
}

func NewCrossSaver(db *broker.DBMultiplexer) *CrossSaver {
	return &CrossSaver{
		db: db,
	}
}

func (s *CrossSaver) Save(ids []string, crossID string) (err error) {
	err = s.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)

		for _, id := range ids {
			_, err = db.Exec(SAVER_SAVE, id, crossID)
			if err != nil {
				return
			}
		}
	})
	return
}

func (s *CrossSaver) Check(ids []string) (crossID string, crossExist bool, err error) {
	err = s.db.Do(func(i multiplexer.Instance) {
		db := i.(*broker.DBInstance)

		ids_ := make([]string, len(ids))
		for i, id := range ids {
			ids_[i] = fmt.Sprintf("\"%s\"", id)
		}
		query := fmt.Sprintf(SAVER_GET, strings.Join(ids_, ","))
		var row *sql.Rows
		row, err = db.Query(query)
		if err != nil {
			return
		}
		defer row.Close()

		for row.Next() {
			row.Scan(&crossID)
			crossExist = true
		}
	})
	return
}
