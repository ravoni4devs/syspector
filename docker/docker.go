package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ravoni4devs/syspector/internal/common"
)

const (
	cgroupBasePath = "/sys/fs/cgroup"
)

var (
	cgroupVersion   = 2
	cgroupDetectErr error
)

type VirtualMemoryStat struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

type DockerStat struct {
	Memory     VirtualMemoryStat `json:"memory"`
	CpuPercent float64           `json:"cpuPercent"`
}

func (m VirtualMemoryStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func init() {
	cgroupVersion, cgroupDetectErr = detectCgroupVersion()
}

func GetStat(duration time.Duration) (DockerStat, error) {
	var stat DockerStat

	m, err := VirtualMemory()
	if err != nil {
		return stat, err
	}

	c, err := CpuPercent(duration)
	if err != nil {
		return stat, err
	}

	stat.Memory = m
	stat.CpuPercent = c
	return stat, nil
}

func VirtualMemory() (VirtualMemoryStat, error) {
	var stat VirtualMemoryStat
	if cgroupDetectErr != nil {
		return stat, cgroupDetectErr
	}

	if cgroupVersion == 1 {
		usedBytes, err := common.ReadFileNoStat(fmt.Sprintf("%s/memory/memory.usage_in_bytes", cgroupBasePath))
		if err != nil {
			return stat, err
		}
		limitBytes, err := common.ReadFileNoStat(fmt.Sprintf("%s/memory/memory.limit_in_bytes", cgroupBasePath))
		if err != nil {
			return stat, err
		}
		stat.Used = common.ParseUint64(string(usedBytes))
		stat.Total = common.ParseUint64(string(limitBytes))
	} else {
		usedBytes, err := common.ReadFileNoStat(fmt.Sprintf("%s/memory.current", cgroupBasePath))
		if err != nil {
			return stat, err
		}
		limitBytes, err := common.ReadFileNoStat(fmt.Sprintf("%s/memory.max", cgroupBasePath))
		if err != nil {
			return stat, err
		}
		stat.Used = common.ParseUint64(string(usedBytes))

		if strings.TrimSpace(string(limitBytes)) == "max" {
			meminfo, err := os.ReadFile("/proc/meminfo")
			if err != nil {
				return stat, err
			}
			stat.Total = parseMemInfoTotal(string(meminfo))
		} else {
			stat.Total = common.ParseUint64(string(limitBytes))
		}
	}

	if stat.Total > stat.Used {
		stat.Free = stat.Total - stat.Used
	} else {
		stat.Free = 0
	}

	stat.Available = stat.Free

	if stat.Total > 0 {
		stat.UsedPercent = common.ParseFloat(fmt.Sprintf("%.2f", (float64(stat.Used)/float64(stat.Total))*100))
	}

	return stat, nil
}

func CpuPercent(seconds time.Duration) (float64, error) {
	if cgroupDetectErr != nil {
		return 0, cgroupDetectErr
	}

	var usage1, usage2 uint64
	var err error

	if cgroupVersion == 1 {
		usage1, err = readCpuacctUsage()
		if err != nil {
			return 0, err
		}
	} else {
		usage1, err = readCgroupV2CpuUsage()
		if err != nil {
			return 0, err
		}
	}

	time.Sleep(seconds)

	if cgroupVersion == 1 {
		usage2, _ = readCpuacctUsage()
	} else {
		usage2, _ = readCgroupV2CpuUsage()
	}

	// uso total entre os dois samples
	deltaUsage := usage2 - usage1
	// delta time em nanos
	deltaTime := seconds.Seconds() * 1e9

	percent := (float64(deltaUsage) / deltaTime) * 100
	return common.ParseFloat(fmt.Sprintf("%.2f", percent)), nil
}

func readCpuacctUsage() (uint64, error) {
	data, err := common.ReadFileNoStat(fmt.Sprintf("%s/cpuacct/cpuacct.usage", cgroupBasePath))
	if err != nil {
		return 0, err
	}
	return common.ParseUint64(string(data)), nil
}

func readCgroupV2CpuUsage() (uint64, error) {
	data, err := common.ReadFileNoStat(fmt.Sprintf("%s/cpu.stat", cgroupBasePath))
	if err != nil {
		return 0, err
	}
	lines := splitLines(string(data))
	for _, line := range lines {
		if hasPrefix(line, "usage_usec") {
			fields := splitFields(line)
			if len(fields) == 2 {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				return val * 1000, nil // converte para nanossegundos
			}
		}
	}
	return 0, errors.New("usage_usec not found in cpu.stat")
}

func detectCgroupVersion() (int, error) {
	if runtime.GOOS != "linux" {
		return 0, errors.New("only runs on Linux")
	}

	if _, err := os.Stat("/sys/fs/cgroup/cgroup.controllers"); err == nil {
		return 2, nil
	}
	return 1, nil
}

func parseMemInfoTotal(meminfo string) uint64 {
	lines := splitLines(meminfo)
	for _, line := range lines {
		if hasPrefix(line, "MemTotal:") {
			fields := splitFields(line)
			if len(fields) >= 2 {
				val, _ := strconv.ParseUint(fields[1], 10, 64)
				return val * 1024 // kB -> bytes
			}
		}
	}
	return 0
}

func splitLines(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}

func splitFields(s string) []string {
	return strings.Fields(s)
}

func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}
