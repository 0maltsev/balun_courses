package main

import (
    "container/heap"
    "testing"

    "github.com/stretchr/testify/assert"
)

type Task struct {
    Identifier int
    Priority   int
}

type PriorityQueue []*Task

type Scheduler struct {
    pq *PriorityQueue
}

func NewScheduler() Scheduler {
    queue := &PriorityQueue{}
    heap.Init(queue)

    return Scheduler{
        pq: queue,
    }
}

func (s *Scheduler) AddTask(task Task) {
    heap.Push(s.pq, &task)
}

func (s *Scheduler) GetTask() Task {
    return *heap.Pop(s.pq).(*Task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
    s.pq.Update(taskID, newPriority)
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(i interface{}) {
    task := i.(*Task)
    *pq = append(*pq, task)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    task := old[n-1]
    *pq = old[:n-1]
    return task
}

func (pq *PriorityQueue) Update(taskID int, newPriority int) {
    for i, task := range *pq {
        if task.Identifier == taskID {
            task.Priority = newPriority
            heap.Fix(pq, i)
            return
        }
    }
}

func TestTrace(t *testing.T) {
    task1 := Task{Identifier: 1, Priority: 10}
    task2 := Task{Identifier: 2, Priority: 20}
    task3 := Task{Identifier: 3, Priority: 30}
    task4 := Task{Identifier: 4, Priority: 40}
    task5 := Task{Identifier: 5, Priority: 50}

    scheduler := NewScheduler()
    scheduler.AddTask(task1)
    scheduler.AddTask(task2)
    scheduler.AddTask(task3)
    scheduler.AddTask(task4)
    scheduler.AddTask(task5)

    task := scheduler.GetTask()
    assert.Equal(t, task5, task)

    task = scheduler.GetTask()
    assert.Equal(t, task4, task)

    scheduler.ChangeTaskPriority(1, 100)

    task = scheduler.GetTask()
    assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

    task = scheduler.GetTask()
    assert.Equal(t, task3, task)
}