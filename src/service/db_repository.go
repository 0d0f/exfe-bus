package main

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"model"
)

type DBRepository struct {
	db     *sql.DB
	Config *model.Config
}

func (r *DBRepository) Connect() error {
	if r.db != nil {
		r.db.Close()
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&keepalive=1",
		r.Config.DB.Username, r.Config.DB.Password, r.Config.DB.Addr, r.Config.DB.Port, r.Config.DB.DbName))
	if err != nil {
		return err
	}
	_, err = db.Query("SELECT 1")
	if err != nil {
		db.Close()
		return err
	}
	r.db = db
	return nil
}

func (r *DBRepository) Query(sql string, v ...interface{}) (*sql.Rows, error) {
	result, err := r.db.Query(sql, v...)
	if err == driver.ErrBadConn {
		err = r.Connect()
		if err == nil {
			result, err = r.db.Query(sql, v...)
		}
	}
	return result, err
}

func (r *DBRepository) Exec(sql string, v ...interface{}) (sql.Result, error) {
	result, err := r.db.Exec(sql, v...)
	if err == driver.ErrBadConn {
		err = r.Connect()
		if err == nil {
			result, err = r.db.Exec(sql, v...)
		}
	}
	return result, err
}
