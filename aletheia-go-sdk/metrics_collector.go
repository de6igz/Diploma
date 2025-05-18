package main

import (
	"runtime"
)

// AllMetricsCollector собирает всевозможные метрики о состоянии среды выполнения Go.
func AllMetricsCollector() map[string]interface{} {
	metricsMap := make(map[string]interface{})

	// Сбор базовых метрик памяти через runtime.MemStats.
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	metricsMap["memory_alloc_bytes"] = memStats.Alloc
	metricsMap["memory_total_alloc_bytes"] = memStats.TotalAlloc
	metricsMap["memory_sys_bytes"] = memStats.Sys
	metricsMap["memory_mallocs"] = memStats.Mallocs
	metricsMap["memory_frees"] = memStats.Frees
	metricsMap["heap_alloc_bytes"] = memStats.HeapAlloc
	metricsMap["heap_sys_bytes"] = memStats.HeapSys
	metricsMap["heap_idle_bytes"] = memStats.HeapIdle
	metricsMap["heap_inuse_bytes"] = memStats.HeapInuse
	metricsMap["heap_objects"] = memStats.HeapObjects

	// Сбор базовых системных метрик.
	metricsMap["goroutine_count"] = runtime.NumGoroutine()
	metricsMap["num_cpu"] = runtime.NumCPU()
	metricsMap["max_procs"] = runtime.GOMAXPROCS(0)
	metricsMap["go_version"] = runtime.Version()
	metricsMap["os"] = runtime.GOOS
	metricsMap["arch"] = runtime.GOARCH

	// Сбор расширенных метрик через пакет runtime/metrics.
	// Все доступные метрики из runtime/metrics добавляются с префиксом "runtime_metrics."
	//for _, m := range metrics.All() {
	//	metricsMap["runtime_metrics."+m.Name] = m.Kind
	//}

	return metricsMap
}
