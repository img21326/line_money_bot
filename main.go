package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/jinzhu/now"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := influxdb2.NewClient(os.Getenv("INFLUXDB_HOST"), os.Getenv("INFLUXDB_TOKEN"))
	writeAPI := client.WriteAPIBlocking("my-org", "my-bucket")
	queryAPI := client.QueryAPI("my-org")
	// queryAPI := client.QueryAPI("my-org")
	// Setup HTTP Server for receiving requests from LINE platform
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
					message_arr := strings.Fields(message.Text)
					log.Printf("Message: %v", message_arr)
					if len(message_arr) == 1 {
						if message_arr[0] == "餘額" {
							_start := "2021-01-01T00:00:01Z"
							query_str := fmt.Sprintf(`from(bucket:"my-bucket")|> range(start: %s) |> filter(fn: (r) => r._measurement == "account_book") |> filter(fn: (r) => r["unit"] == "%s") |> filter(fn: (r) => r["_field"] == "amount") |> sum()`, _start, event.Source.UserID)
							// log.Printf("%v", query_str)
							result, err := queryAPI.Query(context.Background(), query_str)
							if err != nil {
								log.Printf("error from 餘額: %v", err)
							}
							var amount int64
							for result.Next() {
								amount = result.Record().Value().(int64)
							}
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("餘額為: %d", amount))).Do(); err != nil {
								log.Print(err)
								return
							}
						}
						if message_arr[0] == "今日花費" {
							now := time.Now()
							year, month, day := now.Date()
							_start := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
							_end := time.Date(year, month, day, 23, 59, 59, 59, now.Location())
							query_str := fmt.Sprintf(`from(bucket:"my-bucket")|> range(start: %d, stop: %d) |> filter(fn: (r) => r._measurement == "account_book") |> filter(fn: (r) => r["unit"] == "%s") |> filter(fn: (r) => r["_field"] == "amount") |> filter(fn: (r) => r._value < 0) |> sum()`, _start.Unix(), _end.Unix(), event.Source.UserID)
							// log.Printf("%v", query_str)
							result, err := queryAPI.Query(context.Background(), query_str)
							if err != nil {
								log.Printf("error from 今日花費: %v", err)
							}
							var amount int64
							for result.Next() {
								amount = result.Record().Value().(int64)
							}
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("本日花費為: %f", math.Abs(float64(amount))))).Do(); err != nil {
								log.Print(err)
								return
							}

						}
						if message_arr[0] == "本週花費" {
							_start := now.BeginningOfWeek()
							_end := now.EndOfWeek()
							query_str := fmt.Sprintf(`from(bucket:"my-bucket")|> range(start: %d, stop: %d) |> filter(fn: (r) => r._measurement == "account_book") |> filter(fn: (r) => r["unit"] == "%s") |> filter(fn: (r) => r["_field"] == "amount") |> filter(fn: (r) => r._value < 0) |> sum()`, _start.Unix(), _end.Unix(), event.Source.UserID)
							// log.Printf("%v", query_str)
							result, err := queryAPI.Query(context.Background(), query_str)
							if err != nil {
								log.Printf("error from 本週花費: %v", err)
							}
							var amount int64
							for result.Next() {
								amount = result.Record().Value().(int64)
							}
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("本週花費: %f", math.Abs(float64(amount))))).Do(); err != nil {
								log.Print(err)
								return
							}

						}
						if message_arr[0] == "本月花費" {
							_start := now.BeginningOfMonth()
							_end := now.EndOfMonth()
							query_str := fmt.Sprintf(`from(bucket:"my-bucket")|> range(start: %d, stop: %d) |> filter(fn: (r) => r._measurement == "account_book") |> filter(fn: (r) => r["unit"] == "%s") |> filter(fn: (r) => r["_field"] == "amount") |> filter(fn: (r) => r._value < 0) |> sum()`, _start.Unix(), _end.Unix(), event.Source.UserID)
							// log.Printf("%v", query_str)
							result, err := queryAPI.Query(context.Background(), query_str)
							if err != nil {
								log.Printf("error from 本月花費: %v", err)
							}
							var amount int64
							for result.Next() {
								amount = result.Record().Value().(int64)
							}
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(fmt.Sprintf("本月花費: %f", math.Abs(float64(amount))))).Do(); err != nil {
								log.Print(err)
								return
							}
						}
					}
					if message_arr[0][0] != 43 && message_arr[0][0] != 45 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
							log.Print(err)
							return
						}
					}
					if len(message_arr) >= 3 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式輸入錯誤")).Do(); err != nil {
							log.Print(err)
							return
						}
					}
					amount, _ := strconv.Atoi(message_arr[0])

					var p *write.Point
					if len(message_arr) == 2 {
						p = influxdb2.NewPointWithMeasurement("account_book").
							AddTag("unit", event.Source.UserID).
							AddField("amount", amount).
							AddField("desc", message_arr[1]).
							SetTime(time.Now())
					} else {
						p = influxdb2.NewPointWithMeasurement("account_book").
							AddTag("unit", event.Source.UserID).
							AddField("amount", amount).
							SetTime(time.Now())
					}

					writeAPI.WritePoint(context.Background(), p)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("已存入資料庫")).Do(); err != nil {
						log.Print(err)
						return
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
