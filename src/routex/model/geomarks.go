package model

import (
	"database/sql"
	"encoding/json"
)

const (
	GEOMARKS_CREATE = "INSERT IGNORE INTO `geomarks` (`id`, `type`, `cross_id`, `mark`, `touched_at`, `deleted`) VALUES (?, ?, ?, ?, UNIX_TIMESTAMP(), FALSE)"
	GEOMARKS_UPDATE = "UPDATE `geomarks` SET `mark`=?, `touched_at`=UNIX_TIMESTAMP(), `deleted`=FALSE WHERE `id`=? AND `type`=? AND `cross_id`=?"
	GEOMARKS_GET    = "SELECT `mark` FROM `geomarks` WHERE `cross_id`=? AND `deleted`=FALSE"
	GEOMARKS_DELETE = "UPDATE `geomarks` SET `deleted`=TRUE, `touched_at`=UNIX_TIMESTAMP() WHERE `id`=? AND `type`=? AND `cross_id`=? AND `deleted`=FALSE"
)

type GeomarksSaver struct {
	db     *sql.DB
	create *sql.Stmt
	update *sql.Stmt
	get    *sql.Stmt
	del    *sql.Stmt
}

func NewGeomarkSaver(db *sql.DB) (*GeomarksSaver, error) {
	p := NewErrPrepare(db)
	ret := &GeomarksSaver{
		db:     db,
		create: p.Prepare(GEOMARKS_CREATE),
		update: p.Prepare(GEOMARKS_UPDATE),
		get:    p.Prepare(GEOMARKS_GET),
		del:    p.Prepare(GEOMARKS_DELETE),
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *GeomarksSaver) Set(crossId int64, mark Geomark) error {
	b, err := json.Marshal(mark)
	if err != nil {
		return err
	}
	n, err := s.create.Exec(mark.Id, mark.Type, crossId, string(b))
	if err != nil {
		return err
	}
	ret, err := n.RowsAffected()
	if err != nil {
		return err
	}
	if ret == 0 {
		mark.CreatedBy, mark.CreatedAt = mark.UpdatedBy, mark.UpdatedAt
		b, err = json.Marshal(mark)
		if err != nil {
			return err
		}
		_, err := s.update.Exec(string(b), mark.Id, mark.Type, crossId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *GeomarksSaver) Get(crossId int64) ([]Geomark, error) {
	rows, err := s.get.Query(crossId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Geomark
	for rows.Next() {
		var b string
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var mark Geomark
		if err := json.Unmarshal([]byte(b), &mark); err != nil {
			return nil, err
		}
		ret = append(ret, mark)
	}
	return ret, nil
}

func (s *GeomarksSaver) Delete(crossId int64, markType, markId string) error {
	if _, err := s.del.Exec(markId, markType, crossId); err != nil {
		return err
	}
	return nil
}
