package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

var StartTime time.Time // 进程启动时间

// 系统信息
type SystemInfo struct {
	Hostname     string `json:"host_name"`
	Sysname      string `json:"sysname"`
	Release      string `json:"release"`
	Machine      string `json:"machine"`
	CPUModelName string `json:"cpu_model_name"`
	CPUFrequency string `json:"cpu_frequency"`
	CPUCores     int    `json:"cpu_cores"`
	Location     string `json:"Location"`
	GOVersion    string `json:"go_version"`
	ProcessID    int    `json:"process_id"`
	CmdLine      string `json:"command_line"`
}

// 内存和CPU负载信息
type LoadInfo struct {
	TimeString string  `json:"time_string"`
	MemTotal   uint64  `json:"memory_total"`
	MemUsed    uint64  `json:"memory_used"`
	MemFree    uint64  `json:"memory_free"`
	MemBuffers uint64  `json:"memory_buffers"`
	MemCached  uint64  `json:"memory_cached"`
	LoadAVG1   float64 `json:"load_avg_1"`
	LoadAVG5   float64 `json:"load_avg_5"`
	LoadAVG15  float64 `json:"load_avg_15"`
}

// 进程的内存占用等
type ProcessInfo struct {
	TimeString      string `json:"time_string"`
	Uptime          int64  `json:"uptime"`
	VirtualMemory   int64  `json:"virtual_memory"`
	ResisdentMemory int64  `json:"resident_memory"`
	SharedMemory    int64  `json:"shared_memory"`
}

// golang运行时的内存状态
type RuntimeStatus struct {
	TimeString string `json:"time_string"`
	// General statistics
	Alloc   uint64 `json:"alloc_bytes"`
	Sys     uint64 `json:"sys_bytes"`
	Mallocs uint64 `json:"mallocs_bytes"`
	Frees   uint64 `json:"frees_bytes"`

	// Main allocation heap statistics
	HeapAlloc   uint64 `json:"heap_alloc_bytes"`
	HeapSys     uint64 `json:"heap_sys_bytes"`
	HeapIdle    uint64 `json:"heap_idle_bytes"`
	HeapInuse   uint64 `json:"heap_inuse_bytes"`
	HeapObjects uint64 `json:"heap_objests"`

	// stack statistics
	StackInuse  uint64 `json:"stack_inuse_bytes"`
	StackSys    uint64 `json:"stack_sys_bytes"`
	MSpanInuse  uint64 `json:"mspan_inuse_bytes"`
	MSpanSys    uint64 `json:"mspan_sys_bytes"`
	MCacheInuse uint64 `json:"mcache_inuse_bytes"`
	MCacheSys   uint64 `json:"mcache_sys_bytes"`

	// GC status
	GCPause          float64 `json:"gc_pause"`
	GCPausePerSecond float64 `json:"gc_pause_per_second"`
	GCPerSecond      float64 `json:"gc_per_second"`
	GCTotalPause     float64 `json:"gc_total_pause"`

	//Num of goroutines
	Goroutines uint64 `json:"goroutines"`
}

func init() {
	StartTime = time.Now()
}

// 判断文件或目录是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// 读取文件的全部内容
func ReadFileFullContent(filename string) ([]byte, error) {
	if !Exist(filename) {
		return nil, fmt.Errorf("File %s is not exist", filename)
	}

	return ioutil.ReadFile(filename)
}

// 读取文件全部内容，返回所有的行
func ReadFileFullContentByLine(filename string) ([][]byte, error) {
	if !Exist(filename) {
		return nil, fmt.Errorf("File %s is not exist", filename)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return bytes.Split(data, []byte("\n")), nil
}

// 获取当前时间的字符串格式
func GetTimeString() string {
	return time.Now().Format(TimeLayout)
}

// 获取当前时间所属的Location
func GetTimeLocation() string {
	zonename, _ := time.Now().In(time.Local).Zone()
	return zonename
}

func ConvertToSlice(data interface{}) []byte {
	if value, ok := data.([65]int8); ok {
		b := make([]byte, len(value))

		i := 0
		for ; i < len(value); i++ {
			if value[i] == 0 {
				break
			}
			b[i] = byte(value[i])
		}
		return b[:i]
	}
	return nil
}
