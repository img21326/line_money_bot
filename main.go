package main

import (
	"fmt"
	"log"
	"moneybot/handler"
	"moneybot/repo"
	"os"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	db_host := os.Getenv("POSTGRES_HOST")
	db_pwd := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=postgres password=%s dbname=moneybot port=%s sslmode=disable TimeZone=Asia/Taipei", db_host, db_pwd, db_port)
	log.Println(dsn)
	// time.Sleep(10 * time.Second)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&repo.User{}, &repo.Account{}, &repo.Tag{})

	r := gin.Default()
	r.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:         "views",
		Extension:    ".tpl",
		Master:       "layouts/master",
		Funcs:        nil,
		DisableCache: true,
	})

	repo_acc := repo.NewAccountRepo(db)
	repo_tag := repo.NewTagRepo(db)
	repo_user := repo.NewUserRepo(db)

	handler.NewLineHandler(r, bot, *repo_user, *repo_tag, *repo_acc)
	handler.NewTagHandler(r, *repo_user, *repo_tag)
	handler.NewUserHandler(r, *repo_user, *repo_acc)

	// r.POST("/v1/days/sum", func(c *gin.Context) {
	// 	var u ApiSum
	// 	if err := c.BindJSON(&u); err != nil {
	// 		log.Printf("Tags Sum BindJson err: %+v \n", err)
	// 		c.AbortWithStatus(400)
	// 		return
	// 	}
	// 	user := Repo.FindOrCreateUser(u.UserId)
	// 	t := time.Date(u.Year, time.Month(u.Month), 1, 0, 0, 0, 0, time.Now().Location())
	// 	search := &Search{User: user, Start: t, End: now.With(t).EndOfMonth()}
	// 	if u.Tag != "" {
	// 		search.Tag = Tag{Name: u.Tag}
	// 	}
	// 	r := Repo.DayOfSum(search)
	// 	c.JSON(200, r)
	// })

	log.Printf("Server Start at Port: %s", os.Getenv("PORT"))
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
