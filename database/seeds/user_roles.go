package seeds

import (
	"github.com/jinzhu/gorm"

	"github.com/abulhanifah/classroom/models"
)

func SeedUserRole(db *gorm.DB) {
	for _, o := range UserRole {
		db.Where(models.UserRole{ID: o.ID}).Assign(o).FirstOrCreate(&o)
	}
}

var UserRole = []models.UserRole{
	{ID: 1, Name: "Admin"},
	{ID: 2, Name: "Pengajar"},
	{ID: 3, Name: "Murid"},
}
