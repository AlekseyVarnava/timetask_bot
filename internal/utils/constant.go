package utils

import "time"

const (
	GetTelegramUsersURL = "https://api-demo.timetask.ru/api/Telegram/Users"
	PostLoginURL        = "https://api-demo.timetask.ru/api/Account/Login"
	GetRefreshURL       = "https://api-demo.timetask.ru/api/Account/refresh"
	BaseUserInfoURL     = "https://api-demo.timetask.ru/api/Account/User?"
	BaseTaskURL         = "https://api-demo.timetask.ru/api/Task?"

	MaxGoroutines_returnTaskInfo int           = 3 // кол-во макс горутин для returnTaskInfo
	RetriesSendMessage           int         = 2 // кол-во попыток повторной отправки сообщения
	DelayRetriesSendMessage      time.Duration = 2 * time.Second
)

var (
	AuthToken    string
	RefreshToken string
	Email        string
	Password     string

	// интервалы уведомлений
	NotificationOneDay = NotificationIntervals{
		TimeStr: "1440",
		TimeDur: 1440,
	}
	Notification12Hours = NotificationIntervals{
		TimeStr: "720",
		TimeDur: 720,
	}
	Notification1Hour = NotificationIntervals{
		TimeStr: "60",
		TimeDur: 60,
	}
	Notification30Mins = NotificationIntervals{
		TimeStr: "30",
		TimeDur: 30,
	}
)
