package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/codeskyblue/goreq"
	shellquote "github.com/kballard/go-shellquote"
	"github.com/openatx/androidutils"
	"github.com/pkg/errors"
	"github.com/prometheus/procfs"
	"github.com/shogo82148/androidbinary/apk"
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
	memo, err := getMemoryInfo();
	cpuStats[0] = stat
	info.CpuPercent = 100.0 - stat.SystemCPUPercent(last)
	info.MemoPercent = 100.0 - (float64(memo["memfree"])+float64(memo["cached"]))/float64(memo["memtotal"])*100.0;
	return
}