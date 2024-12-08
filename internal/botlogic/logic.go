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
	sliceTask     map[uint32]bool
	taskMutex     sync.Mutex
)

// Запуск бота
func Launch() {
	sliceTask = make(map[uint32]bool)
	go returnTelegramUsers()
	time.Sleep(time.Second * 15)
	for {
		go returnTaskInfo()
		time.Sleep(time.Minute * 2)
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
			timeZone, err := timeStringToInt(user.TimeZoneOffset)
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
	for i := range *telegramUsers {
		taskInfoNew, err := apitimetask.GetTaskInfo((*telegramUsers)[i].UserID)
		if err != nil {
			fmt.Printf("Ошибка получения информации о задаче: %v\n", err)
			continue
		}
		for j := range taskInfoNew {
			taskID := uint32(taskInfoNew[j].ID)

			taskMutex.Lock()
			_, exists := sliceTask[taskID]
			if exists {
				taskMutex.Unlock()
				continue // Уже обработана
			}
			sliceTask[taskID] = false
			taskMutex.Unlock()

			taskTime, err := parseDateTime(taskInfoNew[j].Date, taskInfoNew[j].Time)
			if err != nil {
				fmt.Printf("Ошибка разбора даты и времени: %v\n", err)
				continue
			}

			var notificationTime time.Time
			switch taskInfoNew[j].NotificationIntervals {
			case "1440": // 1 день
				notificationTime = taskTime.Add(-1440 * time.Minute).Add(-(*telegramUsers)[i].TimeZoneTimeDutation)
			case "720": // 12 часов
				notificationTime = taskTime.Add(-720 * time.Minute).Add(-(*telegramUsers)[i].TimeZoneTimeDutation)
			case "60": // 1 час
				notificationTime = taskTime.Add(-60 * time.Minute).Add(-(*telegramUsers)[i].TimeZoneTimeDutation)
			case "30": // 30 минут
				notificationTime = taskTime.Add(-30 * time.Minute).Add(-(*telegramUsers)[i].TimeZoneTimeDutation)
			default: // По умолчанию за 1 час
				notificationTime = taskTime.Add(-60 * time.Minute).Add(-(*telegramUsers)[i].TimeZoneTimeDutation)
			}
			go scheduleMessage(taskID, (*telegramUsers)[i].ChatID, taskInfoNew[j], notificationTime)
		}
	}
}

func scheduleMessage(taskID uint32, chatID string, taskInfo utils.TaskInfoResponseOne, notificationTime time.Time) {
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
		taskMutex.Lock()
		if sliceTask[taskID] {
			taskMutex.Unlock()
			return // Уведомление уже отправлено
		}
		taskMutex.Unlock()

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

func sendMessage(taskID uint32, chatID int64, message string) {
	ok, err := telegram.TgAPI_SendMessage(chatID, message)
	if ok && err == nil {
		taskMutex.Lock()
		sliceTask[taskID] = true
		taskMutex.Unlock()
	}
}

func StringToInt64(input string) (int64, error) {
	result, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось преобразовать строку %q в int64: %w", input, err)
	}
	return result, nil
}

func parseDateTime(dateStr, timeStr string) (time.Time, error) {
	// Предполагается формат: "YYYY-MM-DD" и "HH:MM"
	const dateFormat = "02.01.2006"
	const timeFormat = "15:04"
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	parsedTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	// Объединяем дату и время
	return time.Date(date.Year(), date.Month(), date.Day(),
		parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local), nil
}

func timeStringToInt(timeStr string) (time.Duration, error) {
	// Разбиваем строку по разделителю ":"
	parts := strings.Split(timeStr, ":")
	if len(parts) == 0 {
		return -99, fmt.Errorf("ошибка формата часового пояса: %s", timeStr)
	}

	// Парсим первую часть (часы) в int
	hours, err := time.ParseDuration(parts[0] + "h")
	if err != nil {
		return -99, err
	}

	return hours, nil
}
