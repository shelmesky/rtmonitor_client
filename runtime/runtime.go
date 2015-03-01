package rtruntime

import (
	"fmt"
	"runtime"
	"time"
    "os"
    "github.com/shelmesky/rtmonitor_client/lib"
)


func GetRuntimeStats(channel chan *lib.RuntimeStatus, sleep time.Duration) {
	memStats := &runtime.MemStats{}
	lastSampleTime := time.Now()
	var lastPauseNs uint64 = 0
	var lastNumGc uint32 = 0
    var runtime_status lib.RuntimeStatus

	nsInMs := float64(time.Millisecond)

	for {
		runtime.ReadMemStats(memStats)

		now := time.Now()

        runtime_status.TimeString = lib.GetTimeString()

        runtime_status.Alloc = memStats.Alloc
        runtime_status.Sys= memStats.Sys
        runtime_status.Mallocs= memStats.Mallocs
        runtime_status.Frees = memStats.Frees

        runtime_status.HeapAlloc = memStats.HeapAlloc
        runtime_status.HeapSys = memStats.HeapSys
        runtime_status.HeapIdle = memStats.HeapIdle
        runtime_status.HeapInuse = memStats.HeapInuse
        runtime_status.HeapObjects = memStats.HeapObjects

        runtime_status.StackSys = memStats.StackSys
        runtime_status.StackInuse = memStats.StackInuse
        runtime_status.MSpanInuse = memStats.MSpanInuse
        runtime_status.MSpanSys = memStats.MSpanSys
        runtime_status.MCacheInuse = memStats.MCacheInuse
        runtime_status.MCacheSys = memStats.MCacheSys

        runtime_status.GCTotalPause = float64(memStats.PauseTotalNs)/nsInMs

        /*
            计算平均每秒GC暂停的时间
        */

        // 总的GC暂停时间大于0
		if lastPauseNs > 0 {
            // 获取两次GC间隔的时间
			pauseSinceLastSample := memStats.PauseTotalNs - lastPauseNs
			runtime_status.GCPausePerSecond = float64(pauseSinceLastSample)/nsInMs/sleep.Seconds()
		}
		lastPauseNs = memStats.PauseTotalNs

        // 两次采样之间GC的次数
		countGc := int(memStats.NumGC - lastNumGc)

        /*
            计算平均每秒GC的次数
        */

		if lastNumGc > 0 {
			diff := float64(countGc)
			diffTime := now.Sub(lastSampleTime).Seconds()
			runtime_status.GCPerSecond = diff/diffTime
		}

        /*
            如果两次采样之间发生了GC，则获取暂停时间(ms)
        */

		if countGc > 0 {
            /*
                如果两次采样之间发生了超过256次的GC，就有可能获取不到一些GC动作的暂停时间
                因为memStats.PauseNS是一个环形的缓存，最多只能记录256次GC动作的暂停时间
            */
			if countGc > 256 {
				fmt.Fprintf(os.Stderr, "We're missing some gc pause times")
				countGc = 256
			}

			for i := 0; i < countGc; i++ {
                // 最近一次的GC暂停时间保存在[(NumGC+255)%256]
				idx := int((memStats.NumGC-uint32(i))+255) % 256
				pause := float64(memStats.PauseNs[idx])
                // 暂停时间单位为ns，ns/nsInMs单位为ms
				runtime_status.GCPause = pause/nsInMs
			}
		}

		// 记录前一次的状态，为下次采样做准备
		lastNumGc = memStats.NumGC
		lastSampleTime = now

        channel <- &runtime_status

		time.Sleep(sleep)
	}
}

