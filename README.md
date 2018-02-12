# Log utility for Go
### Level  
``` 
log.Debug("This is a debug message")
log.Infof("count:%d", 3)
```
Output:
```
2018/02/12 09:09:52 [DEBUG] github.com/g/l/log_test.go(TestLogger_Debug):9      | This is a debug message
2018/02/12 09:09:52 [INFO]  github.com/g/l/log_test.go(TestLogger_Debug):10     | count:3
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
2018/02/12 09:09:52 [ERROR] github.com/g/l/log_test.go(TestFieldLogger_WithFields):15   | userID=1 name=Tom     | data not found
2018/02/12 09:09:52 [INFO]  github.com/g/l/log_test.go(TestFieldLogger_WithFields):17   | userID=1 name=Tom count=2     | Try to post topic:Which is the best city
```

### Flags
Set flags to filter log info
``` 
log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lfunction)
log.Info("System started")
```
Output:
``` 
2018/02/12 09:11:26.558626 [INFO]  github.com/gopub/log_test.TestLogger_SetFlags:22     | System started
```