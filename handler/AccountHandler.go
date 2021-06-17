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

type AccHandler struct {
	UserRepo repo.UserRepo
	AccRepo  repo.AccountRepo
	CateRepo repo.CateRepo
}

func NewAccHandler(e *gin.Engine, u repo.UserRepo, a repo.AccountRepo, c repo.CateRepo) {
	handler := AccHandler{
		UserRepo: u,
		AccRepo:  a,
		CateRepo: c,
	}

	e.GET("/acc/list/cate/sum", handler.ShowListCateOfSumPage)
	e.GET("/acc/list/day/sum", handler.ShowListDayOfSumPage)

	e.POST("/v1/acc/month/sum", handler.MonthOfSum)
	e.POST("v1/acc/month/list/cate", handler.ListCate)
	e.POST("v1/acc/month/list/cate/sum", handler.ListMonthOfCateSum)
	e.POST("/v1/acc/days/list/sum", handler.ListDayOfSum)
	e.POST("/v1/acc/day/list/info", handler.ListDayOfInfo)
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

type AccInputListMonthOfCateSum struct {
	UserId string `json:"user_id"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

func (h *AccHandler) ShowListCateOfSumPage(c *gin.Context) {
	c.HTML(http.StatusOK, "listcateofsum", gin.H{
		"liff_id": os.Getenv("LIFF_LISTCATEOFSUM"),
	})
}

func (h *AccHandler) ListCate(c *gin.Context) {
	var user_input AccInputListMonthOfCateSum
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("Acc ListMonthOfCateSum BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), 1, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfMonth()}
	r, err := h.CateRepo.List(user.ID, s)
	if err != nil {
		log.Printf("Acc ListMonthOfCateSum err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, r)
}

func (h *AccHandler) ListMonthOfCateSum(c *gin.Context) {
	var user_input AccInputListMonthOfCateSum
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("Acc ListMonthOfCateSum BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), 1, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfMonth()}
	r, err := h.AccRepo.ListMonthOfCateSum(user.ID, s)
	if err != nil {
		log.Printf("Acc ListMonthOfCateSum err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, r)
}

func (h *AccHandler) ShowListDayOfSumPage(c *gin.Context) {
	q := c.Request.URL.Query()
	var cate string
	if len(q["cate"]) > 0 {
		cate = q["cate"][0]
	}
	c.HTML(http.StatusOK, "listdayofsum", gin.H{
		"liff_id": os.Getenv("LIFF_LISTDAYOFSUM"),
		"cate":    cate,
	})
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

type AccInputDayOfInfo struct {
	UserId string `json:"user_id"`
	Cate   string `json:"cate"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Day    int    `json:"day"`
}

func (h *AccHandler) ListDayOfInfo(c *gin.Context) {
	var user_input AccInputDayOfInfo
	if err := c.BindJSON(&user_input); err != nil {
		log.Printf("Acc ListDayOfInfo BindJson err: %+v \n", err)
		c.AbortWithStatus(400)
		return
	}
	user := h.UserRepo.FindOrCreateUser(user_input.UserId)
	t := time.Date(user_input.Year, time.Month(user_input.Month), user_input.Day, 0, 0, 0, 0, time.Now().Location())
	s := utils.Select{Start: t, End: now.With(t).EndOfDay(), Cate: user_input.Cate}
	accs, err := h.AccRepo.ListDayOfInfo(user.ID, s)
	if err != nil {
		log.Printf("Acc ListDayOfInfo err: %+v \n", err)
		c.AbortWithStatus(500)
		return
	}
	type RetData struct {
		CreatedAt time.Time `json:"created_at"`
		Amount    int       `json:"amount"`
		Tags      []string  `json:"tags"`
	}
	var ret_dat []RetData
	for i := range accs {
		var x RetData
		a := accs[i]
		x.Amount = a.Amount
		x.CreatedAt = a.CreatedAt
		var n []string
		for z := range a.Tags {
			n = append(n, a.Tags[z].Name)
		}
		x.Tags = n
		ret_dat = append(ret_dat, x)
	}
	c.JSON(200, ret_dat)
}
