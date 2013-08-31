package model

import (
	"database/sql"
	"fmt"
)

type Routex struct {
	CrossId   int64 `json:"cross_id,omitempty"`
	Enable    bool  `json:"enable,omitempty"`
	UpdatedAt int64 `json:"updated_at, omitempty"`
}

const (
	ROUTEX_SETUP_INSERT      = "INSERT IGNORE INTO `routex` (`user_id`, `cross_id`, `enable`, `updated_at`) VALUES(?, ?, ?, UNIX_TIMESTAMP())"
	ROUTEX_SETUP_UPDATE      = "UPDATE `routex` SET `enable`=?, `updated_at`=UNIX_TIMESTAMP() WHERE `user_id`=? AND `cross_id`=?"
	ROUTEX_SETUP_SEARCH      = "SELECT `cross_id`, `enable`, `updated_at` FROM `routex` WHERE `cross_id` IN (%s) GROUP By `cross_id` ORDER BY `updated_at` DESC"
	ROUTEX_SETUP_GET         = "SELECT `cross_id`, `enable`, `updated_at` FROM `routex` WHERE `user_id`=? AND `cross_id`=? LIMIT 1"
	ROUTEX_SETUP_ONLY_UPDATE = "UPDATE `routex` SET `updated_at`=UNIX_TIMESTAMP() WHERE `user_id`=? AND `cross_id`=?"
)

type RoutexSaver struct {
	db         *sql.DB
	insert     *sql.Stmt
	update     *sql.Stmt
	get        *sql.Stmt
	onlyUpdate *sql.Stmt
}

func NewRoutexSaver(db *sql.DB) (*RoutexSaver, error) {
	p := NewErrPrepare(db)
	ret := &RoutexSaver{
		db:         db,
		insert:     p.Prepare(ROUTEX_SETUP_INSERT),
		update:     p.Prepare(ROUTEX_SETUP_UPDATE),
		get:        p.Prepare(ROUTEX_SETUP_GET),
		onlyUpdate: p.Prepare(ROUTEX_SETUP_ONLY_UPDATE),
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *RoutexSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	res, err := s.update.Exec(true, userId, crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		if _, err := s.insert.Exec(userId, crossId, true); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoutexSaver) DisableCross(userId, crossId int64) error {
	res, err := s.update.Exec(false, userId, crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		if _, err := s.insert.Exec(userId, crossId, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoutexSaver) Search(crossIds []int64) ([]Routex, error) {
	ids := ""
	for _, id := range crossIds {
		ids = fmt.Sprintf("%s,%d", ids, id)
	}
	ids = ids[1:]
	sql := fmt.Sprintf(ROUTEX_SETUP_SEARCH, ids)
	row, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	var ret []Routex
	for row.Next() {
		var routex Routex
		if err := row.Scan(&routex.CrossId, &routex.Enable, &routex.UpdatedAt); err != nil {
			return nil, err
		}
		ret = append(ret, routex)
	}
	return ret, nil
}

func (s *RoutexSaver) Get(userId, crossId int64) (*Routex, error) {
	row := s.get.QueryRow(userId, crossId)
	var ret Routex
	if err := row.Scan(&ret.CrossId, &ret.Enable, &ret.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (s *RoutexSaver) Update(userId, crossId int64) error {
	_, err := s.onlyUpdate.Exec(userId, crossId)
	if err != nil {
		return err
	}
	return nil
}
