package rmodel

import (
	"database/sql"
)

const (
	BREADCRUMBS_UPDATE_START = "UPDATE `breadcrumbs_windows` SET `end_at`=UNIX_TIMESTAMP()+? WHERE `user_id`=? AND `cross_id`=? AND `end_at`>=UNIX_TIMESTAMP()"
	BREADCRUMBS_INSERT_START = "INSERT INTO `breadcrumbs_windows` (`user_id`, `cross_id`, `start_at`, `end_at`) VALUES(?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()+?)"
	BREADCRUMBS_UPDATE_END   = "UPDATE `breadcrumbs_windows` SET `end_at`=UNIX_TIMESTAMP()-1 WHERE `user_id`=? AND `cross_id`=? AND `end_at`>=UNIX_TIMESTAMP()"
	BREADCRUMBS_GET_END      = "SELECT `end_at` FROM `breadcrumbs_windows` WHERE `user_id`=? AND `cross_id`=? ORDER BY `end_at` DESC LIMIT 1"
	BREADCRUMBS_SAVE         = "INSERT INTO `breadcrumbs` (`user_id`, `lat`, `lng`, `acc`, `timestamp`) VALUES(?, ?, ?, ?, UNIX_TIMESTAMP());"
	BREADCRUMBS_GET          = "SELECT b.lat, b.lng, b.acc, b.timestamp FROM breadcrumbs AS b, breadcrumbs_windows AS w WHERE b.user_id=w.user_id AND b.timestamp BETWEEN w.start_at AND w.end_at AND w.user_id=? AND w.cross_id=? AND b.timestamp<=? AND b.timestamp>? ORDER BY b.timestamp DESC LIMIT 100"
	BREADCRUMBS_UPDATE       = "UPDATE `breadcrumbs` SET lat=?, lng=?, acc=?, timestamp=UNIX_TIMESTAMP() WHERE user_id=? ORDER BY timestamp DESC LIMIT 1"
)

type BreadcrumbsSaver struct {
	db          *sql.DB
	updateStart *sql.Stmt
	insertStart *sql.Stmt
	updateEnd   *sql.Stmt
	getEnd      *sql.Stmt
	save        *sql.Stmt
	get         *sql.Stmt
	update      *sql.Stmt
}

func NewBreadcrumbsSaver(db *sql.DB) (*BreadcrumbsSaver, error) {
	p := NewErrPrepare(db)
	ret := &BreadcrumbsSaver{
		db:          db,
		updateStart: p.Prepare(BREADCRUMBS_UPDATE_START),
		insertStart: p.Prepare(BREADCRUMBS_INSERT_START),
		updateEnd:   p.Prepare(BREADCRUMBS_UPDATE_END),
		getEnd:      p.Prepare(BREADCRUMBS_GET_END),
		save:        p.Prepare(BREADCRUMBS_SAVE),
		get:         p.Prepare(BREADCRUMBS_GET),
		update:      p.Prepare(BREADCRUMBS_UPDATE),
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *BreadcrumbsSaver) EnableCross(userId, crossId int64, afterInSecond int) error {
	res, err := s.updateStart.Exec(afterInSecond, userId, crossId)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r > 0 {
		return nil
	}
	if _, err := s.insertStart.Exec(userId, crossId, afterInSecond); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) DisableCross(userId, crossId int64) error {
	if _, err := s.updateEnd.Exec(userId, crossId); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) GetWindowEnd(userId, crossId int64) (int64, error) {
	row := s.getEnd.QueryRow(userId, crossId)
	var ret int64
	if err := row.Scan(&ret); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return ret, nil
}

func (s *BreadcrumbsSaver) Save(userId int64, l SimpleLocation) error {
	if _, err := s.save.Exec(userId, l.GPS[0], l.GPS[1], l.GPS[2]); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Update(userId int64, l SimpleLocation) error {
	if _, err := s.update.Exec(l.GPS[0], l.GPS[1], l.GPS[2], userId); err != nil {
		return err
	}
	return nil
}

func (s *BreadcrumbsSaver) Load(userId, crossId, afterTimestamp int64) ([]SimpleLocation, error) {
	rows, err := s.get.Query(userId, crossId, afterTimestamp, afterTimestamp-24*60*60)
	if err != nil {
		return nil, err
	}
	var ret []SimpleLocation
	for rows.Next() {
		var l SimpleLocation
		err := rows.Scan(&l.GPS[0], &l.GPS[1], &l.GPS[2], &l.Timestamp)
		if err != nil {
			return nil, err
		}
		ret = append(ret, l)
	}
	return ret, nil
}
