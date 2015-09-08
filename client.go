package rtmonitor_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shelmesky/rtmonitor_client/lib"
	"github.com/shelmesky/rtmonitor_client/load"
	"github.com/shelmesky/rtmonitor_client/process"
	rtruntime "github.com/shelmesky/rtmonitor_client/runtime"
	"github.com/shelmesky/rtmonitor_client/system"
)

var (
	ReportLoadInterval    = 60 * time.Second
	ReportProcessInterval = 1 * time.Second
	ReportRuntimeInterval = 1 * time.Second

	load_channel    chan *lib.LoadInfo      = make(chan *lib.LoadInfo, 100)
	process_channel chan *lib.ProcessInfo   = make(chan *lib.ProcessInfo, 100)
	runtime_channel chan *lib.RuntimeStatus = make(chan *lib.RuntimeStatus, 100)

	apiServer string
	clientKey string
)

func Start(api_server string, client_key string) {
	apiServer = api_server
	clientKey = client_key

	go rtruntime.GetRuntimeStats(runtime_channel, ReportRuntimeInterval)
	go ReportRuntimeInfo()

	go process.GetProcessInfo(process_channel, ReportProcessInterval)
	go ReportProcessInfo()

	go load.GetLoadInfo(load_channel, ReportLoadInterval)
	go ReportLoadInfo()

	ReportSystemInfo()
}

// 报告系统信息，进程启动时报告一次
func ReportSystemInfo() {
	system_info, err := system.GetSystemInfo()
	if err != nil {
		log.Println(err)
	}

	data, err := json.Marshal(system_info)
	if err != nil {
		log.Println(err)
	}

	err = SendToServer(data, "system")
	if err != nil {
		log.Println(err)
	}

}

// 报告系统负载信息
func ReportLoadInfo() {
	for {
		if load_status, ok := <-load_channel; ok {
			data, err := json.Marshal(load_status)
			if err != nil {
				log.Println(err)
			}

			err = SendToServer(data, "load")
			if err != nil {
				log.Println(err)
			}
		} else {
			break
		}
	}
}

// 报告进程的信息
func ReportProcessInfo() {
	for {
		if process_status, ok := <-process_channel; ok {
			data, err := json.Marshal(process_status)
			if err != nil {
				log.Println(err)
			}

			err = SendToServer(data, "process")
			if err != nil {
				log.Println(err)
			}
		} else {
			break
		}
	}
}

// 报告Golang运行时的信息
func ReportRuntimeInfo() {
	for {
		if runtime_status, ok := <-runtime_channel; ok {
			data, err := json.Marshal(runtime_status)
			if err != nil {
				log.Println(err)
			}

			err = SendToServer(data, "runtime")
			if err != nil {
				log.Println(err)
			}
		} else {
			break
		}
	}
}

func SendToServer(data []byte, data_type string) error {
	request_url := apiServer + clientKey + "/report/" + data_type + "/"
	buf := bytes.NewReader(data)

	resp, err := http.Post(request_url, "application/json", buf)
	if err != nil {
		log.Println(err)
		return err
	}

	if resp != nil {
		resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("Server return none 200 code: %d", resp.StatusCode)
		}
	}

	return nil
}
