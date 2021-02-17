package helpers

import (
	"database/sql"
	"log"

	"github.com/abulhanifah/classroom/configs"
	"github.com/abulhanifah/classroom/constants"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Connect() *gorm.DB {
	dbconf, err := configs.GetDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(dbconf.Driver, dbconf.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	if configs.GetConfig("DB_IS_DEBUG").Bool() {
		db = db.Debug()
	}

	maxOpenConns := configs.GetConfig("DB_MAX_OPEN_CONNS").Int()
	maxIdleConns := configs.GetConfig("DB_MAX_IDLE_CONNS").Int()
	connMaxLifetime := configs.GetConfig("DB_CONN_MAX_LIFETIME").Duration()

	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(connMaxLifetime)

	return db
}

func GetDB(ctx Context) *gorm.DB {
	return ctx.Get(constants.CtxTx).(*gorm.DB)
}

func GetResults(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	length := len(columns)
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		current := makeResultReceiver(length)
		if err := rows.Scan(current...); err != nil {
			panic(err)
		}
		value := make(map[string]interface{})
		for i := 0; i < length; i++ {
			value[columns[i]] = *(current[i]).(*interface{})
		}
		result = append(result, value)
	}
	return result
}

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}
