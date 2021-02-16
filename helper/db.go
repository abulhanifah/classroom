package helper

import (
	"database/sql"
	"log"

	"github.com/abulhanifah/classroom/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	dbconf, err := config.GetDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	dbConnPool, err := sql.Open(dbconf.Driver, dbconf.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbConnPool.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	if config.GetConfig("DB_IS_DEBUG").Bool() {
		dbConnPool = dbConnPool.Debug()
	}

	maxOpenConns := config.GetConfig("DB_MAX_OPEN_CONNS").Int()
	maxIdleConns := config.GetConfig("DB_MAX_IDLE_CONNS").Int()
	connMaxLifetime := config.GetConfig("DB_CONN_MAX_LIFETIME").Duration()

	dbConnPool.DB().SetMaxIdleConns(maxIdleConns)
	dbConnPool.DB().SetMaxOpenConns(maxOpenConns)
	dbConnPool.DB().SetConnMaxLifetime(connMaxLifetime)

	return dbConnPool
}
