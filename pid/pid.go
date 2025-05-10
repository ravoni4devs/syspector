//go:build linux || windows || darwin

package pid

import (
	"time"
)

type PidStat struct {
	PID               int     `json:"pid,omitempty"`
	CpuPercent        float64 `json:"cpu_percent"`
	MemoryUsage       uint64  `json:"memory_usage"`
	State             string  `json:"state,omitempty"`
	StateName         string  `json:"state_name,omitempty"`
	UTime             uint64  `json:"utime,omitempty"`               // clock ticks no user space
	STime             uint64  `json:"stime,omitempty"`               // clock ticks no kernel space
	CUTime            uint64  `json:"cutime,omitempty"`              // utime filhos
	CSTime            uint64  `json:"cstime,omitempty"`              // stime filhos
	NumThreads        int     `json:"num_threads,omitempty"`
	VSize             uint64  `json:"vsize,omitempty"`               // bytes
	RSS               int64   `json:"rss,omitempty"`                 // páginas de memória
	CpuTotalTimeSpent uint64  `json:"cpu_total_time_spent,omitempty"` // soma utime+stime+cutime+cstime
}

func GetStat(pidNumber int, interval time.Duration) (PidStat, error) {
	return getStat(pidNumber, interval)
}
