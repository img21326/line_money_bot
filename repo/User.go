package repo

import (
	"time"

	"github.com/jinzhu/now"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LineId   string    `json:"line_id"`
	Accounts []Account `json:"accounts"`
	Tags     []Tag     `json:"tags"`
	Cates    []Cate    `json:"cates"`
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

func (r *UserRepo) FindOrCreateUser(line_id string) (user *User) {
	err := r.db.Where("line_id=?", line_id).Find(&user).Error

	if user.ID == 0 || err == gorm.ErrRecordNotFound {
		user = &User{LineId: line_id}
		r.db.Create(user)
	}
	return user
}

type CreateAccount struct {
	Account Account
	Cate    string
	Date    time.Time
}

func (r *UserRepo) CreateAccountAndUpdateCate(user *User, ac *CreateAccount) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user).Error; err != nil {
			return err
		}
		// user.Total += int64(account.Amount)
		var fcate Cate
		var err error
		if ac.Date.IsZero() {
			ac.Date = time.Now()
		}
		tx.Model(&Cate{}).Where("user_id", user.ID).Where("name", ac.Cate).Where("created_at > ?", now.With(ac.Date).BeginningOfMonth()).Where("created_at <= ?", now.With(ac.Date).EndOfMonth()).Find(&fcate)
		if fcate.ID == 0 {
			// 將上個月的匯入 或是 新建新的一筆
			var pre_cate Cate
			d := now.With(ac.Date).BeginningOfMonth().AddDate(0, 0, -1)
			tx.Model(&Cate{}).Where("user_id", user.ID).Where("name", ac.Cate).Where("created_at > ?", now.With(d).BeginningOfMonth()).Where("created_at <= ?", now.With(d).EndOfMonth()).Find(&pre_cate)
			if pre_cate.ID != 0 {
				fcate = Cate{UserID: user.ID, Name: ac.Cate, Total: int64(int(pre_cate.Total) + ac.Account.Amount)}
			} else {
				fcate = Cate{UserID: user.ID, Name: ac.Cate, Total: int64(ac.Account.Amount)}
			}
			fcate.CreatedAt = ac.Date
			err = tx.Create(&fcate).Error
		} else {
			fcate.Total += int64(ac.Account.Amount)
			err = tx.Save(&fcate).Error
		}
		if err != nil {
			return err
		}
		//
		ac.Account.CateID = fcate.ID
		ac.Account.CreatedAt = ac.Date
		if err := tx.Create(&ac.Account).Error; err != nil {
			return err
		}
		return nil
	})
}
