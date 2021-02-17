package seeds

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/abulhanifah/classroom/helpers"
)

func SeedOauth2Client(db *gorm.DB) {
	for _, o := range Oauth2Client {
		db.Where(helpers.Oauth2Client{ID: o.ID}).Assign(o).FirstOrCreate(&o)
	}
}

var Oauth2Client = []helpers.Oauth2Client{
	{
		ID:     helpers.NewToken(),
		Name:   "default",
		Secret: helpers.RandomString(10),
		UserID: uuid.New().String(),
	},
}
