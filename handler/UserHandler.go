package handler

import (
	"log"
	"moneybot/repo"
	"time"

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
	UserId string   `json:"user_id"`
	Amount int      `json:"amount"`
	Cate   string   `json:"cate"`
	Date   string   `json:"date"`
	Tags   []string `json:"tags"`
}

func (h *UserHandler) InsertAccount(c *gin.Context) {
	var user_input UserInputInsertAccount
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("User InsertAccount BindJson err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	if user.ID == 0 {
		log.Printf("Not found User with userID: %s", user_input.UserId)
		c.AbortWithStatus(500)
		return
	}
	layout := "2006-01-02"
	t, _ := time.Parse(layout, user_input.Date)
	create_acc := &repo.CreateAccount{
		Account: repo.Account{Amount: user_input.Amount, UserID: user.ID},
		Cate:    user_input.Cate,
		Date:    t,
	}
	for _, t := range user_input.Tags {
		create_acc.Account.Tags = append(create_acc.Account.Tags, repo.Tag{Name: t, UserID: user.ID})
	}
	err := h.UserRepo.CreateAccountAndUpdateCate(user, create_acc)
	if err != nil {
		log.Printf("User InsertAccount err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	c.JSON(200, true)
}
