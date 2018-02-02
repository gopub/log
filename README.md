# Log utility for Go
### Level  
``` 
log.Debug("This is a debug message")
log.Infof("count:%d", 3)
```
Output:
```
[DEBUG] 2018/02/02 15:57:16 github.com/g/l/log_test.go(TestLogger_Debug):9      | This is a debug message
[INFO]  2018/02/02 15:57:16 github.com/g/l/log_test.go(TestLogger_Debug):10     | count:3
```
Disable output by setting `log.SetLevel(log.INFO)`
### Fields
``` 
logger := log.WithFields([]*log.Field{{Key: "userID", Value: 1}, {Key: "name", Value: "Tom"}})
logger.Error("data not found")

logger.WithFields([]*log.Field{{Key: "count", Value: 2}}).Infof("Try to post topic:%s", "Which is the best city")
```
Output:
``` 
[ERROR] 2018/02/02 15:57:16 github.com/g/l/log_test.go(TestFieldLogger_WithFields):15   | userID=1 name=Tom     | data not found
[INFO]  2018/02/02 15:57:16 github.com/g/l/log_test.go(TestFieldLogger_WithFields):17   | userID=1 name=Tom count=2     | Try to post topic:Which is the best city
```

### Flags
Set flags to filter log info
``` 
log.SetFlags(log.Lmicroseconds|log.Lfunction)
log.Info("System started")
```
Output:
``` 
[INFO]  16:02:25.289557 github.com/gopub/log_test.TestLogger_SetFlags:22        | System started
```