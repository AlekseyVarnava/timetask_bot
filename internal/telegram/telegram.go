package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var bot *tgbotapi.BotAPI

func InitBot() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Токен не найден в .env файле")
	}
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}
	bot.Debug = false
	fmt.Println("Бот успешно запущен!")
}

func TgAPI_SendMessage(chatID int64, message string) (bool, error) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v\n", err)
		return false, err
	}
	fmt.Printf("Сообщение отправлено: %s\n", message)
	return true, nil
}
