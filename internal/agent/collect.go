package agent

import (
	"math/rand"
	"runtime"
)

func (ag *MyAgent) collectMetrics() {
	ag.collectRuntimeMetrics()
	ag.collectSystemMetrics()
}

func (ag *MyAgent) collectRuntimeMetrics() {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	ag.storage.SetGauge("Alloc", float64(stats.Alloc))
	ag.storage.SetGauge("BuckHashSys", float64(stats.BuckHashSys))
	ag.storage.SetGauge("Frees", float64(stats.Frees))
	ag.storage.SetGauge("GCCPUFraction", stats.GCCPUFraction)
	ag.storage.SetGauge("GCSys", float64(stats.GCSys))
	ag.storage.SetGauge("HeapAlloc", float64(stats.HeapAlloc))
	ag.storage.SetGauge("HeapIdle", float64(stats.HeapIdle))
	ag.storage.SetGauge("HeapInuse", float64(stats.HeapInuse))
	ag.storage.SetGauge("HeapObjects", float64(stats.HeapObjects))
	ag.storage.SetGauge("HeapReleased", float64(stats.HeapReleased))
	ag.storage.SetGauge("HeapSys", float64(stats.HeapSys))
	ag.storage.SetGauge("LastGC", float64(stats.LastGC))
	ag.storage.SetGauge("Lookups", float64(stats.Lookups))
	ag.storage.SetGauge("MCacheInuse", float64(stats.MCacheInuse))
	ag.storage.SetGauge("MCacheSys", float64(stats.MCacheSys))
	ag.storage.SetGauge("MSpanInuse", float64(stats.MSpanInuse))
	ag.storage.SetGauge("MSpanSys", float64(stats.MSpanSys))
	ag.storage.SetGauge("Mallocs", float64(stats.Mallocs))
	ag.storage.SetGauge("NextGC", float64(stats.NextGC))
	ag.storage.SetGauge("NumForcedGC", float64(stats.NumForcedGC))
	ag.storage.SetGauge("NumGC", float64(stats.NumGC))
	ag.storage.SetGauge("OtherSys", float64(stats.OtherSys))
	ag.storage.SetGauge("PauseTotalNs", float64(stats.PauseTotalNs))
	ag.storage.SetGauge("StackInuse", float64(stats.StackInuse))
	ag.storage.SetGauge("StackSys", float64(stats.StackSys))
	ag.storage.SetGauge("Sys", float64(stats.Sys))
	ag.storage.SetGauge("OtherSys", float64(stats.TotalAlloc))
	ag.storage.SetGauge("RandomValue", float64(stats.TotalAlloc))
}
func (ag *MyAgent) collectSystemMetrics() {
	ag.storage.SetGauge("RandomValue", rand.Float64())
	ag.storage.SetCounter("PollCount", 1)
}
