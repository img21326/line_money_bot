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
	handler := UserHandler{
		UserRepo: u,
		AccRepo:  a,
	}
	e.POST("/v1/acc/insert", handler.InsertAccount)
}

type UserInputMonthOfSum struct {
	UserId string `json:"user_id"`
	Tag    string `json:"tag"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

type UserInputInsertAccount struct {
	Amount int    `json:"amount"`
	Cate   string `json:"cate"`
	Date   string `json:"date"`
}

func (h *UserHandler) InsertAccount(c *gin.Context) {

}
