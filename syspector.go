//go:build linux || darwin || windows

package syspector

import (
	"fmt"
	"os"
	"time"

	"github.com/ravoni4devs/syspector/cpu"
	"github.com/ravoni4devs/syspector/docker"
	"github.com/ravoni4devs/syspector/goruntime"
	"github.com/ravoni4devs/syspector/mem"
	"github.com/ravoni4devs/syspector/pid"
	"github.com/ravoni4devs/syspector/system"
)

type Collector interface {
	Stats() (Stats, error)
	GetStatsByPID(pidNumber int) (Stats, error)
}

type statCollector struct{}

func New() Collector {
	return &statCollector{}
}

type Stats struct {
	Memory     mem.VirtualMemoryStat `json:"memory"`
	CpuPercent float64               `json:"cpu"`
	Runtime    goruntime.RuntimeStat `json:"runtime"`
	PID        pid.PidStat           `json:"pid"`
	System     system.SystemStat     `json:"system"`
}

func (c *statCollector) Stats() (Stats, error) {
	var stats Stats
	stats.Runtime = goruntime.GetStat()
	systemStat, err := system.GetStat()
	if err != nil {
		return stats, fmt.Errorf("getting system stats %s", err)
	}
	stats.System = systemStat

	pidStat, err := pid.GetStat(os.Getpid(), time.Second*1)
	if err != nil {
		return stats, fmt.Errorf("getting PID stats %s", err)
	}
	stats.PID = pidStat

	dockerStat, err := docker.GetStat(time.Second * 1)
	if err == nil {
		stats.CpuPercent = dockerStat.CpuPercent
		stats.Memory.Total = dockerStat.Memory.Total
		stats.Memory.Available = dockerStat.Memory.Available
		stats.Memory.Free = dockerStat.Memory.Free
		stats.Memory.Used = dockerStat.Memory.Used
		stats.Memory.UsedPercent = dockerStat.Memory.UsedPercent
		return stats, nil
	}
	memoryStat, err := mem.GetStat()
	if err != nil {
		return stats, fmt.Errorf("getting memory stats %s", err)
	}
	stats.Memory.Total = memoryStat.Total
	stats.Memory.Available = memoryStat.Available
	stats.Memory.Free = memoryStat.Free
	stats.Memory.Used = memoryStat.Used
	stats.Memory.UsedPercent = memoryStat.UsedPercent
	cpuPercent, err := cpu.Percent(time.Second*1, false)
	if err != nil || len(cpuPercent) == 0 {
		return stats, fmt.Errorf("getting cpu stats %s", err)
	}
	stats.CpuPercent = cpuPercent[0]

	return stats, nil
}

func (c *statCollector) GetStatsByPID(pidNumber int) (Stats, error) {
	var metrics Stats
	stat, err := pid.GetStat(pidNumber, time.Second*1)
	if err != nil {
		return metrics, err
	}
	metrics.PID = stat
	return metrics, nil
}
