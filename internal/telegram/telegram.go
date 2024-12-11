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

func StartBot() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // Если получили сообщение
			chatID := update.Message.Chat.ID
			chatType := update.Message.Chat.Type
			log.Printf("Получено сообщение из чата %d (тип: %s): %s", chatID, chatType, update.Message.Text)

			if chatType == "group" || chatType == "supergroup" {
				fmt.Println(chatID)
				msg := "Бот TimeTask добавлен в общий чат! Теперь он может отправлять сюда сообщения."
				TgAPI_SendMessage(chatID, msg)
			}
		}
	}
}