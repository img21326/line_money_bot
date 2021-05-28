package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	fmt.Print(dsn)
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&User{}, &Account{}, &Tag{})

	Repo := NewRepo(db)

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					log.Printf("UserID: %v", event.Source.UserID)
					user := Repo.FindOrCreateUser(event.Source.UserID)
					log.Printf("Find User: %+v", user)
					message_arr := strings.Fields(message.Text)
					log.Printf("Message: %v", message_arr)
					if len(message_arr) == 1 {
						if message_arr[0] == "餘額" {
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("目前餘額為: %d", user.Total))).Do(); err != nil {
								log.Print(err)
							}
							return
						}
						if message_arr[0] == "今日花費" {
							now := time.Now()
							year, month, day := now.Date()
							start := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
							end := time.Date(year, month, day, 23, 59, 59, 59, now.Location())
							total := Repo.GetSumOfAccount(&user, start, end)
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("今日花費: %d", total))).Do(); err != nil {
								log.Print(err)
							}
							return
						}
						if message_arr[0] == "本週花費" {
							start := now.BeginningOfWeek()
							end := now.EndOfWeek()
							total := Repo.GetSumOfAccount(&user, start, end)
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("今日花費: %d", total))).Do(); err != nil {
								log.Print(err)
							}
							return
						}
						if message_arr[0] == "本月花費" {
							start := now.BeginningOfMonth()
							end := now.EndOfMonth()
							total := Repo.GetSumOfAccount(&user, start, end)
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("今日花費: %d", total))).Do(); err != nil {
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
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	log.Printf("Server Start at Port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
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

func (r *Repo) GetSumOfAccount(user *User, start time.Time, end time.Time) int64 {
	type Result struct {
		Total int64
	}
	var result Result
	err := r.db.Model(&Account{}).Select("user_id, sum(amount) as Total").Where("user_id=?", user.ID).Where("amount < 0").Where("created_at>?", start).Where("created_at<=?", end).Group("user_id").Find(&result).Error
	if err != nil {
		log.Fatalf("Get Sum Of Account Error: %+v", err)
	}
	return result.Total
}
