package utils

import (
	"time"

	"gorm.io/gorm"
)

type Select struct {
	Start time.Time
	End   time.Time
	Cate  string
	Tag   string
	Sum   string
}

func SelectSum(w *gorm.DB, s Select) *gorm.DB {
	switch s.Sum {
	case "+":
		w = w.Where("amount > 0")
	case "-":
		w = w.Where("amount < 0")
	case "sum":
	}
	return w
}
