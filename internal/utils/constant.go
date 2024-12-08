package utils

const (
	GetTelegramUsersURL = "https://api-demo.timetask.ru/api/Telegram/Users"
	PostLoginURL        = "https://api-demo.timetask.ru/api/Account/Login"
	GetRefreshURL       = "https://api-demo.timetask.ru/api/Account/refresh"
	BaseUserInfoURL     = "https://api-demo.timetask.ru/api/Account/User?"
	BaseTaskURL         = "https://api-demo.timetask.ru/api/Task?"
)

var (
	AuthToken    string
	RefreshToken string
	Email        string
	Password     string
)
