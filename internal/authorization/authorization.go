package authorization

import (
	"TimeTaskBot/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var TokenMutex sync.RWMutex

func Auth() {
	utils.Email = os.Getenv("EMAIL")
	if utils.Email == "" {
		log.Fatal("Email не найден в .env файле")
	}
	utils.Password = os.Getenv("PASSWORD")
	if utils.Email == "" {
		log.Fatal("Password не найден в .env файле")
	}

	err := login()
	if err != nil {
		fmt.Printf("Ошибка авторизации, err: %v\n", err)
		return
	}
	// Обновление токена каждые 290 секунд
	go func() {
		var countFail uint8
		time.Sleep(290 * time.Second)
		for countFail < 10 {
			// fmt.Println("Обновление токена, запуск")
			err := refreshAuthToken()
			if err != nil {
				fmt.Printf("Ошибка обновления токена, err: %v\n", err)
				countFail++
				time.Sleep(10 * time.Second)
				continue
			}
			// fmt.Println("Токен успешно обновлён")
			time.Sleep(290 * time.Second)
		}
		Auth() // заново авторизуемся, если countFail увеличился до 10
	}()
}

// Функция для авторизации и получения токена
func login() error {
	requestBody, err := json.Marshal(utils.LoginRequest{
		Email:      utils.Email,
		Password:   utils.Password,
		RememberMe: true,
	})
	if err != nil {
		return fmt.Errorf("could not marshal login request: %v", err)
	}

	resp, err := http.Post(utils.PostLoginURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error in login request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read login response: %v", err)
	}

	var loginResponse utils.LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return fmt.Errorf("could not unmarshal login response: %v", err)
	}

	TokenMutex.Lock()
	utils.AuthToken = loginResponse.Token
	utils.RefreshToken = loginResponse.RefreshToken
	TokenMutex.Unlock()

	fmt.Println("Logged in successfully.")
	return nil
}

// Функция для обновления токена с помощью utils.refreshToken
func refreshAuthToken() error {
	requestBody, err := json.Marshal(utils.RefreshRequest{
		RefreshToken: utils.RefreshToken,
	})
	if err != nil {
		return fmt.Errorf("could not marshal refresh token request: %v", err)
	}

	req, err := http.NewRequest("POST", utils.GetRefreshURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("could not create refresh request: %v", err)
	}
	TokenMutex.RLock()
	req.Header.Set("Authorization", "Bearer "+utils.AuthToken)
	TokenMutex.RUnlock()
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error in refresh token request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read refresh token response: %v", err)
	}

	var refreshResponse utils.RefreshResponse
	err = json.Unmarshal(body, &refreshResponse)
	if err != nil {
		return fmt.Errorf("невозможно размаршалить RefreshToken из Response: %v", err)
	}

	TokenMutex.Lock()
	utils.AuthToken = refreshResponse.Token
	utils.RefreshToken = refreshResponse.RefreshToken
	TokenMutex.Unlock()

	// fmt.Println("Token refreshed successfully.")
	return nil
}
