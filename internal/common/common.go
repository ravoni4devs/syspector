package common

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

var (
	Timeout                = 3 * time.Second
	ErrNotImplementedError = errors.New("not implemented yet")
	ErrTimeout             = errors.New("command timed out")
)

type Invoker interface {
	Command(string, ...string) ([]byte, error)
	CommandWithContext(context.Context, string, ...string) ([]byte, error)
}

type Invoke struct{}

func (i Invoke) Command(name string, arg ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	return i.CommandWithContext(ctx, name, arg...)
}

func (i Invoke) CommandWithContext(ctx context.Context, name string, arg ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, arg...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return buf.Bytes(), err
	}

	if err := cmd.Wait(); err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

type FakeInvoke struct {
	Suffix string
	Error  error
}

func (i FakeInvoke) Command(name string, arg ...string) ([]byte, error) {
	if i.Error != nil {
		return []byte{}, i.Error
	}

	arch := runtime.GOOS

	commandName := filepath.Base(name)

	fname := strings.Join(append([]string{commandName}, arg...), "")
	fname = url.QueryEscape(fname)
	fpath := path.Join("testdata", arch, fname)
	if i.Suffix != "" {
		fpath += "_" + i.Suffix
	}
	if PathExists(fpath) {
		return os.ReadFile(fpath)
	}
	return []byte{}, fmt.Errorf("could not find testdata: %s", fpath)
}

func (i FakeInvoke) CommandWithContext(_ context.Context, name string, arg ...string) ([]byte, error) {
	return i.Command(name, arg...)
}

// ReadFileNoStat uses ioutil.ReadAll to read contents of entire file.
// This is similar to ioutil.ReadFile but without the call to os.Stat, because
// many files in /proc and /sys report incorrect file sizes (either 0 or 4096).
// Reads a max file size of 512kB.  For files larger than this, a scanner
// should be used.
func ReadFileNoStat(filename string) ([]byte, error) {
	const maxBufferSize = 1024 * 512

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := io.LimitReader(f, maxBufferSize)
	return io.ReadAll(reader)
}

func ReadFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ReadLines(filename string) ([]string, error) {
	return ReadLinesOffsetN(filename, 0, -1)
}

// ReadLine reads a file and returns the first occurrence of a line that is prefixed with prefix.
func ReadLine(filename string, prefix string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		if strings.HasPrefix(line, prefix) {
			return line, nil
		}
	}

	return "", nil
}

// ReadLinesOffsetN reads contents from file and splits them by new line.
// The offset tells at which line number to start.
// The count determines the number of lines to read (starting from offset):
// n >= 0: at most n lines
// n < 0: whole file
func ReadLinesOffsetN(filename string, offset uint, n int) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for i := uint(0); i < uint(n)+offset || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF && len(line) > 0 {
				ret = append(ret, strings.Trim(line, "\n"))
			}
			break
		}
		if i < offset {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}

func IntToString(orig []int8) string {
	ret := make([]byte, len(orig))
	size := -1
	for i, o := range orig {
		if o == 0 {
			size = i
			break
		}
		ret[i] = byte(o)
	}
	if size == -1 {
		size = len(orig)
	}

	return string(ret[0:size])
}

func UintToString(orig []uint8) string {
	ret := make([]byte, len(orig))
	size := -1
	for i, o := range orig {
		if o == 0 {
			size = i
			break
		}
		ret[i] = byte(o)
	}
	if size == -1 {
		size = len(orig)
	}

	return string(ret[0:size])
}

func ByteToString(orig []byte) string {
	n := -1
	l := -1
	for i, b := range orig {
		// skip left side null
		if l == -1 && b == 0 {
			continue
		}
		if l == -1 {
			l = i
		}

		if b == 0 {
			break
		}
		n = i + 1
	}
	if n == -1 {
		return string(orig)
	}
	return string(orig[l:n])
}

func ParseFloat(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

func ParseUint64(s string) uint64 {
	s = strings.TrimSuffix(s, "\n")
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func ReadInts(filename string) ([]int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []int64{}, err
	}
	defer f.Close()

	var ret []int64

	r := bufio.NewReader(f)

	line, err := r.ReadString('\n')
	if err != nil {
		return []int64{}, err
	}

	i, err := strconv.ParseInt(strings.Trim(line, "\n"), 10, 32)
	if err != nil {
		return []int64{}, err
	}
	ret = append(ret, i)

	return ret, nil
}

func HexToUint32(hex string) uint32 {
	vv, _ := strconv.ParseUint(hex, 16, 32)
	return uint32(vv)
}

func StringsHas(target []string, src string) bool {
	for _, t := range target {
		if strings.TrimSpace(t) == src {
			return true
		}
	}
	return false
}

func StringsContains(target []string, src string) bool {
	for _, t := range target {
		if strings.Contains(t, src) {
			return true
		}
	}
	return false
}

func IntContains(target []int, src int) bool {
	return slices.Contains(target, src)
}

func PathExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

func PathExistsWithContents(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.Size() > 4 && !info.IsDir() // at least 4 bytes
}

func GetEnvWithContext(ctx context.Context, key string, dfault string, combineWith ...string) string {
	var value string
	if value == "" {
		value = os.Getenv(key)
	}
	if value == "" {
		value = dfault
	}

	return combine(value, combineWith)
}

func GetEnv(key string, dfault string, combineWith ...string) string {
	value := os.Getenv(key)
	if value == "" {
		value = dfault
	}

	return combine(value, combineWith)
}

func IsContainerized() (bool, error) {
	const procOneCgroup = "/proc/1/cgroup"
	data, err := os.ReadFile(procOneCgroup)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed to read process cgroups: %w", err)
	}
	return isContainerizedCgroup(data)
}

func isContainerizedCgroup(data []byte) (bool, error) {
	s := bufio.NewScanner(bytes.NewReader(data))
	for n := 0; s.Scan(); n++ {
		line := s.Bytes()
		// being inside a container: https://stackoverflow.com/a/20012536/235203
		if bytes.Contains(line, []byte("docker")) || bytes.Contains(line, []byte(".slice")) || bytes.Contains(line, []byte("lxc")) || bytes.Contains(line, []byte("kubepods")) {
			return true, nil
		}
	}
	return false, s.Err()
}

func combine(value string, combineWith []string) string {
	switch len(combineWith) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith)+1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
}

func HostProc(combineWith ...string) string {
	return GetEnv("HOST_PROC", "/proc", combineWith...)
}

func HostSys(combineWith ...string) string {
	return GetEnv("HOST_SYS", "/sys", combineWith...)
}

func HostEtc(combineWith ...string) string {
	return GetEnv("HOST_ETC", "/etc", combineWith...)
}

func HostVar(combineWith ...string) string {
	return GetEnv("HOST_VAR", "/var", combineWith...)
}

func HostRun(combineWith ...string) string {
	return GetEnv("HOST_RUN", "/run", combineWith...)
}

func HostDev(combineWith ...string) string {
	return GetEnv("HOST_DEV", "/dev", combineWith...)
}

func HostRoot(combineWith ...string) string {
	return GetEnv("HOST_ROOT", "/", combineWith...)
}

func HostProcWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_PROC", "/proc", combineWith...)
}

func HostProcMountInfoWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_PROC_MOUNTINFO", "", combineWith...)
}

func HostSysWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_SYS", "/sys", combineWith...)
}

func HostEtcWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_ETC", "/etc", combineWith...)
}

func HostVarWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_VAR", "/var", combineWith...)
}

func HostRunWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_RUN", "/run", combineWith...)
}

func HostDevWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_DEV", "/dev", combineWith...)
}

func HostRootWithContext(ctx context.Context, combineWith ...string) string {
	return GetEnvWithContext(ctx, "HOST_ROOT", "/", combineWith...)
}

func Round(val float64, n int) float64 {
	pow10 := math.Pow(10, float64(n))
	return math.Round(val*pow10) / pow10
}
