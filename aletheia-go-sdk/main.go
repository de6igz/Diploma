package main

import (
	"aletheia-go-sdk/sdk"
	"aletheia-go-sdk/sdk/config"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {

	cfg := config.AletheiaConfig{
		ProjectId:        "13",
		CollectorAddress: "http://localhost:7070/aletheia-collector-service",
		ServiceName:      "Example Service",
		Environment:      "prod",
		Version:          "1.0.0",
		AuthToken:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU0MDY4MzAsImlhdCI6MTc0NDU0MjgzMCwic3ViIjoiMTEifQ.7aP8ctQjgPYEEQNY9aG8hXdFJtWp-zgubKgjWmwJtnA",
		UseTLS:           false,
		CaCertPath:       "",
		RPCTimeout:       5 * time.Second,
		Logger:           config.NewDefaultLogger(),
	}

	aletheia, err := sdk.NewClient(cfg)
	if err != nil {
		log.Fatalf("SDK initialization error: %v", err)
	}
	defer aletheia.Close()

	// Установка тегов
	aletheia.SetTags(map[string]string{
		"deployment": "production",
		"region":     "eu-west-1",
		"demo":       "1",
	})
	aletheia.SetCustomField("negoden", "1")
	aletheia.SetCustomField("customField1", "1")
	aletheia.SetCustomField("customField2", "2")
	//customFields := map[string]interface{}{
	//	"fields.negoden": "1",
	//}

	// Запуск мониторинга ресурсов
	aletheia.MonitorResources(2*time.Second, AllMetricsCollector)

	//Имитация работы приложения
	for i := 0; i < 10000; i++ {
		go func() {
			time.Sleep(20 * time.Second)
		}()
	}

	// Ожидание и отправка тестовых событий
	time.Sleep(55 * time.Second)

	// Пример контекста
	ctxMap := map[string]interface{}{
		"transaction_id": "12345",
		"component":      "payment-processor",
	}

	//time.Sleep(10 * time.Second)
	// Отправка информационного сообщения
	//aletheia.CaptureMessage("REST API initialized", models.LogLevelInfo, ctxMap)

	err = fmt.Errorf("ERROR PANIC")
	wg := sync.WaitGroup{}
	go func(wg *sync.WaitGroup) {
		//time.Sleep(5 * time.Second)
	}(&wg)
	aletheia.CaptureException(err, "PKMU", "custom", ctxMap)
	wg.Wait()
	// Пример отправки ошибки
	// err = fmt.Errorf("critical database failure")
	// aletheia.CaptureException(err, ctxMap)
}
