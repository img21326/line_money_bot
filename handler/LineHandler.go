package handler

import (
	"fmt"
	"log"
	"moneybot/repo"
	"moneybot/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineHandler struct {
	LineClient  *linebot.Client
	UserRepo    repo.UserRepo
	TagRepo     repo.TagRepo
	AccountRepo repo.AccountRepo
	CateRepo    repo.CateRepo
}

func NewLineHandler(e *gin.Engine, l *linebot.Client, u repo.UserRepo, t repo.TagRepo, a repo.AccountRepo, c repo.CateRepo) {
	handler := &LineHandler{
		LineClient:  l,
		UserRepo:    u,
		TagRepo:     t,
		AccountRepo: a,
		CateRepo:    c,
	}

	e.POST("line/callback", handler.CallBack)
}

func (h *LineHandler) CallBack(c *gin.Context) {
	reg_date, _ := regexp.Compile("20[0-2][0-9]/(0[1-9]|1[0-2])/(0[1-9]|[12][0-9]|3[01])-20[0-2][0-9]/(0[1-9]|1[0-2])/(0[1-9]|[12][0-9]|3[01])")

	events, err := h.LineClient.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			// w.WriteHeader(400)
			c.AbortWithStatus(400)
		} else {
			// w.WriteHeader(500)
			c.AbortWithStatus(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Printf("UserID: %v", event.Source.UserID)
				user := h.UserRepo.FindOrCreateUser(event.Source.UserID)
				// log.Printf("Find User: %+v", user)
				message_arr := strings.Fields(message.Text)
				message_arr[0] = strings.ToLower(message_arr[0])
				log.Printf("Message: %v", message_arr)

				if reg_date.MatchString(message_arr[0]) && len(message_arr) >= 2 {
					var search utils.Select
					var start time.Time
					var end time.Time
					time_arr := strings.Split(message_arr[0], "-")
					log.Printf("%+v\n", time_arr)
					layout := "2006/01/02"
					start, _ = time.Parse(layout, time_arr[0])
					end, _ = time.Parse(layout, time_arr[1])
					end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 59, end.Location())

					// 結束日小於開始日
					if end.Before(start) {
						if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
							log.Print(err)
						}
						return
					}
					search.End = end
					search.Start = start

					search.Sum = message_arr[1]
					if len(message_arr) >= 3 {
						tag_name := message_arr[2]
						d, err := h.TagRepo.NameOfSum(user.ID, tag_name, search)
						if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("總額為: %d", d))).Do(); err != nil {
							log.Print(err)
						}
						return
					}
					d, err := h.AccountRepo.Sum(user.ID, search)
					if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("總額為: %d", d))).Do(); err != nil {
						log.Print(err)
					}
					return
				}
				// 如果輸入的在指令集裡面
				cmds := []string{"today", "今日", "month", "本月", "week", "本週", "year", "今年", "list", "列表"}
				if utils.StringInSlice(message_arr[0], cmds) && len(message_arr) >= 2 {
					var search utils.Select
					if message_arr[0] == "list" || message_arr[0] == "列表" {
						if message_arr[1] == "tag" || message_arr[1] == "tags" || message_arr[1] == "標籤" {
							r, err := h.TagRepo.List(user.ID)
							var s string
							for k, i := range r {
								if k == len(r)-1 {
									s += fmt.Sprintf("%s", i)
								} else {
									s += fmt.Sprintf("%s\n", i)
								}
							}
							if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(s)).Do(); err != nil {
								log.Print(err)
							}
							return
						}

						if message_arr[1] == "cate" || message_arr[1] == "cates" || message_arr[1] == "類別" {
							r, err := h.CateRepo.List(user.ID)
							var s string
							for k, i := range r {
								if k == len(r)-1 {
									s += fmt.Sprintf("%s", i)
								} else {
									s += fmt.Sprintf("%s\n", i)
								}
							}
							if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(s)).Do(); err != nil {
								log.Print(err)
							}
							return
						}
					}

					//////////////////// 計算總和 ////////////////////////////
					var start time.Time
					var end time.Time
					switch message_arr[0] {
					case "today", "今日":
						now := time.Now()
						year, month, day := now.Date()
						start = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
						end = time.Date(year, month, day, 23, 59, 59, 59, now.Location())
					case "week", "本週":
						start = now.BeginningOfWeek()
						end = now.EndOfWeek()
					case "month", "本月":
						start = now.BeginningOfMonth()
						end = now.EndOfMonth()
					case "year", "今年":
						start = now.BeginningOfYear()
						end = now.EndOfYear()
					}
					search.Start = start
					search.End = end
					search.Sum = message_arr[1]
					var d int64
					if len(message_arr) >= 3 {
						tag_name := message_arr[2]
						d, _ = h.TagRepo.NameOfSum(user.ID, tag_name, search)
					} else {
						d, _ = h.AccountRepo.Sum(user.ID, search)
					}
					if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("總額為: %d", d))).Do(); err != nil {
						log.Print(err)
					}
					return
				}

				if message_arr[0] == "餘額" || message_arr[0] == "balance" {
					var cate_name string
					if len(message_arr) == 1 {
						cate_name = "default"
					} else {
						cate_name = message_arr[1]
					}
					r, err := h.CateRepo.Total(user.ID, cate_name)
					if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("目前 %s 餘額為: %d", cate_name, r.Total))).Do(); err != nil {
						log.Print(err)
					}
					return
				}

				// 如果開頭不是+或-
				// 則跳出
				if message_arr[0][0] != 43 && message_arr[0][0] != 45 {
					if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
						log.Print(err)
					}
					return
				}

				// 最多輸入兩行
				// 否則跳出
				if len(message_arr) >= 3 {
					if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
						log.Print(err)
					}
					return
				}

				// Create
				amount, _ := strconv.Atoi(message_arr[0])
				var cate repo.Cate
				if amount != 0 {
					acc := repo.Account{
						Amount: amount,
						UserID: user.ID,
					}
					if len(message_arr) > 1 {
						tags := strings.Split(message_arr[1], ",")
						cate = repo.Cate{Name: tags[0]}
						for _, t := range tags[1:] {
							acc.Tags = append(acc.Tags, repo.Tag{Name: t, UserID: user.ID})
						}
					} else {
						cate = repo.Cate{Name: "default"}
					}

					err := h.UserRepo.CreateAccountAndUpdateUser(user, &acc, &cate)
					if err != nil {
						log.Print("Create Account Error: %+v", err)
						return
					}
				} else {
					cate.Name = "default"
				}
				r, err := h.CateRepo.Total(user.ID, cate.Name)
				if _, err = h.LineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("目前 %s 餘額為: %d", cate.Name, r.Total))).Do(); err != nil {
					log.Print(err)
				}
				return
			}
		}
	}
}
