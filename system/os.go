package system

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func getOSVersion() string {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/version"); err == nil {
			return strings.TrimSpace(string(data))
		}
	case "darwin":
		out, err := exec.Command("sw_vers", "-productVersion").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	case "windows":
		out, err := exec.Command("cmd", "ver").Output()
		if err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return "unknown"
}

func getLinuxDistro() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	for line := range strings.SplitSeq(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[len("PRETTY_NAME="):], "\"")
		}
	}
	return "unknown"
}

func detectContainer() string {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "docker"
	}
	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return "unknown"
	}
	types := []string{"docker", "kubepods", "containerd", "lxc", "podman"}
	for line := range strings.SplitSeq(string(data), "\n") {
		for _, t := range types {
			if strings.Contains(line, t) {
				return t
			}
		}
	}
	out, err := exec.Command("systemd-detect-virt", "--container").Output()
	if err == nil {
		result := strings.TrimSpace(string(out))
		if result != "none" {
			return result
		}
	}
	return "none"
}

func isVirtualized() bool {
	out, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return false
	}
	text := string(out)
	if strings.Contains(text, "hypervisor") {
		return true
	}

	cmd := exec.Command("systemd-detect-virt")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err == nil {
		result := strings.TrimSpace(stdout.String())
		if result != "" && result != "none" && result != "container" {
			return true
		}
	}

	return false
}
