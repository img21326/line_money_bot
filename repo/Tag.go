package repo

import (
	"log"
	"moneybot/utils"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name      string `json:"name"`
	AccountID uint
	UserID    uint
}

type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{db}
}

func (r *TagRepo) NameOfSum(user_id uint, tag_name string, s utils.Select) (i int64, err error) {
	type Result struct {
		Total int64
	}
	var result Result
	var tags []Tag
	r.db.Where("name = ?", tag_name).Where("user_id = ?", user_id).Find(&tags)
	if len(tags) > 0 {
		var ids []uint
		for _, t := range tags {
			ids = append(ids, t.AccountID)
		}
		w := r.db.Model(&Account{}).Select("user_id, sum(amount) as Total").Where("user_id=?", user_id).Where("created_at>?", s.Start).Where("created_at<=?", s.End)
		w = utils.SelectSum(w, s)
		w = w.Where("id IN ?", ids)
		err := w.Group("user_id").Find(&result).Error
		if err != nil {
			log.Fatalf("Advance Search With Tag error: %+v", err)
		}
		return result.Total, err
	} else {
		return 0, nil
	}
}

func (r *TagRepo) List(user_id uint) (names []string, err error) {
	err = r.db.Model(&Tag{}).Distinct().Where("user_id=?", user_id).Pluck("name", &names).Error
	if err != nil {
		log.Printf("Error By ListTags: %+v", err)
	}
	return names, err
}

func (r *TagRepo) ListNameOfSum(user_id uint, s utils.Select) (tags_sum []utils.NameOfSum, err error) {
	var t utils.NameOfSum
	w := r.db.Model(&Tag{}).Select("tags.name, sum(accounts.amount) as Total").Joins("inner join accounts on accounts.id = tags.account_id")
	w = w.Where("tags.user_id", user_id).Where("tags.created_at>?", s.Start).Where("tags.created_at<=?", s.End)
	rows, err := w.Group("tags.name").Order("Total").Rows()
	if err != nil {
		log.Printf("Error By ListTagsSum: %+v", err)
	}
	for rows.Next() {
		r.db.ScanRows(rows, &t)
		tags_sum = append(tags_sum, t)
	}
	return
}

func (r *TagRepo) ListDayOfSum(user_id uint, tag_name string, s utils.Select) (day_sum []utils.DaySum, err error) {
	var rr utils.DaySum
	w := r.db.Model(&Tag{}).Select("date_trunc('day',accounts.created_at) as \"Day\", sum(accounts.amount) as \"Total\"").Joins("inner join accounts on accounts.id = tags.account_id")
	w = w.Where("tags.user_id", user_id).Where("tags.created_at>?", s.Start).Where("tags.created_at<=?", s.End)
	rows, err := w.Where("tags.name", tag_name).Group("Day").Order("Day").Rows()
	for rows.Next() {
		r.db.ScanRows(rows, &rr)
		day_sum = append(day_sum, rr)
	}
	return
}
