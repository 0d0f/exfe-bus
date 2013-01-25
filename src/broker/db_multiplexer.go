package broker

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/googollee/go-logger"
	"github.com/googollee/go-multiplexer"
	"model"
	"time"
)

type DBInstance struct {
	db  *sql.DB
	log *logger.SubLogger
}

func (i *DBInstance) Ping() error {
	_, err := i.db.Exec("SELECT 1")
	if err != nil {
		return err
	}
	return nil
}

func (i *DBInstance) Close() error {
	return i.db.Close()
}

func (i *DBInstance) Error(err error) {
	i.log.Err("%s", err)
}

func (i *DBInstance) Query(sql string, v ...interface{}) (*sql.Rows, error) {
	result, err := i.db.Query(sql, v...)
	return result, err
}

func (i *DBInstance) Exec(sql string, v ...interface{}) (sql.Result, error) {
	result, err := i.db.Exec(sql, v...)
	return result, err
}

type DBMultiplexer struct {
	homo   *multiplexer.Homo
	config *model.Config
}

func NewDBMultiplexer(config *model.Config) *DBMultiplexer {
	if config.DB.MaxConnections == 0 {
		config.Log.Crit("config DB.MaxConnections should not 0!")
		panic("config DB.MaxConnections should not 0!")
	}
	return &DBMultiplexer{
		homo: multiplexer.NewHomo(func() (multiplexer.Instance, error) {
			db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&keepalive=1",
				config.DB.Username, config.DB.Password, config.DB.Addr, config.DB.Port, config.DB.DbName))
			if err != nil {
				return nil, err
			}
			_, err = db.Exec("SELECT 1")
			if err != nil {
				db.Close()
				return nil, err
			}
			return &DBInstance{
				db:  db,
				log: config.Log.SubPrefix("db"),
			}, nil
		}, config.DB.MaxConnections, -1, time.Duration(config.DB.HeartBeatInSecond)*time.Second),
		config: config,
	}
}

func (m *DBMultiplexer) Do(f func(multiplexer.Instance)) error {
	return m.homo.Do(f)
}

func (m *DBMultiplexer) Close() error {
	return m.homo.Close()
}
