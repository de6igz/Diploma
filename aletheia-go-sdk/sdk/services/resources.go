package services

import (
	"aletheia-go-sdk/sdk/config"
	"aletheia-go-sdk/sdk/models"
	"aletheia-go-sdk/sdk/transport"
	"aletheia-go-sdk/sdk/utils"
	"github.com/sony/gobreaker"
	"runtime"
	"time"
)

type ResourceService struct {
	transport      *transport.HTTPTransport
	config         *config.AletheiaConfig
	logger         config.Logger
	stopChan       chan struct{}
	monitoringTags map[string]string
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewResourceService(cfg config.AletheiaConfig, tp *transport.HTTPTransport, logger config.Logger) *ResourceService {
	cbSettings := gobreaker.Settings{
		Name:        "ResourceServiceCB",
		MaxRequests: 5,
		Interval:    0,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures > 3 && float64(counts.TotalFailures)/float64(counts.Requests) > 0.5
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			logger.Infof("Circuit Breaker '%s' изменил состояние с %s на %s", name, from.String(), to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(cbSettings)

	return &ResourceService{
		transport:      tp,
		config:         &cfg,
		logger:         logger,
		stopChan:       make(chan struct{}),
		circuitBreaker: cb,
	}
}

func (s *ResourceService) StartMonitoring(
	interval time.Duration,
	tags map[string]string,
	userID string,
	userInfo map[string]interface{},
	metricsCollector func() map[string]interface{},
) {
	// Если пользователь не передал свою функцию-сборщик метрик, используем дефолтную.
	if metricsCollector == nil {
		metricsCollector = func() map[string]interface{} {
			mem, goroutines := getRuntimeMetrics()
			return map[string]interface{}{
				"memory_alloc_bytes": mem,
				"goroutine_count":    goroutines,
			}
		}
	}

	s.monitoringTags = tags
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				// Вызываем функцию для сбора метрик.
				fields := metricsCollector()

				// Добавляем пользовательский контекст.
				if userID != "" {
					fields["user_id"] = userID
				}
				if userInfo != nil {
					for k, v := range userInfo {
						fields["user_"+k] = v
					}
				}

				req := models.ResourceRequest{
					ProjectId:    s.config.ProjectId,
					ServiceName:  s.config.ServiceName,
					Environment:  s.config.Environment,
					Version:      s.config.Version,
					GoVersion:    runtime.Version(),
					Os:           runtime.GOOS,
					Arch:         runtime.GOARCH,
					EventType:    string(models.EventResourceUsage),
					EventMessage: "Resource usage report",
					Tags:         utils.FormatTags(s.monitoringTags),
					Timestamp:    time.Now().Format(time.RFC3339),
					Fields:       fields,
					Language:     Golang,
				}

				// Отправка с использованием Circuit Breaker.
				_, err := s.circuitBreaker.Execute(func() (interface{}, error) {
					return nil, s.transport.SendEvent(req)
				})

				if err != nil {
					s.logger.Errorf("Failed to send resource metrics: %v", err)
				} else {
					s.logger.Infof("Resource metrics sent %v", req)
				}

			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *ResourceService) StopMonitoring() {
	close(s.stopChan)
}

func getRuntimeMetrics() (uint64, int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc, runtime.NumGoroutine()
}
