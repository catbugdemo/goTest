package cache

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// 实现一个缓存

// 设定存储的内容
type Item struct {
	Object     interface{} // 真正的数据项
	Expiration int64       // 生存时间
}

const (
	// 没有过期时间标志
	NotExpiration time.Duration = -1

	// 默认的过期时间
	DefaultExxpiration time.Duration = 0
)

// 缓存结构体
type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item // 缓存数据项存储在 map 中
	mu                sync.RWMutex    // 读写锁
	gcInterval        time.Duration   // 过期数据项清理周期
	stopGc            chan bool       // 关闭GC回收
}

// 开发顺序
// 1. 如果有结构体，一定要先通过 New 来创建一个结构体
// 2. 根据需求理解来 分解流程同时创建函数
//		(缓存中：get,add,update,delete)
//		(自动控制的Gc回收)
// 3. 结束函数
// 4. 测试验收
// 所有暴露在外的接口都要思考是否要加锁

// New
// 创建缓存结构体
func New(defaultExxpiration, gcInterval time.Duration) *Cache {
	c := &Cache{
		defaultExpiration: defaultExxpiration,
		gcInterval:        gcInterval,
		items:             map[string]Item{},
		stopGc:            make(chan bool),
	}

	// 开始启用过期清理
	go c.goLoop()

	return c
}

// 查看,添加，修改，删除

// Add 添加
// k - key . v - value . d - time
// 加锁
// 如果数据存在，则返回错误
func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	// 加写锁
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.get(k); found {
		return fmt.Errorf("Item %s already exist", k)
	}
	// 设置
	c.set(k, v, d)
	return nil
}

// Set 设置缓存数据项 ，无论存在都覆盖
func (c Cache) Set(k string, v interface{}, d time.Duration) {
	var ex int64
	if d == DefaultExxpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		ex = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{
		Object:     v,
		Expiration: ex,
	}
}

func (c *Cache) set(k string, v interface{}, d time.Duration) {
	var ex int64
	if d == DefaultExxpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		ex = time.Now().Add(d).UnixNano()
	}

	c.items[k] = Item{
		Object:     v,
		Expiration: ex,
	}
}

// Get 获取
// 先通过 key 判断能否获取
// 判断时间是否过期
func (c *Cache) Get(k string) interface{} {
	// 读锁 (保证该数据的一致性)
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[k]
	// 如果未找到, 解开读锁
	if !found {
		return nil
	}
	// 找到了,但是时间过期了
	if item.Expired() {
		return nil
	}
	return item.Object
}

// get 获取 (该方法不需要进行读写锁配置)
// 与 Get 相类似
func (c Cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expired() {
		return nil, false
	}
	return item, true
}

// Expired  判断数据是否过期
// 如果时间未
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

// Replace 替换一个存在的数据
// 判断该数据是否存在
// 存在替换
func (c *Cache) Replace(k string, v interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.get(k); !found {
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	c.set(k, v, d)
	return nil
}

// Delete 删除一个数据
func (c *Cache) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delete(k)
}

func (c *Cache) delete(k string) {
	delete(c.items, k)
}

// gcLoop 过期缓存数据项清理
func (c *Cache) goLoop() {
	// 每个一段时间会向该通道发送当时的时间
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		// ticker.C 只作为一个通道
		case <-ticker.C:
			c.DeleteExpired()
			// 通道关闭了
		case <-c.stopGc:
			ticker.Stop()
			return
		}
	}
}

// 删除过期数据
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			c.delete(k)
		}
	}
}

// Save 将缓存数据写入到 w io.Writer 中
func (c *Cache) Save(w io.Writer) (e error) {
	// 返回一个将编码后数据写入w的*Encoder
	enc := gob.NewEncoder(w)
	// 该错误返回很有趣
	defer func() {
		if x := recover(); x != nil {
			e = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range c.items {
		gob.Register(v.Object)
	}
	e = enc.Encode(&c.items)
	return
}

// 从 io.Reader 中读取数据项
func (c *Cache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]Item{}
	if e := dec.Decode(&items); e != nil {
		return e
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range items {
		ov, found := c.items[k]
		// 如果存储的用户找到并且未过期
		if !found || ov.Expired() {
			c.items[k] = v
		}
	}
	return nil
}

// SaveToFile 保存数据项到文件中
func (c *Cache) SaveToFile(file string) error {
	f, e := os.Create(file)
	defer f.Close()

	if e != nil {
		return e
	}
	if e = c.Save(f); e != nil {
		return e
	}
	return nil
}

// LoadFile 从文件中加载缓存数据项
func (c *Cache) LoadFile(file string) error {
	f, e := os.Open(file)
	defer f.Close()
	if e != nil {
		return e
	}
	if e = c.Load(f); e != nil {
		return e
	}
	return nil
}

// Count 返回数据项的数量
func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// 清空缓存
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]Item{}
}

// StopGc 停止过期缓存清理
func (c *Cache) StopGc() {
	c.stopGc <- true
}
