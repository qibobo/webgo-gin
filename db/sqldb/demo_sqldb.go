package sqldb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/qibobo/webgo-gin/db"
	"github.com/qibobo/webgo-gin/models"
)

type DemoSQLDB struct {
	dbConfig db.DatabaseConfig
	logger   zap.Logger
	sqldb    *sql.DB
}

func NewDemoSQLDB(dbConfig db.DatabaseConfig, logger zap.Logger) (*DemoSQLDB, error) {
	sqldb, err := sql.Open(db.MysqlDriverName, dbConfig.URL)
	if err != nil {
		return nil, err
	}

	err = sqldb.Ping()
	if err != nil {
		sqldb.Close()
		logger.Error("ping demo db", zap.Error(err), zap.Any("dbConfig", dbConfig))
		return nil, err
	}

	sqldb.SetConnMaxLifetime(dbConfig.ConnectionMaxLifetime)
	sqldb.SetMaxIdleConns(dbConfig.MaxIdleConnections)
	sqldb.SetMaxOpenConns(dbConfig.MaxOpenConnections)
	sqldb.SetConnMaxIdleTime(dbConfig.ConnectionMaxIdleTime)

	return &DemoSQLDB{
		dbConfig: dbConfig,
		logger:   logger,
		sqldb:    sqldb,
	}, nil
}

func (ddb *DemoSQLDB) Close() error {
	err := ddb.sqldb.Close()
	if err != nil {
		ddb.logger.Error("Close demo db", zap.Error(err), zap.Any("dbConfig", ddb.dbConfig))
		return err
	}
	return nil
}

func (ddb *DemoSQLDB) GetDemo(id int) (*models.Demo, error) {
	var resultId int
	var resultName string
	query := "SELECT * FROM demo WHERE id = ?"
	err := ddb.sqldb.QueryRow(query, id).Scan(&resultId, &resultName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		ddb.logger.Error("get demo", zap.Error(err), zap.String("query", query), zap.Int("id", id))
		return nil, err
	}
	demo := &models.Demo{
		Id:   resultId,
		Name: resultName,
	}

	return demo, nil
}
func (ddb *DemoSQLDB) CreateDemo(demo *models.Demo) error {
	query := "INSERT INTO demo (id, name) VALUES(?, ?)"
	_, err := ddb.sqldb.Exec(query, demo.Id, demo.Name)
	if err != nil {
		ddb.logger.Error("create demo", zap.Error(err), zap.Int("id", demo.Id), zap.String("name", demo.Name))
	}
	return err
}
func (ddb *DemoSQLDB) DeleteDemo(id string) error {
	query := "DELETE FROM demo WHERE id = ?"
	_, err := ddb.sqldb.Exec(query, id)
	if err != nil {
		ddb.logger.Error("failed to delete demo", zap.Error(err), zap.String("query", query), zap.String("id", id))
	}
	return err
}

func (ddb *DemoSQLDB) GetDBStatus() sql.DBStats {
	return ddb.sqldb.Stats()
}
