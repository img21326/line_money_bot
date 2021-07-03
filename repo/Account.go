package repo

import (
	"log"

	"moneybot/utils"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Amount int   `json:"amount"`
	Tags   []Tag `json:"tags"`
	UserID uint
	CateID uint
}

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{db}
}

func (r *AccountRepo) Sum(user_id uint, s utils.Select) (i int64, err error) {
	type Result struct {
		Total int64
	}
	var result Result
	w := r.db.Model(&Account{}).Select("user_id, sum(amount) as Total").Where("user_id=?", user_id).Where("created_at>?", s.Start).Where("created_at<=?", s.End)
	w = utils.SelectSum(w, s)
	err = w.Group("user_id").Find(&result).Error
	if err != nil {
		log.Fatalf("Get Sum Of Account Error: %+v", err)
	}
	return result.Total, err
}

func (r *AccountRepo) ListMonthOfCateSum(user_id uint, s utils.Select) (name_sum []utils.NameOfSum, err error) {
	var rr utils.NameOfSum
	// w := r.db.Model(&Account{}).Select("cate as name, sum(amount) as total")
	w := r.db.Model(&Account{}).Select("cates.name as name, sum(accounts.amount) as Total").Joins("inner join cates on accounts.cate_id = cates.id")
	w = w.Where("accounts.created_at > ?", s.Start).Where("accounts.created_at <= ?", s.End)
	rows, err := w.Group("name").Rows()
	if err != nil {
		log.Printf("Error By ListMonthCateSum: %+v", err)
	}
	for rows.Next() {
		r.db.ScanRows(rows, &rr)
		name_sum = append(name_sum, rr)
	}
	return

}

func (r *AccountRepo) ListDayOfSum(user_id uint, s utils.Select) (day_sum []utils.DaySum, err error) {
	var rr utils.DaySum
	w := r.db.Model(&Account{}).Select("date_trunc('day',created_at) as \"Day\", sum(amount) as \"Total\"")
	w = w.Where("created_at > ?", s.Start).Where("created_at <= ?", s.End)
	if s.Cate != "" {
		var cate Cate
		err = r.db.Model(&Cate{}).Where("name", s.Cate).Where("created_at > ?", s.Start).Where("created_at <= ?", s.End).Find(&cate).Error
		log.Printf("%+v \n", cate)
		if err == nil || cate.ID != 0 {
			w = w.Where("cate_id = ?", cate.ID)
		} else {
			log.Printf("ListDayOfSum Find Cate Err: %+v", err)
		}
	}
	w = w.Group("Day")
	rows, err := w.Rows()
	for rows.Next() {
		r.db.ScanRows(rows, &rr)
		day_sum = append(day_sum, rr)
	}
	return day_sum, err
}

func (r *AccountRepo) ListDayOfInfo(user_id uint, s utils.Select) (accs []Account, err error) {
	w := r.db.Where("created_at > ?", s.Start).Where("created_at <= ?", s.End)
	if s.Cate != "" {
		var cate Cate
		err = r.db.Model(&Cate{}).Where("name", s.Cate).Where("created_at > ?", now.With(s.Start).BeginningOfMonth()).Where("created_at <= ?", now.With(s.Start).EndOfMonth()).Find(&cate).Error
		log.Printf("%+v \n", cate)
		if err == nil || cate.ID != 0 {
			w = w.Where("cate_id = ?", cate.ID)
		} else {
			log.Printf("ListDayOfSum Find Cate Err: %+v", err)
		}
	}
	w = w.Where("user_id = ?", user_id)
	err = w.Preload("Tags").Find(&accs).Error
	return
}
