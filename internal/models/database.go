package models

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"

	// Load available database drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	// import _ "github.com/jinzhu/gorm/dialects/mysql"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"

	"github.com/teserakt-io/automation-engine/internal/config"
	slibcfg "github.com/teserakt-io/serverlib/config"
)

var (
	// ErrUnsupportedDialect is returned when creating a new database with a Config having an unsupported dialect
	ErrUnsupportedDialect = errors.New("unsupported database dialect")
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
	config config.DBCfg
	logger *log.Logger
}

var _ Database = &gormDB{}

// NewDB creates a new database
func NewDB(config config.DBCfg, logger *log.Logger) (Database, error) {
	var db *gorm.DB
	var err error

	cnxStr, err := config.ConnectionString()
	if err != nil {
		return nil, err
	}

	db, err = gorm.Open(config.Type.String(), cnxStr)
	if err != nil {
		return nil, err
	}

	db.LogMode(config.Logging)
	db.SetLogger(logger)

	return &gormDB{
		db:     db,
		config: config,
		logger: logger,
	}, nil
}

func (gdb *gormDB) Migrate() error {
	gdb.logger.Println("Database Migration Started.")

	switch gdb.config.Type {
	case slibcfg.DBTypeSQLite:
		// Enable foreign key support for sqlite3
		gdb.Connection().Exec("PRAGMA foreign_keys = ON")
	}

	result := gdb.Connection().AutoMigrate(
		Rule{},
		Trigger{},
		TriggerState{},
		Target{},
	)

	if result.Error != nil {
		return result.Error
	}

	gdb.logger.Println("Database Migration Finished.")

	return nil
}

func (gdb *gormDB) Connection() *gorm.DB {
	return gdb.db
}

func (gdb *gormDB) Close() error {
	return gdb.db.Close()
}
