package main

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sync/semaphore"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// IOCounters defines the data structure for I/O counters
type IOCounters struct {
	ReadCount  uint64
	WriteCount uint64
}

// ProcInfo defines the data structure for process information
type ProcInfo struct {
	Name    string
	Cmdline string
}

// LogMessage defines the data structure for log message
type LogMessage struct {
	Timestamp  string
	PID        int32
	Name       string
	Cmdline    string
	CPUPercent float64
	MemPercent float32
	IOCounters *IOCounters
	NetIO      []net.IOCountersStat
}

// Create a cache for storing process information
var procInfoCache = make(map[int32]*ProcInfo)
var cacheMutex sync.RWMutex

// Create an object pool for storing process information
var procInfoPool = sync.Pool{
	New: func() interface{} {
		return &ProcInfo{}
	},
}

// Create an object pool for storing I/O counters
var ioCountersPool = sync.Pool{
	New: func() interface{} {
		return &IOCounters{}
	},
}

// LogMessage channel
var logChan = make(chan LogMessage, 1000) // Buffer size can be adjusted as needed

// Retrieve process I/O counters from the /proc file system
func getIOCounters(pid int32) (*IOCounters, error) {
	filePath := fmt.Sprintf("/proc/%d/io", pid)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	counters := ioCountersPool.Get().(*IOCounters)
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, err
		}

		switch fields[0] {
		case "syscr:":
			counters.ReadCount = value
		case "syscw:":
			counters.WriteCount = value
		}
	}

	return counters, nil
}

// Refresh the process information cache
func refreshProcInfoCache() {
	pids, _ := process.Pids()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			// The process may no longer exist
			continue
		}
		info := procInfoPool.Get().(*ProcInfo)
		name, _ := proc.Name()
		info.Name = name
		cmdline, _ := proc.Cmdline()
		info.Cmdline = cmdline
		procInfoCache[pid] = info
	}
}

// Log writer goroutine
func logWriter() {
	for msg := range logChan {
		fmt.Printf("Timestamp: %s\n", msg.Timestamp)
		fmt.Printf("PID: %d\n", msg.PID)
		fmt.Printf("Name: %s\n", msg.Name)
		fmt.Printf("Cmdline: %s\n", msg.Cmdline)
		fmt.Printf("CPU: %.6f%%\n", msg.CPUPercent)
		fmt.Printf("Mem: %.6f%%\n", msg.MemPercent)

		if msg.IOCounters != nil {
			fmt.Printf("IO: %+v\n", msg.IOCounters)
		} else {
			fmt.Println("IO: <nil>, Please confirm if you are running the code with root privileges.")
		}

		for _, netStat := range msg.NetIO {
			fmt.Printf("Interface: %v\n", netStat.Name)
			fmt.Printf("Bytes Sent: %v\n", netStat.BytesSent)
			fmt.Printf("Bytes Recv: %v\n", netStat.BytesRecv)
		}

		fmt.Println("----------Robotics Monitor------------")
	}
}

// Print process and network information
func printProcessAndNetInfo(sem *semaphore.Weighted, pid int32) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		// The process may no longer exist
		sem.Release(1) // Release semaphore
		return
	}

	cacheMutex.RLock()
	info, ok := procInfoCache[pid]
	cacheMutex.RUnlock()
	if !ok {
		// The process info may not be in the cache
		sem.Release(1) // Release semaphore
		return
	}

	cpuPercent, _ := proc.CPUPercent()
	memPercent, _ := proc.MemoryPercent()
	ioCounters, _ := getIOCounters(pid)
	netIO, _ := net.IOCounters(false)

	// Send log message to the log writer
	logChan <- LogMessage{
		Timestamp:  time.Now().Format(time.RFC3339),
		PID:        pid,
		Name:       info.Name,
		Cmdline:    info.Cmdline,
		CPUPercent: cpuPercent,
		MemPercent: memPercent,
		IOCounters: ioCounters,
		NetIO:      netIO,
	}

	if ioCounters != nil {
		ioCountersPool.Put(ioCounters)
	}

	sem.Release(1)
}

func main() {
	sem := semaphore.NewWeighted(100) // Create a semaphore to limit concurrency

	go func() {
		for {
			refreshProcInfoCache()
			time.Sleep(1 * time.Second) // You can modify this time to change the refresh frequency
		}
	}()

	// Start the log writer goroutine
	go logWriter()

	for {
		pids, _ := process.Pids()
		for _, pid := range pids {
			if err := sem.Acquire(context.Background(), 1); err == nil {
				go printProcessAndNetInfo(sem, pid)
			}
		}

		time.Sleep(1 * time.Second) // You can modify this time to change the printing frequency
	}
}
