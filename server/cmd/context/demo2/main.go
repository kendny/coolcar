package main

import (
	"context"
	"fmt"
	"time"
)

type paramKey struct {
}

func main() {
	c := context.WithValue(context.Background(), paramKey{}, "abc")

	c, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	go mainTask(c)

	var cmd string
	for {
		fmt.Scan(&cmd)
		if cmd == "c" {
			cancel() // 取消主任务
		}
	}
}

func mainTask(c context.Context) {
	fmt.Printf("main task started with param %q\n", c.Value(paramKey{}))
	// 开启协程
	go func() {
		c1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		smallTask(c1, "background_task", 8*time.Second) // 后台任务， 独立于主任务
	}()

	go func() {
		c1, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()
		smallTask(c1, "sub_task", 9*time.Second) // 主任务的子任务
	}()
	smallTask(c, "same_task", 8*time.Second) // 主任务的同一任务的步骤
}

func smallTask(c context.Context, name string, d time.Duration) {
	fmt.Printf("%s started with param %q\n", name, c.Value(paramKey{}))
	select {
	case <-time.After(d):
		fmt.Printf("%s done\n", name)
	case <-c.Done():
		fmt.Printf("%s canceled\n", name)
	}
}
