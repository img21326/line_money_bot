package handler

import (
	"log"
	"moneybot/repo"
	"moneybot/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

type AccHandler struct {
	UserRepo repo.UserRepo
	AccRepo  repo.AccountRepo
}

func NewAccHandler(e *gin.Engine, u repo.UserRepo, a repo.AccountRepo) {
	handler := AccHandler{
		UserRepo: u,
		AccRepo:  a,
	}

	e.POST("/v1/acc/month/sum", handler.MonthOfSum)
	e.POST("/v1/acc/days/list/sum", handler.ListDayOfSum)
}

type AccInputDayOfSum struct {
	UserId string `json:"user_id"`
	Cate   string `json:"cate"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

func (h *AccHandler) MonthOfSum(c *gin.Context) {
	var user_input UserInputMonthOfSum
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("Acc Month of Sum BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), 1, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfMonth()}
	r, err := h.AccRepo.Sum(user.ID, s)
	if err != nil {
		log.Printf("Acc Month of Sum err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	type R struct {
		Total int64 `json:"total"`
	}
	c.JSON(200, R{Total: r})
}

func (h *AccHandler) ListDayOfSum(c *gin.Context) {
	var user_input AccInputDayOfSum
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("Acc Day of Sum BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), 1, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfMonth(), Cate: user_input.Cate}
	r, err := h.AccRepo.ListDayOfSum(user.ID, s)
	if err != nil {
		log.Printf("Acc Day of Sum err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, r)
}
