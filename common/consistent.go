package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

/*
	一致性Hash算法
*/

// 声明新打切片类型
type units []uint32

// 返回切片长度
func (x units) Len() int {
	return len(x)
}

// 比较两个数的大小
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

// 切片中两个值交换
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// 当Hash环上没有数据时，提示错误
var errEmpty = errors.New("Hash 环没有数据！")

// 创建结构体保存一致性Hash信息
type Consistent struct {
	// hash环， key 为哈希值，值存放节点信息
	circle       map[uint32]string
	sortedHashes units
	// 虚拟节点个数，用来增加hash的平衡性
	VirtualNode int
	// map 读写锁
	sync.RWMutex
}

// 创建一致性Hash算法结构体，设置默认节点数量
func NewConsistent() *Consistent {
	return &Consistent{
		// 初始化变量
		circle: make(map[uint32]string),
		// 设置虚拟环节点数
		VirtualNode: 20,
	}
}

// 自动生成Key值
func (c *Consistent) generateKey(element string, index int) string {
	// 副本Key 生成逻辑
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		// 声明一个数组长度为64
		var scratch [64]byte
		// 拷贝数据到数组中
		copy(scratch[:], key)
		// 使用IEEE 多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// 更新排序 方便查找
func (c *Consistent) updateSortHashes() {
	hashes := c.sortedHashes[:0]
	// 判断切片容量是否过大，如果过大则重置
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}
	// 添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	// 对所有节点hash值进行排序， 方便之后进行二分查找
	sort.Sort(hashes)
	// 重新赋值
	c.sortedHashes = hashes
}

// 向Hash环添加节点
func (c *Consistent) Add(element string) {
	// 加锁
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

//
func (c *Consistent) add(element string) {
	// 循环虚拟节点，设置副本
	for i := 0; i < c.VirtualNode; i++ {
		// 根据生成的节点添加到Hash环中
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	// 更新排序
	c.updateSortHashes()
}

// 删除一个节点
func (c *Consistent) Remove(element string) {
	// 加锁
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// 移除节点
func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, 1)))
	}
	c.updateSortHashes()
}

// 顺时针查找最近的节点
func (c *Consistent) search(key uint32) int {
	// 查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	// 使用二分查找方法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	// 如果超出范围则设置为0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

// 根据数据标识获取最近服务器节点信息
func (c *Consistent) Get(name string) (string, error) {
	// 添加锁
	c.RLock()
	// 解锁
	defer c.Unlock()
	// 如果为零则返回错误
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	// 计算hash值
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
