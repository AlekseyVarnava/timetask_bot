package telegram

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var bot *tgbotapi.BotAPI

func InitBot() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Ошибка загрузки .env файла")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return fmt.Errorf("Токен не найден в .env файле")
	}
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("Ошибка инициализации бота: %v", err)
	}
	bot.Debug = false
	log.Println("Бот успешно запущен!")
	return nil
}

func TgAPI_SendMessage(chatID int64, message string) (bool, error) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Ошибка отправки сообщения: %v\n", err)
		return false, err
	}
	return true, nil
}
