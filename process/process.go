package process

import (
    "fmt"
    "os"
    "time"
    "log"
    "strings"
    "strconv"
    "github.com/shelmesky/rtmonitor_client/lib"
)

// 获取进程的信息
func GetProcessInfo(channel chan *lib.ProcessInfo, sleep time.Duration) {
    /*
        Process's memory status in /proc/PID/statm:
        total program pages |
        resident pages |
        shared pages |
        text(code) |
        data/stack |
        library |
        dirty pages

        page size is 4K
    */

    var process_info lib.ProcessInfo

    pagesize := int64(os.Getpagesize())

    for {
        pid := os.Getpid()

        statm_file := fmt.Sprintf("/proc/%d/statm", pid)

        data, err := lib.ReadFileFullContent(statm_file)

        if err != nil {
            log.Fatal(err)
            return
        }

        now_time := lib.GetTimeString()
        process_info.TimeString = now_time

        memory_splited := strings.Split(string(data), " ")

        if len(memory_splited) > 0 {
            // 虚拟内存(VIRT)
            virtual_pages := strings.Trim(memory_splited[0], " ")
            if v, err := strconv.ParseInt(virtual_pages, 10, 0); err == nil {
                process_info.VirtualMemory = v * pagesize
            }

            // 物理内存(RES)
            resisdent_pages := strings.Trim(memory_splited[1], " ")
            if v, err := strconv.ParseInt(resisdent_pages, 10, 0); err == nil {
                process_info.ResisdentMemory = v * pagesize
            }

            // 共享内存(SHR)
            shared_pages := strings.Trim(memory_splited[2], " ")
            if v, err := strconv.ParseInt(shared_pages, 10, 0); err == nil {
                process_info.SharedMemory = v * pagesize
            }
        }

        // 进程运行了多长时间
        process_info.Uptime = int64(time.Now().Sub(lib.StartTime).Seconds())

        channel <- &process_info

        time.Sleep(sleep)
    }
}


