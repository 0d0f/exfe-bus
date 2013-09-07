package rmodel

import (
	"bytes"
	"database/sql"
	"fmt"
)

type Routex struct {
	CrossId   int64 `json:"cross_id,omitempty"`
	Enable    bool  `json:"enable,omitempty"`
	UpdatedAt int64 `json:"updated_at, omitempty"`
}

const (
	ROUTEX_SETUP_INSERT = "INSERT IGNORE INTO `routex` (`cross_id`, `updated_at`) VALUES(?, UNIX_TIMESTAMP())"
	ROUTEX_SETUP_UPDATE = "UPDATE `routex` SET `updated_at`=UNIX_TIMESTAMP() WHERE `cross_id`=?"
	ROUTEX_SETUP_SEARCH = "SELECT `cross_id`, `updated_at` FROM `routex` WHERE `cross_id` IN (%s) ORDER BY `updated_at` DESC"
)

type RoutexSaver struct {
	db     *sql.DB
	insert *sql.Stmt
	update *sql.Stmt
}

func NewRoutexSaver(db *sql.DB) (*RoutexSaver, error) {
	p := NewErrPrepare(db)
	ret := &RoutexSaver{
		db:     db,
		insert: p.Prepare(ROUTEX_SETUP_INSERT),
		update: p.Prepare(ROUTEX_SETUP_UPDATE),
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *RoutexSaver) Search(crossIds []int64) ([]Routex, error) {
	ids := bytes.NewBuffer(nil)
	for _, id := range crossIds {
		ids.WriteString(fmt.Sprintf("%d,", id))
	}
	if ids.Len() > 0 {
		ids.Truncate(ids.Len() - 1)
	}
	sql := fmt.Sprintf(ROUTEX_SETUP_SEARCH, ids.String())
	row, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	var ret []Routex
	for row.Next() {
		routex := Routex{
			Enable: true,
		}
		if err := row.Scan(&routex.CrossId, &routex.UpdatedAt); err != nil {
			return nil, err
		}
		ret = append(ret, routex)
	}
	return ret, nil
}

func (s *RoutexSaver) Update(crossId int64) error {
	res, err := s.insert.Exec(crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		if _, err := s.update.Exec(crossId); err != nil {
			return err
		}
	}
	return nil
}
