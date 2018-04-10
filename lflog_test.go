package lflog

import (
	"testing"
	"time"
	"runtime"
	"fmt"
)

func TestLog(t *testing.T) {

	LoadConfig("config.xml")
	//defer Close()
	date := time.Now().Format("2006-01-02")

	Debug("这是debug,当前时间:%s",date)

	Info("uid:%d,name:%s",10057,"jfeige")

	runtime.Gosched()

	fmt.Println("finish....")
}
