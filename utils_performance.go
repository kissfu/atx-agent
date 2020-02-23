package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PerformanceInfo struct {
	CpuPercent  float64 `json:"cpuPercent"`
	MemoPercent float64 `json:"memoPercent"`
}

func getMemoryInfo() (info map[string]int, err error) {
	output, err := Command{
		Args:    []string{"cat", "proc/meminfo"},
		Timeout: 10 * time.Second,
	}.CombinedOutputString()
	if err != nil {
		return
	}
	re := regexp.MustCompile(`(\w[\w ]+):\s*(\d+)`)
	matches := re.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		err = errors.New("Invalid dumpsys meminfo output")
		return
	}
	info = make(map[string]int, len(matches))
	for _, m := range matches {
		key := strings.ToLower(m[1])
		val, _ := strconv.Atoi(m[2])
		info[key] = val
	}
	return
}

func readPerformanceInfo() (info PerformanceInfo, err error) {
	// android进程从1开始。
	last, ok := cpuStats[0]
	if !ok || // need fresh history data
		last.UpdateTime.Add(5*time.Second).Before(time.Now()) {
		last, err = NewCPUStat(0)
		if err != nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
		log.Println("Update data")
	}
	stat, err := NewCPUStat(0)
	if err != nil {
		return
	}
	memo, err := getMemoryInfo()
	cpuStats[0] = stat
	info.CpuPercent = 100.0 - stat.SystemCPUPercent(last)
	info.MemoPercent = 100.0 - (float64(memo["memfree"])+float64(memo["cached"]))/float64(memo["memtotal"])*100.0
	return
}

func handlePerformanceWebsocket(w http.ResponseWriter, r *http.Request) {
	l := logrus.WithField("remoteaddr", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l.WithError(err).Error("Unable to upgrade connection")
		return
	}
	defer func() {
		log.Println("finally!")
		if r := recover(); r != nil {
			log.Println("conn close", r)
			conn.Close()
		}
	}()
	//go func() {
	//for {
	ch := time.Tick(5 * time.Second)
	for range ch {
		info, err := readPerformanceInfo()
		if err != nil {
			l.WithError(err).Error("read performance error")
			return
		}
		log.Println("conn.WriteJSON", info)
		if err := conn.WriteJSON(info); err != nil {
			l.WithError(err).Error("write json error")
			return
		}
	}
	//}
	//}()
	log.Println("Done!")
}
