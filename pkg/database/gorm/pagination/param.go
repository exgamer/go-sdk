package pagination

import "gorm.io/gorm"

type Param struct {
	DB     *gorm.DB
	Paging *Paging
}
