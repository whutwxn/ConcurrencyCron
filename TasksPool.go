package ConcurrencyCron

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"runtime"
	"time"
)

/**
 *@author  wxn
 *@project ConcurrencyCron
 *@package ConcurrencyCron
 *@date    19-8-2 上午11:23
 */
type TasksPool interface {
	Seconds() *task                                        //Run every few seconds
	Minutes() *task                                        //Run every few minutes
	Hours() *task                                          //Run every few hours
	Days() *task                                           //Run every few days
	Weekdays() *task                                       //Run every few weeks
	Monday() *task                                         //Run every few weeks on Monday
	Tuesday() *task                                        //Run every few weeks on Tuesday
	Wednesday() *task                                      //Run every few weeks on Wednesday
	Thursday() *task                                       //Run every few weeks on Thursday
	Friday() *task                                         //Run every few weeks on Friday
	Saturday() *task                                       //Run every few weeks on Saturday
	Sunday() *task                                         //Run every few weeks on Sunday
	JudgeRun(tm time.Time) bool                            //Determine if it is going to run
	Run(ticket TicketsPool, tm time.Time)                  //Run and judge the next run time
	GetNext() time.Time                                    //Get nest run time
	Do(taskFunc interface{}, params ...interface{}) string //Add a run function
	GetUuid() string                                       //get uuid
}

type task struct {
	uuid      string        //uuid for each task
	interval  uint64        //interval foreach task
	atTime    time.Duration //The time of the task starts running e.g: 8:05
	latest    time.Time     //The last run time of the task
	next      time.Time     //The next run time of the task
	startDay  time.Weekday  //The weeks of the task starts running e.g: Monday
	funcName  string        //The name of the function that needs to be run
	funcVal   interface{}   //the function that needs to be run
	funcParam []interface{} //Function parameters
	unit      string        //The unit of running interval
}

//create a task
func NewTask(interval uint64) TasksPool {
	jobId, err := uuid.NewV4()
	if err != nil {
		fmt.Println("get uuid error")
	}
	id := jobId.String()
	return &task{
		id,
		interval,
		0,
		time.Unix(0, 0),
		time.Unix(0, 0),
		time.Sunday,
		"",
		nil,
		nil,
		"",
	}
}

func (t *task) setUnit(unit string) *task {
	t.unit = unit
	return t
}

func (t *task) weekday(startDay time.Weekday) *task {
	t.startDay = startDay
	return t.Weekdays()
}

func (t *task) periodDuration() time.Duration {
	interval := time.Duration(t.interval)
	switch t.unit {
	case "seconds":
		return interval * time.Second
	case "minutes":
		return interval * time.Minute
	case "hours":
		return interval * time.Hour
	case "days":
		return time.Duration(interval * time.Hour * 24)
	case "weeks":
		return time.Duration(interval * time.Hour * 24 * 7)
	}
	panic("unspecified job period")
}

func (t *task) getNextRun() {
	now := time.Now()
	if t.latest == time.Unix(0, 0) {
		t.latest = now
	}
	switch t.unit {
	case "seconds":
		t.next = t.latest.Add(time.Duration(t.interval) * time.Second)
	case "minutes":
		t.next = t.latest.Add(time.Duration(t.interval) * time.Minute)
	case "hours":
		t.next = t.latest.Add(time.Duration(t.interval) * time.Hour)
	case "days":
		t.next = time.Date(t.latest.Year(), t.latest.Month(), t.latest.Day(), 0, 0, 0, 0, time.Local)
		t.next = t.next.Add(t.atTime)
	case "weeks":
		t.next = time.Date(t.latest.Year(), t.latest.Month(), t.latest.Day(), 0, 0, 0, 0, time.Local)
		dayDiff := int(t.startDay)
		dayDiff -= int(t.next.Weekday())
		if dayDiff != 0 {
			t.next = t.next.Add(time.Duration(dayDiff) * 24 * time.Hour)
		}
		t.next = t.next.Add(t.atTime)

	}

	for t.next.Before(now) || t.next.Before(t.latest) {
		t.next = t.next.Add(t.periodDuration())
	}
}

func (t *task) GetNext() time.Time {
	return t.next
}

func (t *task) JudgeRun(tm time.Time) bool {
	return tm.Unix() >= t.next.Unix()
}

func (t *task) Run(ticket TicketsPool, tm time.Time) {
	defer func() {
		if p := recover(); p != nil {
			err, ok := interface{}(p).(error)
			var errMsg string
			if ok {
				errMsg = fmt.Sprintf("Async Call Panic! (error: %s)", err)
			} else {
				errMsg = fmt.Sprintf("Async Call Panic! (clue: %#v)", p)
			}
			fmt.Println(errMsg)
		}
		ticket.Return()
	}()
	taskFunc := reflect.ValueOf(t.funcVal)
	if len(t.funcParam) != taskFunc.Type().NumIn() {
		fmt.Printf("param number error:need %d params given %d params", taskFunc.Type().NumIn(), len(t.funcParam))
	}
	params := make([]reflect.Value, len(t.funcParam))
	for i, param := range t.funcParam {
		params[i] = reflect.ValueOf(param)
	}
	t.latest = tm
	t.getNextRun()
	taskFunc.Call(params)

}

func (t *task) Do(taskFunc interface{}, params ...interface{}) string {
	tp := reflect.TypeOf(taskFunc)
	if tp.Kind() != reflect.Func {
		panic("only function can be schedule into the job queue.")
	}
	t.funcName = runtime.FuncForPC(reflect.ValueOf(taskFunc).Pointer()).Name()
	t.funcVal = taskFunc
	t.funcParam = params
	t.getNextRun()
	return t.uuid
}

func (t *task) Seconds() *task {
	return t.setUnit("seconds")
}

func (t *task) Minutes() *task {
	return t.setUnit("minutes")
}

func (t *task) Hours() *task {
	return t.setUnit("hours")
}

func (t *task) Days() *task {
	return t.setUnit("days")
}

func (t *task) Weekdays() *task {
	return t.setUnit("weeks")
}

func (t *task) Monday() (job *task) {
	return t.weekday(time.Monday)
}

func (t *task) Tuesday() *task {
	return t.weekday(time.Tuesday)
}

func (t *task) Wednesday() *task {
	return t.weekday(time.Wednesday)
}

func (t *task) Thursday() *task {
	return t.weekday(time.Thursday)
}

func (t *task) Friday() *task {
	return t.weekday(time.Friday)
}

func (t *task) Saturday() *task {
	return t.weekday(time.Saturday)
}

func (t *task) Sunday() *task {
	return t.weekday(time.Sunday)
}

func (t *task) GetUuid() string {
	return t.uuid
}
