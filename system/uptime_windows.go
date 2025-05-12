//go:build windows

package system

import "syscall"

func Uptime() (float64, error) {
	lib := syscall.NewLazyDLL("kernel32.dll")
	proc := lib.NewProc("GetTickCount64")

	ret, _, err := proc.Call()
	if err != nil && err.Error() != "The operation completed successfully." {
		return 0, err
	}

	// GetTickCount64 retorna milissegundos desde a inicialização
	uptimeMs := uint64(ret)
	return float64(uptimeMs) / 1000.0, nil
}
