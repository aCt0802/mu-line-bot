package app

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

//Secret TOMLファイル
type Secret struct {
	Line   LineSecret
	OpenAI OpenAISecret
}

//LineSecret SecretファイルのLine関連部分
type LineSecret struct {
	BotChannelSecret      string `toml:"botChannelSecret"`
	BotChannelAccessToken string `toml:"botChannelAccessToken"`
}

//OpenAISecret SecretファイルのOpenAI関連部分
type OpenAISecret struct {
	ApiKey string `toml:"apiKey"`
}

func HandleMassage(c *gin.Context) {

	// secret情報をTOMLファイルか読み出し
	var secret Secret
	_, err := toml.DecodeFile("./secret.toml", &secret)
	if err != nil {
		fmt.Println(err)
	}

	// LINE Botクライアント生成する
	// BOT にはチャネルシークレットとチャネルトークンをTOMLファイルから読み込んで渡す
	bot, err := linebot.New(
		secret.Line.BotChannelSecret,
		secret.Line.BotChannelAccessToken,
	)

	// エラーに値があればログに出力し終了する
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			log.Print(err)
		}
		return
	}

	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				aiMessageResponse := getOpenAIResponse(secret.OpenAI.ApiKey, message.Text)
				log.Print(aiMessageResponse.Choices[0].Messages.Content)
				replyMessage := aiMessageResponse.Choices[0].Messages.Content
				// 人工知能からの返事をbotが送信
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Print(err)
				}

			}
		}
	}
}
