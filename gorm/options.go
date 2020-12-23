/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

const defaultTablePrefix = ""

type options struct {
	TablePrefix           string        `json:"tablePrefix" yaml:"tablePrefix"`
	Host                  string        `json:"host,omitempty" yaml:"host" description:"db service host address"`
	Username              string        `json:"username,omitempty" yaml:"username"`
	Password              string        `json:"-" yaml:"password"`
	Type                  string        `json:"type" yaml:"type"`
	DBName                string        `json:"dbName" yaml:"dbName"`
	Debug                 bool          `json:"debug" yaml:"debug"`
	Port                  string        `json:"port" yaml:"port"`
	MaxConnMaxIdleTime    time.Duration `json:"maxConnMaxIdleTime,omitempty" yaml:"maxConnMaxIdleTime"`
	MaxOpenConnections    int           `json:"maxOpenConnections,omitempty" yaml:"maxOpenConnections"`
	MaxIdleConn           int           `json:"maxIdleConn,omitempty" yaml:"maxIdleConn"`
	MaxConnectionLifeTime time.Duration `json:"maxConnectionLifeTime,omitempty" yaml:"maxConnectionLifeTime"`
}

func newDatabaseOptions() *options {
	return &options{
		TablePrefix:           defaultTablePrefix,
		Host:                  "127.0.0.1",
		Username:              "root",
		Password:              "",
		DBName:                "",
		Type:                  "mysql",
		Port:                  "3306",
		Debug:                 true,
		MaxConnMaxIdleTime:    time.Duration(10) * time.Second,
		MaxOpenConnections:    100,
		MaxIdleConn:           100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
	}
}

func (m *options) getDSN() gorm.Dialector {
	switch m.Type {
	case "mysql":
		return mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&allowNativePasswords=true&parseTime=true", m.Username, m.Password, m.Host, m.Port, m.DBName))
	case "sqlite3":
		return sqlite.Open(m.DBName)
	case "postgres":
		fallthrough
	default:
		return postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", m.Host, m.Port, m.Username, m.DBName, m.Password))
	}
}
