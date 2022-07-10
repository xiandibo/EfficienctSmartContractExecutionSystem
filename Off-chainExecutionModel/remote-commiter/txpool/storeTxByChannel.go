package txpool

// Active Object方式：只有1个goroutine在执行写操作。避免多个goroutine竞争锁。
// 适合业务场景复杂，性能要求高的场景。

// Service active object对象
type Service struct {
	channel chan chanInput `desc:"即将加入到数据slice的数据"`
	data    map[int]*Tx    `desc:"数据slice"`
}

type chanInput struct {
	key int
	value *Tx
}

// 新建一个size大小缓存的active object对象
func NewService(size int, done func()) *Service {
	s := &Service{
		channel: make(chan chanInput, size),
		data:    make(map[int]*Tx, 0),
	}

	go func() {
		s.schedule()
		done()
	}()
	return s
}

// 把管道中的数据append到slice中
func (s *Service) schedule() {
	for v := range s.channel {
		s.data[v.key] = v.value
		//s.data = append(s.data, v)
	}
}

// 增加一个值
func (s *Service) Add(key int, value *Tx) {
	elem := chanInput{
		key:   key,
		value: value,
	}
	s.channel <- elem
}

// 读取一个值
func (s *Service) Get(key int) *Tx {
	return s.data[key]
}

// 管道使用完关闭
func (s *Service) Close() {
	close(s.channel)
}

// 返回slice
func (s *Service) Slice() map[int]*Tx {
	return s.data
}
