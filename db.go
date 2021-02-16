package classroom

import (
	"database/sql"
	"log"
)

func Connect() *sql.DB {
	dbconf, err := GetDBConfig()
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

	if GetConfig("DB_IS_DEBUG").Bool() {
		dbConnPool = dbConnPool.Debug()
	}

	maxOpenConns := GetConfig("DB_MAX_OPEN_CONNS").Int()
	maxIdleConns := GetConfig("DB_MAX_IDLE_CONNS").Int()
	connMaxLifetime := GetConfig("DB_CONN_MAX_LIFETIME").Duration()

	dbConnPool.DB().SetMaxIdleConns(maxIdleConns)
	dbConnPool.DB().SetMaxOpenConns(maxOpenConns)
	dbConnPool.DB().SetConnMaxLifetime(connMaxLifetime)

	return dbConnPool
}
