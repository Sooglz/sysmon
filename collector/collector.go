package collector

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Metrics struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
	DiskPercent   float64 `json:"disk_percent"`
	Timestamp     int64   `json:"timestamp"`
}

var (
	cpuMu     sync.Mutex
	lastTotal uint64
	lastIdle  uint64
)

func Collect(diskPath string) (Metrics, error) {
	cpu, err := cpuPercent()
	if err != nil {
		return Metrics{}, err
	}
	mem, err := memPercent()
	if err != nil {
		return Metrics{}, err
	}
	disk, err := diskPercent(diskPath)
	if err != nil {
		return Metrics{}, err
	}
	return Metrics{
		CPUPercent:    cpu,
		MemoryPercent: mem,
		DiskPercent:   disk,
		Timestamp:     time.Now().Unix(),
	}, nil
}

func cpuPercent() (float64, error) {
	cpuMu.Lock()
	defer cpuMu.Unlock()

	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return 0, fmt.Errorf("empty /proc/stat")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, fmt.Errorf("unexpected /proc/stat format")
	}

	var total, idle uint64
	for i, v := range fields[1:] {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		total += n
		if i == 3 || i == 4 {
			idle += n
		}
	}

	if lastTotal == 0 {
		lastTotal = total
		lastIdle = idle
		return 0, nil
	}

	totalDiff := total - lastTotal
	idleDiff := idle - lastIdle
	lastTotal = total
	lastIdle = idle

	if totalDiff == 0 {
		return 0, nil
	}
	return float64(totalDiff-idleDiff) / float64(totalDiff) * 100, nil
}

func memPercent() (float64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var total, available uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			total = parseKB(line)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			available = parseKB(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, fmt.Errorf("MemTotal not found")
	}
	used := total - available
	return float64(used) / float64(total) * 100, nil
}

func parseKB(line string) uint64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	n, _ := strconv.ParseUint(fields[1], 10, 64)
	return n
}

func diskPercent(path string) (float64, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, err
	}
	total := stat.Blocks
	free := stat.Bfree
	if total == 0 {
		return 0, nil
	}
	used := total - free
	return float64(used) / float64(total) * 100, nil
}
