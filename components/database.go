package components

import (
	"fmt"

	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/instance-id/GoVerifier-dgo/appconfig"
)

type DbConfig struct {
	Db   *appconfig.DbSettings
	Xorm *XormDB
}

type XormDB struct {
	Engine      *xorm.Engine
	dbChnl      chan dbQuery
	closeWorker chan error
	runit       bool
}

func (xdb *DbConfig) ConnectDB(d *appconfig.DbSettings) *DbConfig {
	dbConfig := &DbConfig{
		Db: d,
		Xorm: &XormDB{
			Engine: func() *xorm.Engine {
				eng, err := xorm.NewEngine(d.Database, DetermineConnection(d))
				if err != nil {
					logrus.Fatalf("Database Connection Error: %s", err)
				}
				return eng
			}(),
			dbChnl:      make(chan dbQuery, 32),
			closeWorker: make(chan error),
		},
	}
	return dbConfig
}

func DetermineConnection(d *appconfig.DbSettings) string {
	var connString string
	switch d.Database {
	case "mysql":
		connString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", d.Data.Username, d.Data.Password, d.Data.Address, d.Data.DbName)
	case "mssql":
		connString = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", d.Data.Address, d.Data.Username, d.Data.Password, d.Data.DbName)
	case "postgres":
		connString = fmt.Sprintf("%s:%s@%s:5432/%s?sslmode=disable", d.Data.Username, d.Data.Password, d.Data.Address, d.Data.DbName)
	case "sqlite":
		connString = fmt.Sprintf("%s:%s@%s:5432/%s?sslmode=disable", d.Data.Username, d.Data.Password, d.Data.Address, d.Data.DbName)
	}
	return connString
}

func (x *XormDB) Run() {
	for x.dbChnl != nil {
		ev, ok := <-x.dbChnl
		if !ok {
			break
		}
		ev.Query()
		ev.Done()
	}
	// close
	x.closeWorker <- x.Engine.Close()
	x.Engine = nil
}

func (x *XormDB) Close() (err error) {
	c := x.dbChnl
	x.dbChnl = nil
	close(c)
	err = <-x.closeWorker
	close(x.closeWorker)
	return
}

// xorm reverse mysql instance:!WE2er#$@tcp(instance.id:3306)/verify?charset=utf8 templates/goxorm
