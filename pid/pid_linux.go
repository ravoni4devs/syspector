//go:build linux

package pid

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

func getStat(pidNumber int, interval time.Duration) (PidStat, error) {
	const clkTck = 100
	filename := "/proc/" + strconv.Itoa(pidNumber) + "/stat"

	data, err := os.ReadFile(filename)
	if err != nil {
		return PidStat{}, err
	}

	pidStatSample1, err := parsePidStatLinux(string(data))
	if err != nil {
		return PidStat{}, err
	}

	time.Sleep(interval)

	data, err = os.ReadFile(filename)
	if err != nil {
		return PidStat{}, err
	}

	pidStatSample2, err := parsePidStatLinux(string(data))
	if err != nil {
		return PidStat{}, err
	}

	totalTime1 := pidStatSample1.CpuTotalTimeSpent
	totalTime2 := pidStatSample2.CpuTotalTimeSpent

	deltaTime := float64(totalTime2 - totalTime1)
	cpuUsage := (deltaTime / (interval.Seconds() * float64(clkTck))) * 100

	pidStatSample2.CpuPercent = cpuUsage

	return pidStatSample2, nil
}

func parsePidStatLinux(data string) (PidStat, error) {
	fields := strings.Fields(data)
	if len(fields) < 24 {
		return PidStat{}, errors.New("invalid stat format")
	}

	utime, _ := strconv.ParseUint(fields[13], 10, 64)
	stime, _ := strconv.ParseUint(fields[14], 10, 64)
	cutime, _ := strconv.ParseUint(fields[15], 10, 64)
	cstime, _ := strconv.ParseUint(fields[16], 10, 64)
	numThreads, _ := strconv.Atoi(fields[19])

	totalCpuTime := utime + stime + cutime + cstime

	return PidStat{
		UTime:             uint64(utime),
		STime:             uint64(stime),
		CUTime:            uint64(cutime),
		CSTime:            uint64(cstime),
		NumThreads:        numThreads,
		CpuTotalTimeSpent: totalCpuTime,
	}, nil
}
