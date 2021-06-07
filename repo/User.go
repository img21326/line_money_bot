package repo

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LineId   string    `json:"line_id"`
	Accounts []Account `json:"accounts"`
	tags     []Tag     `json:"tags"`
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

func (r *UserRepo) CreateAccountAndUpdateUser(user *User, account *Account, cate *Cate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user).Error; err != nil {
			return err
		}
		// user.Total += int64(account.Amount)
		var fcate Cate
		var err error
		tx.Model(&Cate{}).Where("user_id", user.ID).Where("name", cate.Name).Find(&fcate)
		if fcate.ID == 0 {
			fcate = Cate{UserID: user.ID, Name: cate.Name, Total: int64(account.Amount)}
			err = tx.Create(&fcate).Error
		} else {
			fcate.Total += int64(account.Amount)
			err = tx.Save(&fcate).Error
		}
		if err != nil {
			return err
		}
		//
		account.CateID = fcate.ID
		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		return nil
	})
}
