package utils

import "time"

type NameOfSum struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

type DaySum struct {
	Day   time.Time `json:"day"`
	Total int64     `json:"total"`
}
