package re

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 复习 go-cache
// 今天，可以来学习写下 go-cache
// 本项目代码全部内容可在 github.com/patrickmn/go-cache 中查询

// 需要的前置知识
//	1.什么是读锁，什么是写锁， sync.Mutex 和 sync.RWMutex 区别
//	2.go runtime 协程参考
//	3.

// 写 go-cahce 的顺序
// 	1.确定数据需要存储的内容，和其过期时间
//	2.将cache封装为一个结构体
//	3.New
//	4.主方法包括:
//		添加，修改，删除，查询，GC垃圾回收
//	5.测试
//	6.提交

// 1. 确定数据需要存储的内容:
//		数据，过期时间

type Item struct {
	// 能够存储的任何种类的数据
	Object interface{}
	// 源码中过期时间,为什么是 int64 (因为 time.Duration 的代指就是 int64 )
	Expiration int64
}

// 判断是否过期
func (item Item) Expired() bool {
	// 如果过期时间0,则表示未设置过期时间
	if item.Expiration == 0 {
		return false
	}
	// 当前时间大于过期时间 则未过期
	return time.Now().UnixNano() > item.Expiration
}

const (
	// 设置永久过期时间为 -1
	NoExpiration time.Duration = -1
	// 设置默认过期时间为 0
	DefaultExpiration time.Duration = 0
)

// 该结构体方便于外部调用
type Cache struct {
	*cache
}

// 2.cache 中存储的内容:
//		默认过期时间，存储的 item, 读写锁, janitor
type cache struct {
	defaultExpiration time.Duration             // 默认过期时间
	items             map[string]Item           // 内容
	mu                sync.RWMutex              // 读写锁
	onEvicted         func(string, interface{}) //
	janitor           *janitor                  // gc 垃圾回收
}

// 4.主方法

// Set (添加) 该方法因为在外部层面需要保证数据的一致性，防止多个操作
// 该方法伟伦是否有值，都会重新设置
// k - key , x - value , d - 过期时间
func (c cache) Set(k string, x interface{}, d time.Duration) {
	// 最后需要存入的数据内容
	var e int64
	// 如果为 0，则设置为之前设置的默认过期时间
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	// 设置过期时间从当前时间开始
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	// 添加写锁，保证数据的原子性,并写入内容
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[k] = Item{
		Object:     x, //存储内容
		Expiration: e, //过期时间
	}
}

// set 与 Set 相同，因为作用在 Add 中 已经加了锁，所以不需要加锁
func (c *cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

// SetDefault 设置之前的时间为过期时间
func (c *cache) SetDefault(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

// Add 如果数据存在则不会进行操作
func (c cache) Add(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.get(k); found { // 已经存储了数据
		return fmt.Errorf("Item %s already exists", k)
	}
	c.set(k, x, d)
	return nil
}

// Replace 修改
func (c *cache) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.get(k); !found {
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	// 否则进行设置
	c.set(k, x, d)
	return nil
}

// Get 获取 该数据为
func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.Unlock()
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}

// GetWithExpiration 该方法同时返回数据的过期时间
func (c cache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[k]
	if !found {
		return nil, time.Time{}, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, time.Time{}, false
		}
		//
		return item.Object, time.Unix(0, item.Expiration), true
	}
	// 返回过期时间未永久的
	return item.Object, time.Time{}, true
}

// Delete 删除缓存, evicted 用来存储驱逐数据
func (c *cache) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, evicted := c.delete(k)
	if evicted {
		c.onEvicted(k, v)
	}
}

// delete 删除缓存
func (c *cache) delete(k string) (interface{}, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

// get 不需要读写锁的原因同 set
func (c *cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 { //设置了过期时间，但是过期了
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}


type keyAndValue struct {
	key   string
	value interface{}
}

func (c *cache) DeleteExpired() {
	var evictedItems []keyAndValue
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			ov, eveicted := c.delete(k)
			if eveicted {
				evictedItems = append(evictedItems, keyAndValue{k, ov})
			}
		}
	}
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvicted = f
}

// Items 获取所有数据
func (c *cache) Items() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := make(map[string]Item, len(c.items))
	now := time.Now().UnixNano()
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			continue
		}
		m[k] =v
	}
	return m
}

func (c cache) ItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Flush 删除所有
func (c cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = map[string]Item{}
}

// janitor gc的垃圾回收
type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache)  {
	// 每过一段时间，返回时间通道
	ticker := time.NewTicker(j.Interval)
	for  {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}

// 3. New 一个 cache
func newCache(de time.Duration, m map[string]Item) *cache {
	// 如果过期时间是0，创建永久时间
	if de == 0 {
		de = -1
	}
	return &cache{
		defaultExpiration: de,
		items:             m,
	}
}


func newCacheWithJanitor(de, ci time.Duration, m map[string]Item) *Cache {
	c := newCache(de, m)

	C := &Cache{c}

	if ci > 0 {
		// 运行垃圾回收
		runJanitor(c,ci)
		// 设置垃圾回收
		runtime.SetFinalizer(C,stopJanitor)
	}
	return C
}

func New(defaultExpiration, cleanInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithJanitor(defaultExpiration, cleanInterval, items)
}
