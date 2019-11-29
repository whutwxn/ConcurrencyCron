## ConcurrentCron: A Golang Job Scheduling Package.

[![GgoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/whutwxn/ConcurrencyCron)
[![Go ReportCard](https://goreportcard.com/report/github.com/whutwxn/ConcurrencyCron)](https://goreportcard.com/report/github.com/whutwxn/ConcurrencyCron)

ConcurrentCron is a task scheduler that supports high concurrency at the same time which lets you run Go functions periodically at pre-determined interval using a simple, human-friendly syntax.

You can run this scheduler in the following way

```go
package main

import (
	"ConcurrencyCron"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

var (
	scheduler ConcurrencyCron.Scheduler
)

func test(num string) {
	//fmt.Println("before:im a task:", num)
	//time.Sleep(10 * time.Second)
	fmt.Println("after:im a task:", num, " current time:", time.Now().Format("15:04"))
}

func init() {
	var err error
	scheduler, err = ConcurrencyCron.NewScheduler(200, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	scheduler.Every(1).Minutes().Do(test, time.Now().Format("15:04"))
	scheduler.Start(ctx)
}

func main() {
	r := gin.Default()
	r.PUT("/addOnce", func(c *gin.Context) {
		tm := time.Now()
		hour := tm.Hour()
		min := tm.Minute() + 1
		tim := fmt.Sprintf("%2d:%2d", hour, min)
		uuid := scheduler.Once().At(tim).Do(test, tim)
		fmt.Println(uuid)
		c.String(200, uuid)
	})
	r.PUT("/addInterval", func(c *gin.Context) {

	})
	r.DELETE("/removeOnce/:uuid", func(c *gin.Context) {
		uuid := c.Param("uuid")
		scheduler.RemoveByUuid(uuid)
	})
	r.Run(":12315")
	ch := make(chan bool)
	<-ch //test
}
```
This article refers to some of [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)'s ideas and things, the specific timing tasks are the same as gocron, you can refer to his project

Thank you for the support and understanding ,jasonlvhit!
