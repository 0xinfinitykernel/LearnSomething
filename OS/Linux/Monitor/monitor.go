package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sync/semaphore"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// IOCounters 定义 I/O 计数器的数据结构
type IOCounters struct {
	ReadCount  uint64
	WriteCount uint64
}

// ProcInfo 定义进程信息的数据结构
type ProcInfo struct {
	Name    string
	Cmdline string
}

// LogMessage 定义日志消息的数据结构
type LogMessage struct {
	Timestamp  time.Time
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

// Process name filter
var processNameFilter string

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
	newCache := make(map[int32]*ProcInfo)
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
		newCache[pid] = info
	}
	cacheMutex.Lock()
	procInfoCache = newCache
	cacheMutex.Unlock()
}

// Log writer goroutine
func logWriter(logger *log.Logger, done chan struct{}) {
	for msg := range logChan {
		logger.Printf("Timestamp: %s\n", msg.Timestamp.Format(time.RFC3339))
		logger.Printf("PID: %d\n", msg.PID)
		logger.Printf("Name: %s\n", msg.Name)
		logger.Printf("Cmdline: %s\n", msg.Cmdline)
		logger.Printf("CPU: %.6f%%\n", msg.CPUPercent)
		logger.Printf("Mem: %.6f%%\n", msg.MemPercent)

		if msg.IOCounters != nil {
			logger.Printf("IO: %+v\n", msg.IOCounters)
		} else {
			logger.Println("IO: <nil>, Please confirm if you are running the code with root privileges.")
		}

		for _, netStat := range msg.NetIO {
			logger.Printf("Interface: %v\n", netStat.Name)
			logger.Printf("Bytes Sent: %v\n", netStat.BytesSent)
			logger.Printf("Bytes Recv: %v\n", netStat.BytesRecv)
		}

		logger.Println("----------Robotics Monitor------------")
	}
	done <- struct{}{}
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

	// Check if process name matches the filter
	processNames := strings.Split(processNameFilter, ",")
	match := false
	if len(processNames) == 1 && processNames[0] == "" {
		// If no filter is set, match all processes
		match = true
	} else {
		for _, name := range processNames {
			if strings.Contains(info.Name, name) {
				match = true
				break
			}
		}
	}
	if !match {
		sem.Release(1) // Release semaphore
		return
	}

	cpuPercent, _ := proc.CPUPercent()
	memPercent, _ := proc.MemoryPercent()
	ioCounters, _ := getIOCounters(pid)
	netIO, _ := net.IOCounters(false)

	// Send log message to the log writer
	logChan <- LogMessage{
		Timestamp:  time.Now(),
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
	// Get the log file path from the command line arguments
	var logFilePath string
	flag.StringVar(&logFilePath, "log", "", "Log file path")
	flag.StringVar(&processNameFilter, "proc", "", "Process name filter (comma separated)")
	flag.Parse()

	sem := semaphore.NewWeighted(1000) // Create a semaphore to limit concurrency

	logger := log.New(os.Stdout, "", log.LstdFlags) // Default to standard output

	if logFilePath != "" {
		logger = log.New(&lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    2 * 1024, // megabytes
			MaxAge:     2,        // days
			MaxBackups: 2,        // The maximum number of old log files to retain
		}, "", log.LstdFlags)
	}

	done := make(chan struct{})
	defer func() {
		close(logChan)
		<-done // Wait for log writing to complete
		if logFilePath != "" {
			if lumberjackLogger, ok := logger.Writer().(*lumberjack.Logger); ok {
				_ = lumberjackLogger.Close()
			}
		}
	}()

	go func() {
		// Capture the interrupt signal to exit gracefully
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		done <- struct{}{}
	}()

	go func() {
		for {
			refreshProcInfoCache()
			time.Sleep(1 * time.Second) // You can modify this time to change the refresh frequency
		}
	}()

	// Start the log writer goroutine
	go logWriter(logger, done)

	for {
		select {
		case <-done:
			return
		default:
			pids, _ := process.Pids()
			for _, pid := range pids {
				if err := sem.Acquire(context.Background(), 1); err == nil {
					go printProcessAndNetInfo(sem, pid)
				}
			}

			time.Sleep(1 * time.Second) // You can modify this time to change the printing frequency
		}
	}
}
