package main

import (
	"ConcurrencyCron"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

/**
 *@author  wxn
 *@project ConcurrencyCron
 *@package ConcurrencyCron
 *@date    19-8-1 下午5:58
 */

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
	//ConcurrencyCron.DefaultWriter=os.Stdout
	scheduler, err = ConcurrencyCron.NewScheduler(200)
	if err != nil {
		fmt.Println(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	scheduler.Every(1).Minutes().Do(test, time.Now().Format("15:04"))
	scheduler.Start(ctx)
}

func main() {
	//scheduler, err := ConcurrencyCron.NewScheduler(200) //200 is the number of tasks that can be run in parallel
	//if err != nil {
	//	fmt.Println(err)
	//}
	//for i := 0; i < 200; i++ {
	//scheduler.Every(1).Seconds().Do(test, 1)
	//scheduler.Every(1).Minutes().Do(test, 1000+1)
	//scheduler.Every(1).Hours().Do(test, 10000+1)
	//scheduler.Once().At("9:38").Do(test, 10000+2)
	//}
	//fmt.Println("started")
	//ctx, _ := context.WithCancel(context.Background())
	//scheduler.Start(ctx)
	//ch := make(chan bool)
	//<-ch             //test
	//scheduler.Stop() //stop the tasks
	//cancel()
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
	r.GET("/list", func(c *gin.Context) {
		res := scheduler.ListTasks()
		c.JSON(200, res)
	})
	r.Run(":12315")
	ch := make(chan bool)
	<-ch //test
}
