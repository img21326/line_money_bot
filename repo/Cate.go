package repo

import (
	"log"

	"gorm.io/gorm"
)

type Cate struct {
	gorm.Model
	Name    string `json:"name"`
	Total   int64  `json:"total"`
	UserID  uint
	Account []Account `json:"account"`
}

type CateRepo struct {
	db *gorm.DB
}

func NewCateRepo(db *gorm.DB) *CateRepo {
	return &CateRepo{db}
}

func (r *CateRepo) Total(user_id uint, name string) (cate *Cate, err error) {
	err = r.db.Where("user_id", user_id).Where("name", name).Find(&cate).Error
	if err != nil {
		log.Printf("Error By Cate Total: %+v", err)
	}
	return
}

func (r *CateRepo) List(user_id uint) (names []string, err error) {
	err = r.db.Model(&Cate{}).Distinct().Where("user_id=?", user_id).Pluck("name", &names).Error
	if err != nil {
		log.Printf("Error By Cate List: %+v", err)
	}
	return names, err
}
