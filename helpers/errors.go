package helpers

import (
	"github.com/jinzhu/gorm"
)

func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
