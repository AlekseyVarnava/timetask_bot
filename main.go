package main

import (
	"TimeTaskBot/internal/authorization"
	"TimeTaskBot/internal/botlogic"
	"TimeTaskBot/internal/telegram"
)

func main() {

	// Авторизация и обновление токена
	err := telegram.InitBot()
	if err != nil {
		panic(err)
	}
	authorization.Auth()

	// Запуск бота
	botlogic.Launch()
}
