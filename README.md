# lflog
模仿log4go实现的一款很简单的日志库，还在继续完善中

安装:
```
go get github.com/jfeige/lflog
```
使用:

配置文件格式参考config.xml

应用程序入口加载配置文件:
```
        LoadConfig("config.xml")
```
使用:
```
        Info("当前时间:%s",time.Now().Format("2006-01-02 15:04:05"))
        Info("开始读取数据......")
        
        args := []interface{}{10057,100,"活跃用户"}
        Debug("uid:%d,cnt:%d,memo",args...)
        Debug("receive a post request")
```
API列表:
```
        Debug(args0 interface{}, args ...interface{})

        Info(args0 interface{}, args ...interface{})

        Warn(args0 interface{}, args ...interface{})

        Error(args0 interface{}, args ...interface{})
```
