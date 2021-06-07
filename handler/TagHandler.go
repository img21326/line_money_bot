package handler

import (
	"log"
	"moneybot/repo"
	"moneybot/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

type TagHandler struct {
	UserRepo repo.UserRepo
	TagRepo  repo.TagRepo
}

func NewTagHandler(e *gin.Engine, u repo.UserRepo, t repo.TagRepo) {
	handler := TagHandler{
		UserRepo: u,
		TagRepo:  t,
	}

	e.GET("/tag/list/name/sum", handler.ShowTagsOfSumPage)
	e.POST("/v1/tags/month/sum", handler.MonthOfSum)
}

func (h *TagHandler) ShowTagsOfSumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "tagsum", gin.H{
		"liff_id": os.Getenv("LIFF_TAGSUM"),
	})
}

type TagInputMonthOfSum struct {
	UserId string `json:"user_id"`
	Tag    string `json:"tag"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

func (h *TagHandler) MonthOfSum(c *gin.Context) {
	var user_input TagInputMonthOfSum
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("Tag Month of Sum BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), 1, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfMonth()}
	r, err := h.TagRepo.ListNameOfSum(user.ID, s)
	if err != nil {
		log.Printf("Tags Month of Sum List err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, r)
}
