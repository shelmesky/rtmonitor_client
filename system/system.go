package system

import (
	"github.com/shelmesky/rtmonitor_client/lib"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
)

func GetSystemInfo() (*lib.SystemInfo, error) {
	var system_info lib.SystemInfo
	var uname_info syscall.Utsname

	if err := syscall.Uname(&uname_info); err != nil {
		log.Fatal(err)
		return nil, err
	}

	system_info.Hostname = string(lib.ConvertToSlice(uname_info.Nodename))
	system_info.Sysname = string(lib.ConvertToSlice(uname_info.Sysname))
	system_info.Release = string(lib.ConvertToSlice(uname_info.Release))
	system_info.Machine = string(lib.ConvertToSlice(uname_info.Machine))

	cpuinfo_file := "/proc/cpuinfo"

	cpuinfo_data, err := lib.ReadFileFullContentByLine(cpuinfo_file)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	cpuinfo_model_name_splited := strings.Split(string(cpuinfo_data[4]), ":")
	cpuinfo_frequency_splited := strings.Split(string(cpuinfo_data[6]), ":")

	cpuinfo_model_name := strings.Trim(cpuinfo_model_name_splited[1], " ")
	cpuinfo_frequency := strings.Trim(cpuinfo_frequency_splited[1], " ")

	system_info.CPUModelName = cpuinfo_model_name
	system_info.CPUFrequency = cpuinfo_frequency

	system_info.CPUCores = runtime.NumCPU()

	system_info.Location = lib.GetTimeLocation()
	system_info.GOVersion = runtime.Version()
	system_info.ProcessID = os.Getpid()

	return &system_info, nil
}
