package main

import (
	"sync"

	"github.com/abulhanifah/classroom/configs"
	"github.com/abulhanifah/classroom/database"
	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/routes"
	"github.com/jinzhu/gorm"
)

var (
	syncOnce sync.Once
	dbh      *gorm.DB
)

func main() {
	syncOnce.Do(DBConnection)
	defer dbh.Close()

	r := routes.Init(dbh)
	err := r.Start(":" + configs.GetConfig("APP_PORT").String())
	if err != nil {
		r.Logger.Fatal(err)
	}
}

func DBConnection() {
	dbh = helpers.Connect()
	database.MigrationDB = dbh
	database.Migrate()
	database.Seed()
}
