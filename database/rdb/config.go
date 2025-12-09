package rdb

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Driver      string `json:"driver,omitempty"`
	Host        string `json:"host"`
	DataBase    string `json:"database"`
	UserName    string `json:"userName"`
	Password    string `json:"password,omitempty"`
	Port        int    `json:"port,omitempty"`
	MaxOpenCons int    `json:"maxOpenCons,omitempty"`
	MaxIdleCons int    `json:"maxIdleCons,omitempty"`
}

func (c *Config) Init() (err error) {
	switch c.Driver {
	case "", driverMysql:
		err = c.initMysql()
	case driverSqlite:
		err = c.initSqlite()
	default:
		err = fmt.Errorf("unSupport database driver %s", c.Driver)
	}
	return
}

func (c *Config) initSqlite() (err error) {
	if _, err = os.Stat(c.DataBase); err == nil {
		err = fmt.Errorf("sqlite db file %s not found", c.DataBase)
	} else {
		if c.MaxOpenCons <= 0 {
			c.MaxOpenCons = defaultSqliteMaxOpenCons
		}
		if c.MaxIdleCons <= 0 {
			c.MaxIdleCons = defaultSqliteMaxIdleCons
		}
	}
	return
}

func (c *Config) Dial() (dial gorm.Dialector) {
	switch c.Driver {
	case "", driverMysql:
		dial = mysql.Open(fmt.Sprintf(mysqlDataSourceNameFormat, c.UserName, c.Password, c.Host, c.Port, c.DataBase))
	case driverSqlite:
		dial = sqlite.Open(fmt.Sprintf(sqliteDataSourceNameFormat, c.DataBase))
	}
	return
}

func (c *Config) initMysql() (err error) {
	if c.MaxOpenCons <= 0 {
		c.MaxOpenCons = defaultMysqlMaxOpenCons
	}
	if c.MaxIdleCons <= 0 {
		c.MaxIdleCons = defaultMysqlMaxIdleCons
	}
	if c.Port <= 0 {
		c.Port = defaultMysqlPort
	}
	return
}

const (
	mysqlDataSourceNameFormat  = "%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local"
	sqliteDataSourceNameFormat = "file:%s?cache=shared&mode=rwc&_journal_mode=WAL"
	defaultMysqlPort           = 3306
	defaultMysqlMaxOpenCons    = 6
	defaultMysqlMaxIdleCons    = 6
	defaultSqliteMaxOpenCons   = 1
	defaultSqliteMaxIdleCons   = 1
)

const (
	driverMysql  = "mysql"
	driverSqlite = "sqlite"
)
