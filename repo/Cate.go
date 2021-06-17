package repo

import (
	"log"
	"moneybot/utils"

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

func (r *CateRepo) GetCate(user_id uint, s utils.Select) (cate *Cate, err error) {
	w := r.db.Where("user_id", user_id).Where("created_at > ?", s.Start).Where("created_at <= ?", s.End)
	err = w.Where("name", s.Cate).Find(&cate).Error
	if err != nil {
		log.Printf("Error By Get Cate: %+v", err)
	}
	return
}

func (r *CateRepo) Total(user_id uint, s utils.Select) (cate *Cate, err error) {
	err = r.db.Where("user_id", user_id).Where("created_at > ?", s.Start).Where("created_at <= ?", s.End).Where("name", s.Cate).Find(&cate).Error
	if err != nil {
		log.Printf("Error By Cate Total: %+v", err)
	}
	return
}

func (r *CateRepo) List(user_id uint, s utils.Select) (names []string, err error) {
	w := r.db.Model(&Cate{}).Where("created_at > ?", s.Start).Where("created_at <= ?", s.End).Distinct()
	err = w.Pluck("name", &names).Error
	if err != nil {
		log.Printf("Error By Cate List: %+v", err)
	}
	return names, err
}
