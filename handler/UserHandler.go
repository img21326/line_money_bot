package handler

import (
	"log"
	"moneybot/repo"
	"moneybot/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
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

	e.POST("v1/user/month/sum", handler.MonthOfSum)
}

type UserInputMonthOfSum struct {
	UserId string `json:"user_id"`
	Tag    string `json:"tag"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

func (h *UserHandler) MonthOfSum(c *gin.Context) {
	var user_input UserInputMonthOfSum
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("User Month of Sum BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), 1, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfMonth()}
	r, err := h.AccRepo.Sum(user.ID, s)
	if err != nil {
		log.Printf("User Month of Sum err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	type R struct {
		Total int64 `json:"total"`
	}
	c.JSON(200, R{Total: r})
}
