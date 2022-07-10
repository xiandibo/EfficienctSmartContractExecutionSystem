package txpool

import (
	"github.com/ethereum/go-ethereum/common"
	cmap "github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
	"github.com/vadiminshakov/committer/addrtoip"
	pb "github.com/vadiminshakov/committer/coordinator"
	"github.com/vadiminshakov/committer/queue"
	"sort"
	"strconv"
	"sync"
)

type Tx struct {
	Tx *pb.Transaction
	AccessList []common.Address
	State bool
	ResponseCount uint64
	LockList SortAddress
}

type AddrAndResponse struct {
	Addr common.Address
	response bool
}

type Txpool struct {
	UnProcessTxQueue *queue.Queue //待处理的交易
	//ProcessingQueue *queue.Queue //正在进行2PC的交易
	//ProcessingMap map[uint64]*Tx //等待被启动2PC阶段的Tx
	//DoneProcess []*Tx
	TxIn2PCMap *Service //正在2PC阶段的Tx,Service是一个更高性能的支持并发写的map
	TxIn2PCCMap cmap.ConcurrentMap

	mu sync.Mutex
	RWmu sync.RWMutex


	// 合约地址与mutex的映射,主要想借助Mutex的阻塞唤醒机制
	AddrMutexMap map[common.Address]*sync.Mutex
	muForAddrMutexMap sync.RWMutex


}

func NewTxpool() *Txpool {
	size := 10000 // 通道缓冲容量
	c := make(chan struct{})
	s := NewService(size, func() { c <- struct{}{} }) // 新建一个支持并发读写的map

	m := cmap.New()

	txpool := &Txpool{
		UnProcessTxQueue: queue.New(),
		//ProcessingMap: make(map[uint64]*Tx),
		//DoneProcess: make([]*Tx, 0),
		AddrMutexMap: make(map[common.Address]*sync.Mutex),
		TxIn2PCMap: s,
		TxIn2PCCMap: m,
	}
	return txpool
}

// InitAddrMutexMap 为了保证AddrMutexMap对应预初始化的合约都是存在的
func (t *Txpool) InitAddrMutexMap() {
	for addr, _ := range addrtoip.AddressToIP {
		t.AddrMutexMap[addr] = new(sync.Mutex)
	}
}

func (t *Txpool) AddUnProcessingTx(tx *Tx) bool {
	//t.RWmu.Lock()
	//defer t.RWmu.Unlock()
	//
	//if _, ok := t.ProcessingMap[tx.Tx.Id]; !ok {
	//	t.ProcessingMap[tx.Tx.Id] = tx
	//	//log.Info("New Transaction, AddUnProcessingTx.")
	//	return true
	//} else {
	//	log.Info("Transaction already exist! Drop. Transaction ID: ", tx.Tx.Id)
	//	return false
	//}

	//2021/8/5 新的实现方案，使用无锁的并发写map方案
	//id := int(tx.Tx.Id)
	//v := t.TxIn2PCMap.Get(id)
	//if v == nil {
	//	t.TxIn2PCMap.Add(id, tx)
	//	return true
	//} else {
	//	log.Info("Transaction already exist! Drop. Transaction ID: ", id)
	//	return false
	//}

	id := strconv.Itoa(int(tx.Tx.Id))
	if _, ok := t.TxIn2PCCMap.Get(id); !ok {
		t.TxIn2PCCMap.Set(id, tx)
		return true
	} else {
		log.Info("Transaction already exist! Drop. Transaction ID: ", id)
		return false
	}


}

func (t *Txpool) GetUnProcessingTx(id uint64) *Tx {

	//t.RWmu.RLock()
	////defer log.Info("GetUnProcessingTx t.RWmu.RUnlock()")
	//defer t.RWmu.RUnlock()
	//
	//if _, ok := t.ProcessingMap[id]; !ok {
	//	return nil
	//} else {
	//
	//	return t.ProcessingMap[id]
	//}

	//2021/8/5 新的实现方案，使用无锁的并发写map方案

	//return t.TxIn2PCMap.Get(int(id))

	temp := strconv.Itoa(int(id))
	tx, _ := t.TxIn2PCCMap.Get(temp)
	return tx.(*Tx)

}


// IsSerializable 检查tx的accessList涉及到的合约地址有没有被其他已经在进行中的事务锁定
// 即尝试申请锁
// 如何高效地进行锁的申请与管理？
// TODO：事务完成以后必须释放锁
func (t *Txpool) IsSerializable(tx *Tx) bool {
	AddrListForLock := SortAddress{}
	for i, addr := range tx.AccessList {
		// 第一个addr属于main SC ，不需要按顺序加锁，直接加锁
		if i == 0 {
			t.AddrMutexMap[addr].Lock()
			continue
		}
		if find(AddrListForLock, addr) == false {
			// append
			AddrListForLock = append(AddrListForLock, addr)
		}
	}
	// 对AddrListForLock排序
	sort.Sort(AddrListForLock)

	// 按照顺序申请锁，防止死锁
	for _, addr := range AddrListForLock {
		// t.AddrMutexMap 初始是空的，申请前先检查map中是否存在对应的合约
		// TODO：这里不加读写锁，将来会不会存在对AddrMutexMap的读写冲突产生panic？
		//log.Info("正在申请到所需的锁， txID====", tx.Tx.Id, "LockIndex====", i, addr)
		t.AddrMutexMap[addr].Lock()
		//log.Info("成功申请到所需的锁， txID====", tx.Tx.Id, "LockIndex====", i, addr)

		//t.muForAddrMutexMap.RLock()
		//if _, exist := t.AddrMutexMap[addr]; !exist {
		//	t.muForAddrMutexMap.RUnlock()
		//	t.muForAddrMutexMap.Lock()
		//	if _, ok := t.AddrMutexMap[addr]; ok { // 再次确认
		//		t.AddrMutexMap[addr].Lock()
		//		t.muForAddrMutexMap.Unlock()
		//		continue
		//	}
		//
		//	// 不存在合约
		//	t.AddrMutexMap[addr] = new(sync.Mutex)
		//	t.AddrMutexMap[addr].Lock()
		//	t.muForAddrMutexMap.Unlock()
		//	continue
		//}
		//
		//// 申请锁
		//// 遵循谁申请，谁释放的原则
		//t.AddrMutexMap[addr].Lock()
		//t.muForAddrMutexMap.RUnlock()
	}

	// 把AddrListForLock 与Tx绑定，方便事务完成以后对相应的合约解锁
	tx.LockList = SortAddress{}
	tx.LockList = append(tx.LockList, tx.AccessList[0])
	tx.LockList = append(tx.LockList, AddrListForLock...)
	// tx.LockList = AddrListForLock

	return true
}

// 检查list中是否存在addr
func find(addrList []common.Address, addr common.Address) bool {
	for _, elem := range addrList {
		if elem == addr {
			return true
		}
	}
	return false
}


// 实现sort接口
type SortAddress []common.Address

func (s SortAddress) Len() int {
	return len(s)
}

func (s SortAddress) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortAddress) Less(i, j int) bool {
	for c := 0; c < 20; c++ {
		if s[i][c] < s[j][c] {
			return true
		} else if s[i][c] == s[j][c] {
			continue
		} else {
			return false
		}
	}
	return false
}