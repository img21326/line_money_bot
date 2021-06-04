package handler

import (
	"moneybot/repo"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepo repo.UserRepo
	AccRepo  repo.AccountRepo
}

func NewUserHandler(e *gin.Engine, u repo.UserRepo, a repo.AccountRepo) {
	_ = UserHandler{
		UserRepo: u,
		AccRepo:  a,
	}
}

type UserInputMonthOfSum struct {
	UserId string `json:"user_id"`
	Tag    string `json:"tag"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}
