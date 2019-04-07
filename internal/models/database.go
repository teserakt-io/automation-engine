package models

import (
	"errors"

	"github.com/jinzhu/gorm"

	// Load available database drivers
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//import _ "github.com/jinzhu/gorm/dialects/mysql"
	// import _ "github.com/jinzhu/gorm/dialects/postgres"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	// ErrUnsupportedDialect is returned when creating a new database with a Config having an unsupported dialect
	ErrUnsupportedDialect = errors.New("unsupported database dialect")
)

const (
	// DBDialectSQLite ...
	DBDialectSQLite = "sqlite3"
)

// Database describes a generic database implementation
type Database interface {
	Close() error
	Connection() *gorm.DB
}

// DBConfig holds generic database options and configuration
type DBConfig struct {
	Dialect   string
	CnxString string
	LogMode   bool
	Models    []interface{}
}

type gormDB struct {
	db *gorm.DB
}

var _ Database = &gormDB{}

// NewDB creates a new database
func NewDB(config DBConfig) (Database, error) {

	var db *gorm.DB
	var err error

	switch config.Dialect {
	case DBDialectSQLite:
		db, err = gorm.Open(config.Dialect, config.CnxString)
	default:
		return nil, ErrUnsupportedDialect
	}

	if err != nil {
		return nil, err
	}

	if result := db.AutoMigrate(config.Models...); result.Error != nil {
		return nil, err
	}

	db.LogMode(config.LogMode)

	return &gormDB{
		db: db,
	}, nil

}

func (gdb *gormDB) Connection() *gorm.DB {
	return gdb.db
}

func (gdb *gormDB) Close() error {
	return gdb.db.Close()
}
