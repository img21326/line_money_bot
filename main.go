package main

import (
	"fmt"
	"log"
	"moneybot/handler"
	"moneybot/repo"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile("env.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Read ENV file err: %+v", err)
	}
}

func main() {
	CHANNEL_SECRET := viper.GetString("line.secret")
	CHANNEL_TOKEN := viper.GetString("line.token")
	bot, err := linebot.New(
		CHANNEL_SECRET,
		CHANNEL_TOKEN,
	)
	if err != nil {
		log.Fatal(err)
	}

	db_host := viper.GetString("db.host")
	db_pwd := viper.GetString("db.password")
	db_port := viper.GetString("db.port")
	dsn := fmt.Sprintf("host=%s user=postgres password=%s dbname=moneybot port=%s sslmode=disable TimeZone=Asia/Taipei", db_host, db_pwd, db_port)
	// time.Sleep(10 * time.Second)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&repo.User{}, &repo.Account{}, &repo.Tag{}, &repo.Cate{})

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
	repo_cate := repo.NewCateRepo(db)

	liff_conf := viper.GetStringMap("liff")
	handler.NewLineHandler(r, bot, *repo_user, *repo_tag, *repo_acc, *repo_cate)
	handler.NewTagHandler(r, *repo_user, *repo_tag, liff_conf)
	handler.NewUserHandler(r, *repo_user, *repo_acc)
	handler.NewAccHandler(r, *repo_user, *repo_acc, *repo_cate, liff_conf)

	PORT := viper.GetString("port")
	log.Printf("Server Start at Port: %s", PORT)
	r.Run(fmt.Sprintf(":%s", PORT))
}
