package lflog

import (
	"testing"
	"time"
	"fmt"
)

func TestLog(t *testing.T) {

	LoadConfig("config.xml")
	//defer Close()
	date := time.Now().Format("2006-01-02")

	Debug("这是debug,当前时间:%s",date)

	Info("uid:%d,name:%s",10057,"jfeige")

	args := []interface{}{10057,100,"活跃用户"}
	Debug("uid:%d,cnt:%d,memo:%s",args...)


	time.Sleep(2*time.Second)

	fmt.Println("finish....")
}
