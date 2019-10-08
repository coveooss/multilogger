# multilogger
Multilogger is a small utility based on [logrus](https://github.com/sirupsen/logrus) that writes logs simultaneously to the console and to a file. Each of those logging methods can be given a different logging level. It exposes the same API as a regular logrus logger

## Example
```go
log = multilogger.New(logrus.InfoLevel, logrus.DebugLevel, "my_log_file.out", "app_name")
log.Info("test")
log.Debug("test")
log.Fatal("test")
```

## Usages
This project is used in internal Coveo projects such as:

* [gotemplate](https://github.com/coveooss/gotemplate)
* [terragrunt](https://github.com/coveooss/tgf)
* [tgf](https://github.com/coveooss/tgf)

It is used to reduce visual clutter in CI systems while keeping debug logs when errors arise