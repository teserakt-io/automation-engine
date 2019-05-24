package models

import (
	"errors"

	"github.com/jinzhu/gorm"

	// Load available database drivers
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mysql"
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
	Migrate() error
}

// DBConfig holds generic database options and configuration
type DBConfig struct {
	Dialect   string
	CnxString string
	LogMode   bool
}

type gormDB struct {
	db     *gorm.DB
	config DBConfig
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
		err = ErrUnsupportedDialect
	}

	if err != nil {
		return nil, err
	}

	db.LogMode(config.LogMode)

	return &gormDB{
		db:     db,
		config: config,
	}, nil
}

func (gdb *gormDB) Migrate() error {

	if gdb.config.Dialect == DBDialectSQLite {
		// Enable foreign key support for sqlite3
		gdb.Connection().Exec("PRAGMA foreign_keys = ON")
	}

	result := gdb.Connection().AutoMigrate(
		Rule{},
		Trigger{},
		Target{},
	)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (gdb *gormDB) Connection() *gorm.DB {
	return gdb.db
}

func (gdb *gormDB) Close() error {
	return gdb.db.Close()
}
