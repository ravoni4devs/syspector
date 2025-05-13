//go:build windows

package system

import (
	"syscall"
)

func Uptime() (float64, error) {
	// Windows Vista/Server 2008
	lib := syscall.NewLazyDLL("kernel32.dll")
	procGetTickCount64 := lib.NewProc("GetTickCount64")

	if procGetTickCount64.Find() == nil {
		ret, _, err := procGetTickCount64.Call()
		if err != syscall.Errno(0) {
			return 0, err
		}
		uptimeMs := uint64(ret)
		return float64(uptimeMs) / 1000.0, nil
	}

	// Fallback to GetTickCount
	procGetTickCount := lib.NewProc("GetTickCount")
	ret, _, err := procGetTickCount.Call()
	if err != syscall.Errno(0) {
		return 0, err
	}

	uptimeMs := uint32(ret)
	return float64(uptimeMs) / 1000.0, nil
}
