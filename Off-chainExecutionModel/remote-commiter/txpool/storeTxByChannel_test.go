package txpool

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewService(t *testing.T) {
	// 1. 新建一个active object, 并增加结束信号
	c := make(chan struct{})
	s := NewService(100, func() { c <- struct{}{} })

	// 2. 起n个goroutine不断执行增加操作
	n := 10000
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(a int) {
			tx := &Tx{
				Tx:            nil,
				AccessList:    nil,
				State:         false,
				ResponseCount: 0,
				LockList:      nil,
			}
			s.Add(a, tx)
			wg.Done()
		}(i)
	}
	wg.Wait()



	s.Close()

	<-c

	mm := s.Slice()
	v := s.Get(20000)
	fmt.Println(v)
	// 3. 校验所有结果是否都被添加上
	fmt.Println("done len:", len(mm))
}