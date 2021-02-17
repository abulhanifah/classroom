package database

import (
	"sort"

	"github.com/abulhanifah/classroom/database/seeds"
	"github.com/abulhanifah/classroom/models"
)

func Seed() {
	hasNewSeed := false
	setting := models.Setting{Key: "db.seed.version"}
	MigrationDB.Where(models.Setting{Key: setting.Key}).FirstOrCreate(&setting)

	index := make([]string, 0)
	for i, _ := range seed {
		index = append(index, i)
	}
	sort.Strings(index)
	for _, i := range index {
		if setting.Value == "" || setting.Value < i {
			seed[i]()
			setting.Value = i
			hasNewSeed = true
		}
	}
	if hasNewSeed {
		MigrationDB.Where(models.Setting{Key: setting.Key}).Assign(setting).FirstOrCreate(&setting)
	}
}

var seed = map[string]func(){
	"0001": func() { seeds.SeedOauth2Client(MigrationDB) },
	"0002": func() { seeds.SeedUserRole(MigrationDB) },
	"0003": func() { seeds.SeedUser(MigrationDB) },
}
