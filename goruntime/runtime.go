package goruntime

import (
	"runtime"
	"syscall"
)

var (
	ru  syscall.Rusage
	rtm runtime.MemStats
)

type RuntimeStat struct {
	MemoryMaxRss int64  `json:"memory_max_rss"`
	NumGoroutine int    `json:"num_goroutine"`
	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	Mallocs      uint64 `json:"mallocs"`
	Frees        uint64 `json:"free"`
	LiveObjects  uint64 `json:"live_objects"`
	PauseTotalNs uint64 `json:"pause_total_ns"`
	NumGC        uint32 `json:"num_gc"`
}

func GetStat() RuntimeStat {
	syscall.Getrusage(syscall.RUSAGE_SELF, &ru)
	maxRss := ru.Maxrss

	runtime.ReadMemStats(&rtm)

	// https://golang.org/pkg/runtime/
	var stat RuntimeStat
	stat.NumGoroutine = runtime.NumGoroutine()
	stat.MemoryMaxRss = maxRss
	stat.Alloc = rtm.Alloc
	stat.TotalAlloc = rtm.TotalAlloc
	stat.Sys = rtm.Sys
	stat.Mallocs = rtm.Mallocs
	stat.Frees = rtm.Frees
	stat.LiveObjects = stat.Mallocs - stat.Frees
	stat.PauseTotalNs = rtm.PauseTotalNs
	stat.NumGC = rtm.NumGC
	return stat
}
