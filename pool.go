package main

import (
	"fmt"
	"time"
)

//定义一个任务类型 Task
type Task struct {
	f func() error //定义一个Task里面应该有一个具体的业务， 业务的名称就叫f
	//这里可以加个任务的优先级
}

//创建一个Task任务
func NewTask(arg_f func() error) *Task {
	t := Task{
		f: arg_f,
	}
	return &t
}

//Task也需要一个执行业务的方法
func (t *Task) Execute() {
	t.f() //调用任务中已经绑定好的业务方法
}

//----------有关协程池 Pool角色的功能
//定义一个Pool协程池的类型
type Pool struct {
	//对外的Task入口 EntryChannel
	EntryChannel chan *Task

	//内部的Task队列 JobsChannel
	JobsChannel chan *Task

	//协程池中最大的worker的数量
	worker_num int
}

//创建Pool的函数
func NewPool(cap int) *Pool {
	//创建一个Pool
	p := Pool{
		EntryChannel: make(chan *Task),
		JobsChannel:  make(chan *Task),
		worker_num:   cap,
	}

	return &p
}

//协程池创建一个Worker, 并且让这个Worker去工作
func (p *Pool) worker(worker_ID int) {
	//一个worker具体的工作

	//1 永久的从JobsChannel去取任务
	for task := range p.JobsChannel {
		//task 就是当前Worker 从 JobsChannel 中拿到的任务
		//2 一旦取到任务, 执行这个任务 这里可以做优先级的封装
		task.Execute()
		fmt.Println("worker ID", worker_ID, " 执行完了一个任务")
	}
}

//让协程池， 开始真正的工作, 协程池一个启动方法
func (p *Pool) run() {
	//1 根据worker_num 来创建worker去工作
	for i := 0; i < p.worker_num; i++ {
		//每个worker都应该是一个goroutine
		go p.worker(i)
	}

	//2 从EntryChannel中去任务, 将取到的任务, 发送给JobsChannel
	for task := range p.EntryChannel {
		//一旦有task读到
		p.JobsChannel <- task
	}
}

//主函数 来测试协程池的工作
func main() {
	//1 创建一些任务
	t := NewTask(func() error {
		fmt.Println(time.Now())   //我们需要操作的代码
		return nil
	})

	//2 创建一个Pool 协程池
	p := NewPool(4)

	//3 将这些任务 交给协程池Pool
	task_num := 0  //统计任务的数量的初始值
	go func() {
		for {
			//不断的向p中写入任务t, 每个任务就是打印当前的时间
			p.EntryChannel <- t
			task_num += 1       //统计人数的数量
			fmt.Println("当前一共执行了 ", task_num, "个任务")
		}
	}()

	//4 启动Pool， 让Pool开始工作, 此时pool会创建worker， 让worker工作
	p.run()
}
