package botlogic

import (
	"TimeTaskBot/internal/apitimetask"
	"TimeTaskBot/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"TimeTaskBot/internal/telegram"
)

var (
	telegramUsers *utils.TelegramUsers
	sliceTask     sync.Map
)

// Запуск бота
func Launch() {
	go returnTelegramUsers()
	time.Sleep(time.Second * 15)

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	semaphore := make(chan struct{}, utils.MaxGoroutines_returnTaskInfo) // Семафор для ограничения горутин returnTaskInfo

	for range ticker.C {
		semaphore <- struct{}{} // Захватываем слот
		go func() {
			defer func() { <-semaphore }() // Освобождаем слот
			returnTaskInfo()
		}()
	}
}

// возвращаю всех юзеров зарегистрированных через ТГ
func returnTelegramUsers() {
	for {
		telegramUsersNew, err := apitimetask.GetTelegramUsers()
		if err != nil {
			fmt.Printf("Ошибка получения пользователей зарегистрированных в Телеграм: %v\n", err)
			time.Sleep(time.Minute * 20)
			continue
		}
		for i, user := range *telegramUsersNew {
			timeZone, err := parseStringToTimeDuration(user.TimeZoneOffset)
			if err != nil {
				fmt.Printf("Ошибка парсинга часового пояса: %v\n", err)
				continue
			}
			(*telegramUsersNew)[i].TimeZoneTimeDutation = timeZone
		}
		// if len(telegramUsers) != len(telegramUsersNew) { // подумать тут, он не будет обновлять если юзер сменил часовой пояс
		// 	fmt.Println("Найдены новые пользователи Telegram, список пользователей обновлён")
		telegramUsers = telegramUsersNew
		// }
		fmt.Println("Telegram users:", telegramUsers)
		time.Sleep(time.Minute * 5)
	}
}

func returnTaskInfo() {
	for _, telegramUser := range *telegramUsers {
		fmt.Println("sliceTask contents:")
		sliceTask.Range(func(key, value interface{}) bool {
			fmt.Printf("Key: %v, Value: %v\n", key, value)
			return true // продолжить обход
		})

		taskInfoNew, err := apitimetask.GetTaskInfo(telegramUser.UserID)
		if err != nil {
			fmt.Printf("Ошибка получения информации о задаче: %v\n", err)
			continue
		}
		for _, task := range taskInfoNew {

			_, exists := sliceTask.Load(task.ID)
			if exists {
				continue // Уже обработана
			}
			sliceTask.Store(task.ID, false)

			taskTime, boolSend, err := parseDateTime(task.Date, task.Time)
			if err != nil {
				fmt.Printf("Ошибка разбора даты и времени: %v\n", err)
				continue
			} else if !boolSend {
				continue // время не указано, отправлять уведомление не надо
			}

			var notificationTime time.Time
			intervalMap := map[string]time.Duration{
				utils.NotificationOneDay.TimeStr:  utils.NotificationOneDay.TimeDur * time.Minute,
				utils.Notification12Hours.TimeStr: utils.Notification12Hours.TimeDur * time.Minute,
				utils.Notification1Hour.TimeStr:   utils.Notification1Hour.TimeDur * time.Minute,
				utils.Notification30Mins.TimeStr:  utils.Notification30Mins.TimeDur * time.Minute,
			}
			interval, ok := intervalMap[task.NotificationIntervals]
			if !ok {
				interval = 60 * time.Minute // Интервал по умолчанию
			}
			notificationTime = taskTime.Add(interval).Add(-telegramUser.TimeZoneTimeDutation)
			go scheduleMessage(task.ID, telegramUser.ChatID, task, notificationTime)
		}
	}
}

func scheduleMessage(taskID int, chatID string, taskInfo utils.TaskInfoResponseOne, notificationTime time.Time) {
	fmt.Printf("Для задачи taskID:`%d`, chatID:`%s` уведомление будет отправлено: %v\n", taskID, chatID, notificationTime)
	delay := time.Until(notificationTime)
	if delay < 0 {
		fmt.Printf("Время уведомления для задачи %d уже прошло\n", taskID)
		return
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		value, exists  := sliceTask.Load(taskID)
		if exists && value == true {
			return // Уведомление уже отправлено
		}

		chatIDInt64, err := StringToInt64(chatID)
		if err != nil {
			fmt.Printf("Ошибка преобразования ChatID: %v\n", err)
			return
		}
		var message string
		switch taskInfo.Description {
		case "":
			message = fmt.Sprintf("⏰ Напоминание\nЗадача: %s\n• Дата: %s\n• Время: %s\n• Ссылка на задачу: https://demo.timetask.ru/%d",
				taskInfo.Title, taskInfo.Date, taskInfo.Time, taskID)
		default:
			message = fmt.Sprintf("⏰ Напоминание\nЗадача: %s\nОписание: %s\n• Дата: %s\n• Время: %s\n• Ссылка на задачу: https://demo.timetask.ru/%d",
				taskInfo.Title, taskInfo.Description, taskInfo.Date, taskInfo.Time, taskID)
		}
		sendMessage(taskID, chatIDInt64, message)
	}
}

func sendMessage(taskID int, chatID int64, message string) {
	for i := 0; i < utils.RetriesSendMessage; i++ {
		ok, err := telegram.TgAPI_SendMessage(chatID, message)
		if ok && err == nil {
			sliceTask.Store(taskID, true)
		}
		time.Sleep(utils.DelayRetriesSendMessage)
	}
}

func StringToInt64(input string) (int64, error) {
	result, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать строку %q в int64: %w", input, err)
	}
	return result, nil
}

func parseDateTime(dateStr, timeStr string) (time.Time, bool, error) {
	// Предполагается формат: "YYYY.MM.DD" и "HH:MM"
	const dateFormat = "02.01.2006"
	const timeFormat = "15:04"
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return time.Time{}, false, err
	}
	if timeStr == "" {
		return date, false, nil
	}
	parsedTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return time.Time{}, false, err
	}
	// Объединяем дату и время
	return time.Date(date.Year(), date.Month(), date.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local), true, nil
}

func parseStringToTimeDuration(timeStr string) (time.Duration, error) {
	// Разбиваем строку по разделителю ":"
	parts := strings.Split(timeStr, ":")
	if len(parts) == 0 {
		return 0, fmt.Errorf("ошибка формата часового пояса: %s", timeStr)
	}

	// Парсим первую часть (часы) в int
	hours, err := time.ParseDuration(parts[0] + "h")
	if err != nil {
		return 0, err
	}

	return hours, nil
}
