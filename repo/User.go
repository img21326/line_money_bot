package repo

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LineId   string    `json:"line_id"`
	Accounts []Account `json:"accounts"`
	Total    int64     `json:"total"`
	tags     []Tag
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
		user = &User{LineId: line_id, Total: 0}
		r.db.Create(user)
	}
	return user
}

func (r *UserRepo) CreateAccountAndUpdateUser(user *User, account *Account) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user).Error; err != nil {
			return err
		}
		user.Total += int64(account.Amount)
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		return nil
	})
}
