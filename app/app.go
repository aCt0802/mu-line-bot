package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

//Secret TOMLファイル
type Secret struct {
	Line LineSecret
}

//LineSecret Secretファイルのline関連部分
type LineSecret struct {
	BotChannelSecret      string
	BotChannelAccessToken string
}

func HandleMassage(c *gin.Context) {

	// secret情報をTOMLファイルか読み出し
	var secret Secret
	_, err := toml.DecodeFile("./secret.toml", &secret)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(secret.Line.BotChannelSecret)
	fmt.Println(secret.Line.BotChannelAccessToken)

	// LINE Botクライアント生成する
	// BOT にはチャネルシークレットとチャネルトークンをTomlファイルから読み込んで渡す
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

	// "可愛い" 単語を含む場合、返信される
	var replyText string
	replyText = "可愛い"

	// チャットの回答
	var response string
	response = "ありがとう！！"

	// "おはよう" 単語を含む場合、返信される
	var replySticker string
	replySticker = "おはよう"

	// スタンプで回答が来る
	responseSticker := linebot.NewStickerMessage("11537", "52002757")

	// "猫" 単語を含む場合、返信される
	var replyImage string
	replyImage = "猫"

	// 猫の画像が表示される
	responseImage := linebot.NewImageMessage("https://i.gyazo.com/2db8f85c496dd8f21a91eccc62ceee05.jpg", "https://i.gyazo.com/2db8f85c496dd8f21a91eccc62ceee05.jpg")

	// "ディズニー" 単語を含む場合、返信される
	var replyLocation string
	replyLocation = "ディズニー"

	// ディズニーが地図表示される
	responseLocation := linebot.NewLocationMessage("東京ディズニーランド", "千葉県浦安市舞浜", 35.632896, 139.880394)

	for _, event := range events {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// メッセージがテキスト形式の場合
			case *linebot.TextMessage:
				replyMessage := message.Text
				// テキストで返信されるケース
				if strings.Contains(replyMessage, replyText) {
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(response)).Do()
					// スタンプで返信されるケース
				} else if strings.Contains(replyMessage, replySticker) {
					bot.ReplyMessage(event.ReplyToken, responseSticker).Do()
					// 画像で返信されるケース
				} else if strings.Contains(replyMessage, replyImage) {
					bot.ReplyMessage(event.ReplyToken, responseImage).Do()
					// 地図表示されるケース
				} else if strings.Contains(replyMessage, replyLocation) {
					bot.ReplyMessage(event.ReplyToken, responseLocation).Do()
				}
				// 上記意外は、おうむ返しで返信
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}
