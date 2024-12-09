package main

import (
	"TimeTaskBot/internal/authorization"
	"TimeTaskBot/internal/botlogic"
	"TimeTaskBot/internal/telegram"
)

func main() {

	// Авторизация и обновление токена
	telegram.InitBot()
	authorization.Auth()

	// Запуск бота
	botlogic.Launch()
} 