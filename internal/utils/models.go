package utils

import "time"

type NotificationIntervals struct {
	TimeStr string
	TimeDur time.Duration
}

type TelegramUsers []struct {
	UserID               string `json:"userId"`
	ChatID               string `json:"telegramId"`
	TimeZoneOffset       string `json:"timeZoneOffset"`
	TimeZoneTimeDutation time.Duration
}

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

type LoginResponse struct {
	ID           string `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type UserInfoResponse struct {
	UserID           string `json:"userId"`
	Nickname         string `json:"nickname"`
	Name             string `json:"name"`
	Lastname         string `json:"lastname"`
	Birthday         string `json:"birthday"`
	Email            string `json:"email"`
	Position         string `json:"position"`
	RegistrationDate string `json:"registrationDate"`
	LastActivity     string `json:"lastActivity"`
}

type TaskInfoResponse []struct {
	ID                    int    `json:"id"`
	UserID                string `json:"userId"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	IsCompleted           bool   `json:"isCompleted"`
	IsDelete              bool   `json:"isDelete"`
	Priority              int    `json:"priority"`
	Hours                 int    `json:"hours"`
	Minutes               int    `json:"minutes"`
	Date                  string `json:"date"`
	Time                  string `json:"time"`
	RunTime               int    `json:"runTime"`
	IsOpened              bool   `json:"isOpened"`
	TimerStart            int    `json:"timerStart"`
	TimerTime             int    `json:"timerTime"`
	TimerIsActive         bool   `json:"timerIsActive"`
	TransferCount         int    `json:"transferCount"`
	IsRepit               bool   `json:"isRepit"`
	Tags                  []any  `json:"tags"`
	NotificationIntervals string `json:"notificationIntervals"`
}

type TaskInfoResponseOne struct {
	ID                    int    `json:"id"`
	UserID                string `json:"userId"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	IsCompleted           bool   `json:"isCompleted"`
	IsDelete              bool   `json:"isDelete"`
	Priority              int    `json:"priority"`
	Hours                 int    `json:"hours"`
	Minutes               int    `json:"minutes"`
	Date                  string `json:"date"`
	Time                  string `json:"time"`
	RunTime               int    `json:"runTime"`
	IsOpened              bool   `json:"isOpened"`
	TimerStart            int    `json:"timerStart"`
	TimerTime             int    `json:"timerTime"`
	TimerIsActive         bool   `json:"timerIsActive"`
	TransferCount         int    `json:"transferCount"`
	IsRepit               bool   `json:"isRepit"`
	Tags                  []any  `json:"tags"`
	NotificationIntervals string `json:"notificationIntervals"`
}
