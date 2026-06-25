package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"md-note/backend/internal/config"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	gdb, err := gorm.Open(gormmysql.Open(cfg.MySQLDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)

	return gdb, nil
}

func RunMigrations(sqlDB *sql.DB, migrationsFS embed.FS) error {
	srcDriver, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("load migrations source: %w", err)
	}

	dbDriver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("create migrate db driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", srcDriver, "mysql", dbDriver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}
