
# ClessLi/log

`ClessLi/log`是是一个生产可用的日志包，基于 `zap` 包封装。除了实现 `Go` 日志包的基本功能外，还实现了很多高级功能，`ClessLi/log`具有如下特性：

- 支持自定义配置。
- 支持文件名和行号。
- 支持日志级别 Debug、Info、Warn、Error、Panic、Fatal。
- 支持输出到本地文件和标准输出。支持 JSON 和 TEXT 格式的日志输出，支持自定义日志格式。
- 支持选项模式。
## Usage/Examples

```go
import log "github.com/ClessLi/ansible-role-manager/pkg/log/v1"

func main() {
    defer log.Flush()

    // Debug、Info(with field)、Warnf、Errorw使用
    log.Debug("this is a debug message")
    log.Info("this is a info message", log.Int32("int_key", 10))
    log.Warnf("this is a formatted %s message", "warn")
}
```

执行代码：

```bash
$ go run example.go
```

  
## Acknowledgements

 - [Awesome Readme Templates](https://awesomeopensource.com/project/elangosundar/awesome-README-templates)
 - [Awesome README](https://github.com/matiassingers/awesome-readme)
 - [How to write a Good readme](https://bulldogjob.com/news/449-how-to-write-a-good-readme-for-your-github-project)

  
## License

[MIT](https://choosealicense.com/licenses/mit/)

  