package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-gpt3"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("TELEGAM_BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	gptClient, err := gpt3.NewClient("OPENAI_API")
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		response, err := gptClient.Complete(update.Message.Text, gpt3.ModelEngineDavinci, gpt3.CompleteRequest{
			MaxTokens:   100,
			Stop:        "",
			Temperature: 0.5,
			TopP:        1,
		})
		if err != nil {
			log.Println(err)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response.Choices[0].Text)
		bot.Send(msg)
	}
}
