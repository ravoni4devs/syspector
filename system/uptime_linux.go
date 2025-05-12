package system

import "syscall"

func Uptime() (float64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return 0, err
	}

	// Uptime in secs
	return float64(info.Uptime), nil
}
