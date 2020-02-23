package main

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"testing"
	"time"
)

func TestTempFileName(t *testing.T) {
	tmpDir := os.TempDir()
	filename := TempFileName(tmpDir, ".apk")
	t.Log(filename)
}

func TestStrings(t *testing.T){
	l := logrus.WithField("remoteaddr", "127.0.0.1")
	l.Debug("debug\n")
	error := errors.New("error: name is null")

	l.WithError(error).Error("1Unable to upgrade connection\n")
	fmt.Println("1```=========================")
	UpdateTime := time.Now()
	UpdateTime2 := UpdateTime.Add(60*time.Second)
	fmt.Println(UpdateTime,UpdateTime2,UpdateTime2.Before(time.Now()))
	fmt.Println("2```=========================")
	procStat := string(`cpu  8211422 2392858 7281974 11527978 23046 1147454 506220 0 0 0
	cpu0 1953480 932043 1958578 11516005 23014 334601 116773 0 0 0`)
	idx := strings.Index(procStat, "\n")
	fmt.Println(idx)
	fields := strings.Fields(procStat[:idx])
	fmt.Println(fields)
	var total, idle uint
	for i, raw := range fields[1:] {
		var v uint
		fmt.Sscanf(raw, "%d", &v)
		if i == 3 { // idle
			idle = v
		}
		total += v
	}
	fmt.Println(total,idle)
}

//var ch = make(chan string, 10) // 创建大小为 10 的缓冲信道

func download(url string) {
	fmt.Println("start to download", url)
	time.Sleep(1*time.Second)
	//fmt.Println("send", url)
	//ch <- url // 将 url 发送给信道
}


func TestChannal(t *testing.T) {

	ch := time.Tick(2 * time.Second)
	i := 0
	for range ch {
		if i == 3{
			break
		}
		download("a.com/" + string(i+'0'))
		i++
	}
	//for i := 0; i < 3; i++ {
	//	go download("a.com/" + string(i+'0'))
	//}
	//for range ch {
	//	msg := <-ch // 等待信道返回消息。
	//	fmt.Println("finish", msg)
	//}

	//for i := 0; i < 3; i++ {
	//	msg := <-ch // 等待信道返回消息。
	//	fmt.Println("finish", msg)
	//}
	fmt.Println("Done!")
}
