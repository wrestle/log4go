# log4go

## Description

This repository is reconstructed from alecthomas's log4go, which is a logging package similar to log4j for the Go programming language.

Two new features are supported, one is Json config style, and the other is deferent output accordding to category.

## Features

-   **Log to console**
-   **Log to file, support rotate by size or time**
-   **log to network, support tcp and udp**
-   **support xml config**
-   **原项目使用过于繁琐，经过一些小改变，开箱即用，日常使用我的默认设置即可，无需自行调整配置**
-   **默认支持按天滚动，每个日志文件 800M 上限，超过即滚动**

---------------------------

-   **Support Json style configuration**
-   **Add Category for log**
    * Classify your logs for different output and different usage.
-   **Compatible with the old**

## Usage

First, get the code from this repo. 

```go get github.com/wrestle/log4go```

Then import it to you project.

```import log "github.com/wrestle/log4go"```



Code example:

```
package main

import (
	log "github.com/jeanphorn/log4go"
        "os"	
)

func main() {
        // 感谢 alecthomas 和 jeanphorn 提供的基础项目

	// 默认使用 INFO 级别以上的日志
	// 日志位置在 当前目录的上级中的 log 里
	// e.g : 假设执行目录是 /data/code/test_program_dir/bin/exec_program_bin
	// 那么必须有 /data/code/test_program_dir/log 目录
	// 日志会打在这文件夹里面，文件名为 /data/code/test_program_dir/log/exec_program_bin.log
        log.SetUniqueLogName(os.Args[0])
	log.Info("nomal info test ...")

	log.Close()
}

```

 


## Thanks

Thanks alecthomas for providing the [original resource](https://github.com/alecthomas/log4go).
