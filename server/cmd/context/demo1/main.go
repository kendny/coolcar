package main

import (
	"context"
	"fmt"
	"time"
)

/*
对context的测试
*/

type paramKey struct {
}

func main() {
	c := context.WithValue(context.Background(), paramKey{}, "abc") //旧的context会被垃圾回收机制进行回收
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	mainTask(c)
}

func mainTask(c context.Context) {
	fmt.Printf("main task started with params %q\n", c.Value(paramKey{}))
	//smallTask(context.Background(), "task1", 4*time.Second) // 后台任务，新开的context不受约束
	c1, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()
	// task1用2s，如果没完成就取消， 剩下的时间task2接着进行， 用子任务 【c1是子任务】
	smallTask(c1, "task1", 4*time.Second)
	smallTask(c, "task2", 2*time.Second)
}

func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s started with param: %q\n", name, c.Value(paramKey{}))

	select {
	case <-time.After(d): // 等待6s
		fmt.Printf("%s done\n", name)
	case <-c.Done():
		fmt.Printf("%s canceled\n", name)
	}
}

/***
说明： %!q(<nil>)： 指nil不是%!q能展示出来的
*/
