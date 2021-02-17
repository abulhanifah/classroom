package seeds

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/abulhanifah/classroom/helpers"
	"github.com/abulhanifah/classroom/models"
)

func SeedUser(db *gorm.DB) {
	for _, o := range User {
		temp := []models.User{}
		db.Where(models.User{Email: o.Email}).Limit(1).Find(&temp)
		if len(temp) == 0 {
			hash, err := helpers.Hash(o.Password)
			if err == nil {
				o.Password = string(hash)
			}
			db.Create(&o)
		}
	}
}

var User = []models.User{
	{ID: uuid.New().String(), Name: "Admin", Email: "admin@classroom.id", Password: "admin123456", RoleID: 1, IsEmailVerified: true},
	{ID: uuid.New().String(), Name: "Mr. Clark Kent", Email: "kent.super@classroom.id", Password: "teacher123456", RoleID: 2, IsEmailVerified: true},
	{ID: uuid.New().String(), Name: "Son Goku", Email: "gokuson@dragon.id", Password: "student123456", RoleID: 3, IsEmailVerified: true},
	{ID: uuid.New().String(), Name: "Son Gohan", Email: "gohan@dragon.id", Password: "student123456", RoleID: 3, IsEmailVerified: true},
	{ID: uuid.New().String(), Name: "Naruto", Email: "bunshin@ninja.co.jp", Password: "student123456", RoleID: 3, IsEmailVerified: true},
}
