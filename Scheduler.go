package ConcurrencyCron

import (
	"context"
	"sync"
	"time"
)

/**
 *@author  wxn
 *@project ConcurrencyCron
 *@package ConcurrencyCron
 *@date    19-8-2 下午1:53
 */

type Scheduler interface {
	Every(interval uint64) TasksPool
	Start(ctx context.Context)
	Stop()
	RemoveByUuid(uuid string)
	GetCurrentTicketsNum() uint32
	GetTaskNum() int
}

type scheduler struct {
	mutex   sync.Mutex
	tickets TicketsPool
	tasks   []TasksPool
	size    int
}

func (s *scheduler) GetTaskNum() int {
	return s.size
}

func (s *scheduler) GetCurrentTicketsNum() uint32 {
	return s.tickets.Remain()
}

func NewScheduler(threads uint32) (Scheduler, error) {
	tickets, err := NewTicketsPool(threads)
	if err != nil {
		return nil, err
	}
	return &scheduler{tickets: tickets, tasks: []TasksPool{}, size: 0}, nil
}

func (s *scheduler) Len() int {
	return s.size
}

func (s *scheduler) Swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
}

func (s *scheduler) Less(i, j int) bool {
	return s.tasks[j].GetNext().Unix() >= s.tasks[i].GetNext().Unix()
}

func (s *scheduler) Every(interval uint64) TasksPool {
	task := NewTask(interval)
	s.mutex.Lock()
	//s.tasks[s.size] = task
	s.tasks = append(s.tasks, task)
	s.size++
	defer s.mutex.Unlock()
	return task
}

func (s *scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for true {
			select {
			case <-ticker.C:
				s.startRun()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *scheduler) Stop() {
	s.tickets.Close()
}

func (s *scheduler) RemoveByUuid(uuid string) {
	i := 0
	//found := false
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for ; i < s.size; i++ {
		if s.tasks[i].GetUuid() == uuid {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			//found = true
			break
		}
	}
	//if !found {
	//	return
	//}
	//for j := i + 1; j < s.size; j++ {
	//	s.tasks[i] = s.tasks[j]
	//	i++
	//}
	s.size = s.size - 1
}

func (s *scheduler) startRun() {
	//sort.Sort(s)
	tm := time.Now()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := 0; i < s.size; i++ {
		if s.tasks[i].JudgeRun(tm) {
			s.tickets.Take()
			go s.tasks[i].Run(s.tickets, tm)
		}
	}
}
