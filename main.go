package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LineId   string    `json:"line_id"`
	Accounts []Account `json:"accounts"`
	Total    int64     `json:"total"`
	tags     []Tag
}

type Account struct {
	gorm.Model
	Amount int   `json:"amount"`
	Tags   []Tag `json:"tags"`
	UserID uint
}

type Tag struct {
	gorm.Model
	Name      string `json:"name"`
	AccountID uint
	UserID    uint
	// Accounts []Account `json:"accounts",gorm:"many2many:account_tags;"`
}

type Search struct {
	User  User
	Tag   Tag
	Start time.Time
	End   time.Time
	Sum   string
}

type ApiTagSum struct {
	UserId string `json:"user_id"`
	Year   int    `json:"year"`
	Month  int    `json:"month"`
}

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// queryAPI := client.QueryAPI("my-org")
	// Setup HTTP Server for receiving requests from LINE platform
	db_host := os.Getenv("POSTGRES_HOST")
	db_pwd := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=postgres password=%s dbname=moneybot port=%s sslmode=disable TimeZone=Asia/Taipei", db_host, db_pwd, db_port)
	log.Println(dsn)
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	time.Sleep(10 * time.Second)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&User{}, &Account{}, &Tag{})

	Repo := NewRepo(db)

	reg_date, _ := regexp.Compile("20[0-2][0-9]/(0[1-9]|1[0-2])/(0[1-9]|[12][0-9]|3[01])-20[0-2][0-9]/(0[1-9]|1[0-2])/(0[1-9]|[12][0-9]|3[01])")

	r := gin.Default()
	r.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         "views",
		Extension:    ".tpl",
		Master:       "layouts/master",
		Funcs:        nil,
		DisableCache: true,
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{})
	})

	r.POST("/v1/tags/sum", func(c *gin.Context) {
		// fake
		var u ApiTagSum
		if err := c.BindJSON(&u); err != nil {
			log.Printf("Tags Sum BindJson err: %+v \n", err)
			c.AbortWithStatus(400)
			return
		}
		user := Repo.FindOrCreateUser(u.UserId)
		t := time.Date(u.Year, time.Month(u.Month), 1, 0, 0, 0, 0, time.Now().Location())
		search := &Search{User: user, Start: t, End: now.With(t).EndOfMonth()}
		r := Repo.ListTagsSum(search)
		c.JSON(200, r)
	})

	// http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
	r.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
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
					user := Repo.FindOrCreateUser(event.Source.UserID)
					// log.Printf("Find User: %+v", user)
					message_arr := strings.Fields(message.Text)
					message_arr[0] = strings.ToLower(message_arr[0])
					log.Printf("Message: %v", message_arr)

					if reg_date.MatchString(message_arr[0]) && len(message_arr) >= 2 {
						var search Search
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
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
								log.Print(err)
							}
							return
						}
						search.User = user
						search.End = end
						search.Start = start

						search.Sum = message_arr[1]
						if len(message_arr) >= 3 {
							search.Tag.Name = message_arr[2]
						}
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("總額為: %d", Repo.AdvanceSearch(&search)))).Do(); err != nil {
							log.Print(err)
						}
						return
					}
					// 如果輸入的在指令集裡面
					cmds := []string{"today", "今日", "month", "本月", "week", "本週", "year", "今年", "list", "列表"}
					if StringInSlice(message_arr[0], cmds) && len(message_arr) >= 2 {
						var search Search
						if message_arr[0] == "list" || message_arr[0] == "列表" {
							if message_arr[1] == "tag" || message_arr[1] == "tags" || message_arr[1] == "標籤" {
								search.User = user
								r := Repo.ListTags(&search)
								var s string
								for k, i := range r {
									if k == len(r)-1 {
										s += fmt.Sprintf("%s", i)
									} else {
										s += fmt.Sprintf("%s\n", i)
									}
								}
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(s)).Do(); err != nil {
									log.Print(err)
								}
								return
							}
						}

						//////////////////// 計算總和 ////////////////////////////
						var start time.Time
						var end time.Time
						search.User = user
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
						if len(message_arr) >= 3 {
							search.Tag.Name = message_arr[2]
						}
						search.Sum = message_arr[1]
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("總額為: %d", Repo.AdvanceSearch(&search)))).Do(); err != nil {
							log.Print(err)
						}
						return
					}

					if len(message_arr) == 1 {
						if message_arr[0] == "餘額" || message_arr[0] == "balance" {
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("目前餘額為: %d", user.Total))).Do(); err != nil {
								log.Print(err)
							}
							return
						}
					}

					// 如果開頭不是+或-
					// 則跳出
					if message_arr[0][0] != 43 && message_arr[0][0] != 45 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
							log.Print(err)
						}
						return
					}

					// 最多輸入兩行
					// 否則跳出
					if len(message_arr) >= 3 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
							log.Print(err)
						}
						return
					}

					// Create
					amount, _ := strconv.Atoi(message_arr[0])
					if amount != 0 {
						acc := Account{
							Amount: amount,
							UserID: user.ID,
						}
						if len(message_arr) > 1 {
							tags := strings.Split(message_arr[1], ",")
							for _, t := range tags {
								acc.Tags = append(acc.Tags, Tag{Name: t, UserID: user.ID})
							}
						} else {
							acc.Tags = append(acc.Tags, Tag{Name: "default", UserID: user.ID})
						}

						err := Repo.CreateAccountAndUpdateUser(&user, &acc)
						if err != nil {
							log.Print("Create Account Error: %+v", err)
							return
						}
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("目前餘額為: %d", user.Total))).Do(); err != nil {
						log.Print(err)
					}
					return
				}
			}
		}
	})

	// })
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	log.Printf("Server Start at Port: %s", os.Getenv("PORT"))
	// if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
	// 	log.Fatal(err)
	// }
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) FindOrCreateUser(line_id string) (user User) {
	err := r.db.Where("line_id=?", line_id).Find(&user).Error

	if user.ID == 0 || err == gorm.ErrRecordNotFound {
		user = User{LineId: line_id, Total: 0}
		r.db.Create(&user)
	}
	return user
}

// func (r *Repo) CreateAccount(account *Account) error {
// 	err := r.db.Create(&account).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *Repo) CreateAccountAndUpdateUser(user *User, account *Account) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user).Error; err != nil {
			return err
		}
		user.Total += int64(account.Amount)
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		return nil
	})
}

type AccountResult struct {
	Total int64
}

func (r *Repo) GetSumOfAccount(s *Search) int64 {
	var result AccountResult
	var err error
	w := r.db.Model(&Account{}).Select("user_id, sum(amount) as Total").Where("user_id=?", s.User.ID).Where("created_at>?", s.Start).Where("created_at<=?", s.End)
	switch s.Sum {
	case "+":
		w = w.Where("amount > 0")
	case "-":
		w = w.Where("amount < 0")
	case "sum":
	}
	err = w.Group("user_id").Find(&result).Error
	if err != nil {
		log.Fatalf("Get Sum Of Account Error: %+v", err)
	}
	return result.Total
}

func (r *Repo) GetSumOfTag(s *Search) int64 {
	var result AccountResult
	var tags []Tag
	r.db.Where("name = ?", s.Tag.Name).Where("user_id = ?", s.User.ID).Find(&tags)
	if len(tags) > 0 {
		var ids []uint
		for _, t := range tags {
			ids = append(ids, t.AccountID)
		}
		w := r.db.Model(&Account{}).Select("user_id, sum(amount) as Total").Where("user_id=?", s.User.ID).Where("created_at>?", s.Start).Where("created_at<=?", s.End)
		// if s.Plus {
		// 	w = w.Where("amount > 0")
		// } else {
		// 	w = w.Where("amount < 0")
		// }
		switch s.Sum {
		case "+":
			w = w.Where("amount > 0")
		case "-":
			w = w.Where("amount < 0")
		case "sum":
			w = w
		}
		w.Where("id IN ?", ids)
		err := w.Group("user_id").Find(&result).Error
		if err != nil {
			log.Fatalf("Advance Search With Tag error: %+v", err)
		}
		return result.Total
	} else {
		return 0
	}
}

func (r *Repo) AdvanceSearch(s *Search) int64 {
	if s.Tag.Name != "" {
		return r.GetSumOfTag(s)
	} else {
		return r.GetSumOfAccount(s)
	}
}

func (r *Repo) ListTags(s *Search) []string {
	// search -> User
	var names []string
	err := r.db.Model(&Tag{}).Distinct().Where("user_id=?", s.User.ID).Pluck("name", &names).Error
	if err != nil {
		log.Printf("Error By ListTags: %+v", err)
	}
	return names
}

type TagSum struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

func (r *Repo) ListTagsSum(s *Search) []TagSum {
	var rs []TagSum
	var t TagSum
	w := r.db.Model(&Tag{}).Select("tags.name, sum(accounts.amount) as Total").Joins("inner join accounts on accounts.id = tags.account_id")
	w = w.Where("tags.user_id", s.User.ID).Where("tags.created_at>?", s.Start).Where("tags.created_at<=?", s.End)
	rows, err := w.Group("tags.name").Rows()
	if err != nil {
		log.Printf("Error By ListTagsSum: %+v", err)
	}
	for rows.Next() {
		r.db.ScanRows(rows, &t)
		rs = append(rs, t)
	}
	return rs
}
