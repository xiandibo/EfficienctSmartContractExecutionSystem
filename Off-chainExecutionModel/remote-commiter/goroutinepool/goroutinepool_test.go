package goroutinepool

////主函数
//func TestGoroutinePool(t *testing.T) {
//	//创建一个Task
//	ta := NewTask(func() error {
//		fmt.Println(time.Now())
//		return nil
//	})
//
//	//创建一个协程池,最大开启3个协程worker
//	p := NewPool(3)
//
//	//开一个协程 不断的向 Pool 输送打印一条时间的task任务
//	go func() {
//		for {
//			p.EntryChannel <- ta
//		}
//	}()
//
//	//启动协程池p
//	p.Run()
//
//}
