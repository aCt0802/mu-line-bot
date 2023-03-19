package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/aCt0802/mu-line-bot/app"
	"github.com/gin-gonic/gin"
)

//Config TOMLファイル
type Config struct {
	App AppConfig
}

//AppConfig Configファイルのアプリケーション関連
type AppConfig struct {
	Port string `toml:"port"`
}

func main() {
	// config情報をTOMLファイルか読み出し
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		fmt.Println(err)
	}

	router := gin.New()
	router.Use(gin.Logger())

	// LINE Messaging API ルーティング
	router.POST("/callback", app.HandleMassage)
	router.Run(":" + config.App.Port)
}
