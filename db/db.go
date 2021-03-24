package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/qibobo/webgo-gin/models"
)

const PostgresDriverName = "postgres"
const MysqlDriverName = "mysql"

type OrderType uint8

const (
	DESC OrderType = iota
	ASC
)
const (
	DESCSTR string = "DESC"
	ASCSTR  string = "ASC"
)

var ErrAlreadyExists = fmt.Errorf("already exists")
var ErrDoesNotExist = fmt.Errorf("doesn't exist")
var ErrConflict = fmt.Errorf("conflicting entry exists")

type DatabaseConfig struct {
	URL                   string        `yaml:"url"`
	MaxOpenConnections    int           `yaml:"max_open_connections"`
	MaxIdleConnections    int           `yaml:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `yaml:"connection_max_lifetime"`
	ConnectionMaxIdleTime time.Duration `yaml:"connection_max_idletime"`
}

type DatabaseStatus interface {
	GetDBStatus() sql.DBStats
}
type DemoDB interface {
	DatabaseStatus
	GetDemo(id int) (*models.Demo, error)
	CreateDemo(demo *models.Demo) error
	// SaveDemoInBulk(demo []*models.Demo) error
	// DeleteDemo(id string) error
	Close() error
}
