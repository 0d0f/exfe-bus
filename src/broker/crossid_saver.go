package broker

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	SAVER_SAVE = "INSERT IGNORE INTO `cross_saver` (`reference_id`, `cross_id`, `touched_at`) VALUES (?, ?, NOW())"
	SAVER_GET  = "SELECT cross_id FROM cross_saver WHERE reference_id IN (%s) limit 1"
)

type CrossSaver struct {
	db *sql.DB
}

func NewCrossSaver(db *sql.DB) *CrossSaver {
	return &CrossSaver{
		db: db,
	}
}

func (s *CrossSaver) Save(ids []string, crossID string) error {
	for _, id := range ids {
		_, err := s.db.Exec(SAVER_SAVE, id, crossID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CrossSaver) Check(ids []string) (string, bool, error) {
	ids_ := make([]string, len(ids))
	for i, id := range ids {
		ids_[i] = fmt.Sprintf("\"%s\"", id)
	}
	query := fmt.Sprintf(SAVER_GET, strings.Join(ids_, ","))
	var row *sql.Rows
	row, err := s.db.Query(query)
	if err != nil {
		return "", false, err
	}
	defer row.Close()

	crossID, crossExist := "", false
	for row.Next() {
		row.Scan(&crossID)
		crossExist = true
	}
	return crossID, crossExist, nil
}
