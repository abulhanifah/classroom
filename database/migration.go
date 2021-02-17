package database

import (
	"sort"

	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/models"
	"github.com/jinzhu/gorm"
)

var MigrationDB *gorm.DB

func Migrate() {
	hasNewMigration := false
	setting := models.Setting{Key: "db.migration.version"}
	MigrationDB.AutoMigrate(&setting)
	MigrationDB.Where(models.Setting{Key: setting.Key}).FirstOrCreate(&setting)

	index := make([]string, 0)
	for i, _ := range migration {
		index = append(index, i)
	}
	sort.Strings(index)
	for _, i := range index {
		if setting.Value == "" || setting.Value < i {
			migration[i]()
			setting.Value = i
			hasNewMigration = true
		}
	}
	if hasNewMigration {
		MigrationDB.Where(models.Setting{Key: setting.Key}).Assign(setting).FirstOrCreate(&setting)
	}
}

var migration = map[string]func(){
	"0001": func() { MigrationDB.AutoMigrate(&models.UserRole{}) },
	"0002": func() { MigrationDB.AutoMigrate(&helpers.Oauth2Client{}) },
	"0003": func() { MigrationDB.AutoMigrate(&helpers.Oauth2RedirectUrl{}) },
	"0004": func() { MigrationDB.AutoMigrate(&helpers.Oauth2Token{}) },
	"0005": func() { MigrationDB.AutoMigrate(&models.User{}) },
	"0006": func() { MigrationDB.AutoMigrate(&models.Classroom{}) },
	"0007": func() { MigrationDB.AutoMigrate(&models.Seat{}) },
	"0008": func() { MigrationDB.AutoMigrate(&models.Learning{}) },
}
