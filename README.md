# Log utility for Go
### Level  
``` 
log.Debug("This is a debug message")
log.Infof("count:%d", 3)
```
Output:
```
2018/02/12 09:09:52 [DEB] g/g/l/log_test.go(TestLogger_Debug):9      | This is a debug message
2018/02/12 09:09:52 [INF]  g/g/l/log_test.go(TestLogger_Debug):10     | count:3
```
Disable output by setting `log.SetLevel(log.INFO)`
### Fields
``` 
logger := log.With("userID", 1, "name", "Tom")
logger.Error("data not found")

logger.With("count", 2).Infof("Try to post topic:%s", "Which is the best city")
```
Output:
``` 
2018/02/12 09:09:52 [ERR] g/g/l/log_test.go(TestFieldLogger_WithFields):15   | userID=1 name=Tom     | data not found
2018/02/12 09:09:52 [INF] g/g/l/log_test.go(TestFieldLogger_WithFields):17   | userID=1 name=Tom count=2     | Try to post topic:Which is the best city
```

### Flags
Set flags to filter log info
``` 
log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lfunction)
log.Info("System started")
```
Output:
``` 
2018/02/12 09:11:26.558626 [INF]  g/g/log_test.TestLogger_SetFlags:22     | System started
```

### Module Name
Give logger a name, which make the log clearer
``` 
l := log.GetLogger("TestModule")
l.Info("This is a log")
```
Output:
``` 
2019-01-13 12:53:47.545+0800 [INF] [TestModule] g/g/l/log_test.go(TestLogger_GetLogger):30      | This is a log
```

### Custom log output
All logs are output to os.Stderr by default, however you could change the output destination
``` 
w, err := os.OpenFile(fileName)
...
log.Default().SetOutput(w)
```

### Write logs into files
If environment value LOG_DIR is defined, logs will be saved into files under LOG_DIR. The format of file name is yyyyMMdd.{Num}.log. E.g. 20200118.1.log.