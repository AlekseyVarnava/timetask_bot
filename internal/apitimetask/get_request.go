package apitimetask

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"TimeTaskBot/internal/authorization"
	"TimeTaskBot/internal/utils"
)

// получает информацию о всех пользователях, которые используют телеграм бота. Отсюда брать UserID и TelegramID. telegramUsers возвращает данные структуры TelegramUsers
func GetTelegramUsers() (*utils.TelegramUsers, error) {
	req, err := http.NewRequest("GET", utils.GetTelegramUsersURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create TelegramUsers info request: %v", err)
	}
	authorization.TokenMutex.RLock()
	req.Header.Set("Authorization", "Bearer "+utils.AuthToken)
	authorization.TokenMutex.RUnlock()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in TelegramUsers info request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read TelegramUsers info response: %v", err)
	}
	var telegramUsers utils.TelegramUsers
	err = json.Unmarshal(body, &telegramUsers)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal TelegramUsers info response: %v", err)
	}
	return &telegramUsers, nil
}

// Функция для получения информации о пользователе
func getUserInfo(userId string) error {
	getUserInfoURL := fmt.Sprintf("%sUserId=%s", utils.BaseUserInfoURL, userId)
	req, err := http.NewRequest("GET", getUserInfoURL, nil)
	if err != nil {
		return fmt.Errorf("could not create user info request: %v", err)
	}
	authorization.TokenMutex.RLock()
	req.Header.Set("Authorization", "Bearer "+utils.AuthToken)
	authorization.TokenMutex.RUnlock()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error in user info request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read user info response: %v", err)
	}

	var userInfo utils.UserInfoResponse
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return fmt.Errorf("could not unmarshal user info response: %v", err)
	}

	fmt.Printf("User Info: %+v\n", userInfo)
	return nil
}

// тут получаем информацию о задачах пользователя в структуру TaskInfoResponse. Выглядит так:[{ID:5421 UserID:eb4d4f39-5631-4817-b4e7-26f90358da2a Title:Название задачи тут! Description:Описание задачи тут! IsCompleted:false IsDelete:false Priority:3 Hours:0 Minutes:0 Date:17.11.2024 Time:17:00 RunTime:0 IsOpened:false TimerStart:0 TimerTime:0 TimerIsActive:false TransferCount:0 IsRepit:false Tags:[] NotificationIntervals:1}]
func GetTaskInfo(userId string) (utils.TaskInfoResponse, error) {
	dateNow := time.Now().Format("2006-01-02")
	getTaskURL := fmt.Sprintf("%sUserId=%s&Date=%s", utils.BaseTaskURL, userId, dateNow)
	req, err := http.NewRequest("GET", getTaskURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create Task info request: %v", err)
	}
	authorization.TokenMutex.RLock()
	req.Header.Set("Authorization", "Bearer "+utils.AuthToken)
	authorization.TokenMutex.RUnlock()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in Task info request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read Task info response: %v", err)
	}
	var taskInfo utils.TaskInfoResponse
	err = json.Unmarshal(body, &taskInfo)
	if err != nil {
		return nil, fmt.Errorf("GetTaskInfo ошибка unmarshal: %v, responseBody:%v", err, string(body))
	}
	// if len(taskInfo) > 0 {
	// 	fmt.Println("Task Info:", taskInfo)
	// }
	return taskInfo, nil
}
