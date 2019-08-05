## ConcurrentCron: A Golang Job Scheduling Package.

[![GgoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/whutwxn/ConcurrencyCron)

ConcurrentCron is a task scheduler that supports high concurrency at the same time which lets you run Go functions periodically at pre-determined interval using a simple, human-friendly syntax.

You can run this scheduler in the following way

```go

func test(num int) {
	fmt.Println("before:im a task", num)
	time.Sleep(10 * time.Second)
	fmt.Println("after:im a task", num, time.Now())
}

func main() {
	scheduler, err := ConcurrencyCron.NewScheduler(200) //200 is the number of tasks that can be run in parallel
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 200; i++ {
		scheduler.Every(1).Seconds().Do(test, i)
		scheduler.Every(1).Minutes().Do(test, 1000+i)
		scheduler.Every(1).Hours().Do(test, 10000+i)  you can refer to [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)
	}
	ctx, cancel := context.WithCancel(context.Background())
	scheduler.Start(ctx)
	ch := make(chan bool)
	<-ch //test
	defer cancel() //you can shutdown the tasks by this function gratefully

}


This article refers to some of [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)'s ideas and things, the specific timing tasks are the same as gocron, you can refer to his project
