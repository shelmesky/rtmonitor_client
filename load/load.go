package load

import (
	"bytes"
	"github.com/shelmesky/rtmonitor_client/lib"
	"log"
	"strconv"
	"strings"
	"time"
)

func GetLoadInfo(channel chan *lib.LoadInfo, sleep time.Duration) {
	var load_info lib.LoadInfo

	for {
		meminfo_file := "/proc/meminfo"

		meminfo_data, err := lib.ReadFileFullContentByLine(meminfo_file)

		if err != nil {
			log.Fatal(err)
			return
		}

		meminfo_total_splited := bytes.Fields(meminfo_data[0])
		meminfo_free_splited := bytes.Fields(meminfo_data[1])
		meminfo_buffers_splited := bytes.Fields(meminfo_data[2])
		meminfo_cached_splited := bytes.Fields(meminfo_data[3])

		if v, err := strconv.ParseUint(string(meminfo_total_splited[1]), 10, 0); err == nil {
			load_info.MemTotal = v
		}

		if v, err := strconv.ParseUint(string(meminfo_free_splited[1]), 10, 0); err == nil {
			load_info.MemFree = v
		}

		if v, err := strconv.ParseUint(string(meminfo_buffers_splited[1]), 10, 0); err == nil {
			load_info.MemBuffers = v
		}

		if v, err := strconv.ParseUint(string(meminfo_cached_splited[1]), 10, 0); err == nil {
			load_info.MemCached = v
		}

		load_info.MemUsed = load_info.MemTotal - load_info.MemFree

		loadavg_file := "/proc/loadavg"

		loadavg_data, err := lib.ReadFileFullContent(loadavg_file)

		if err != nil {
			log.Fatal(err)
			return
		}

		loadavg_splited := strings.Split(string(loadavg_data), " ")

		if len(loadavg_splited) > 0 {
			// 1分钟内的负载
			load1 := strings.Trim(loadavg_splited[0], " ")
			if v, err := strconv.ParseFloat(load1, 64); err == nil {
				load_info.LoadAVG1 = v
			}

			// 5分钟内的负载
			load5 := strings.Trim(loadavg_splited[1], " ")
			if v, err := strconv.ParseFloat(load5, 64); err == nil {
				load_info.LoadAVG5 = v
			}

			// 15分钟内的负载
			load15 := strings.Trim(loadavg_splited[2], " ")
			if v, err := strconv.ParseFloat(load15, 64); err == nil {
				load_info.LoadAVG15 = v
			}
		}

		channel <- &load_info
		time.Sleep(sleep)
	}
}
