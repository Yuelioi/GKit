package kv

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ---------------------------
// Config / Options
// ---------------------------

type Option func(*config)

type config struct {
	interval   time.Duration // 背景保存/清理间隔
	pretty     bool          // 是否使用 MarshalIndent
	loadOnInit bool          // 是否在 New 时加载文件 (默认 true)
}

// WithSaveInterval 设置后台保存与清理的间隔 (默认 1m)
func WithSaveInterval(d time.Duration) Option {
	return func(c *config) { c.interval = d }
}

// WithPrettyJSON 设置保存时是否美化 JSON（默认 false）
func WithPrettyJSON(pretty bool) Option {
	return func(c *config) { c.pretty = pretty }
}

// WithLoadOnInit 控制 New 是否加载已有文件（默认 true）
func WithLoadOnInit(b bool) Option {
	return func(c *config) { c.loadOnInit = b }
}

// ---------------------------
// Core types
// ---------------------------

// item 内部存储结构，使用 UnixNano 表示时间（0 表示不过期）
type item[V any] struct {
	Value    V     `json:"value"`
	ExpireAt int64 `json:"expire_at,omitempty"`
}

// KVStore 泛型持久化 KV
type KVStore[V any] struct {
	mu       sync.RWMutex
	data     map[string]item[V]
	filePath string

	// 状态
	dirty bool

	// 背景任务
	saveInterval time.Duration
	stopCh       chan struct{}
	wg           sync.WaitGroup

	// 配置
	pretty bool
}

// NewKVStore 创建 KVStore。
// filePath 为空表示 memory-only 模式（不做磁盘 IO）。
// 默认 saveInterval = 1 minute, pretty = false, loadOnInit = true
func NewKVStore[V any](filePath string, opts ...Option) (*KVStore[V], error) {
	cfg := config{
		interval:   time.Minute,
		pretty:     false,
		loadOnInit: true,
	}
	for _, o := range opts {
		o(&cfg)
	}

	// 确保目录存在（仅当 filePath 非空 且 loadOnInit 或将来保存时才需要）
	if filePath != "" {
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	s := &KVStore[V]{
		data:         make(map[string]item[V]),
		filePath:     filePath,
		saveInterval: cfg.interval,
		stopCh:       make(chan struct{}),
		pretty:       cfg.pretty,
	}

	if cfg.loadOnInit && filePath != "" {
		if err := s.load(); err != nil {
			return nil, err
		}
	}

	// 启动后台循环（仅当间隔 > 0）
	if s.saveInterval > 0 {
		s.wg.Add(1)
		go s.backgroundLoop()
	}

	return s, nil
}

// Close 优雅关闭：停止后台任务并强制保存一次（如果文件可写）
func (s *KVStore[V]) Close() error {
	// 关闭后台
	close(s.stopCh)
	s.wg.Wait()
	// 强制保存（memory-only 模式下无操作）
	return s.Save()
}

// ---------------------------
// 基础操作 API
// ---------------------------

// Set 设置键，不设置过期（覆盖旧值）
func (s *KVStore[V]) Set(key string, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = item[V]{Value: value, ExpireAt: 0}
	s.dirty = true
}

// SetWithTTL 设置键并设置 ttl（零或负值表示不过期）
func (s *KVStore[V]) SetWithTTL(key string, value V, ttl time.Duration) {
	var expireAt int64
	if ttl > 0 {
		expireAt = time.Now().Add(ttl).UnixNano()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = item[V]{Value: value, ExpireAt: expireAt}
	s.dirty = true
}

// Get 获取键（惰性过期：如果过期则视为不存在，但不在此处写磁盘）
// 返回 (zero, false) 当不存在或已过期
func (s *KVStore[V]) Get(key string) (V, bool) {
	s.mu.RLock()
	it, ok := s.data[key]
	s.mu.RUnlock()

	var zero V
	if !ok {
		return zero, false
	}
	if it.ExpireAt > 0 && time.Now().UnixNano() > it.ExpireAt {
		// 已过期，惰性返回不存在
		return zero, false
	}
	return it.Value, true
}

// Delete 删除键（如果存在则标记 dirty）
func (s *KVStore[V]) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		s.dirty = true
	}
}

// Exists 判断键是否存在且未过期
func (s *KVStore[V]) Exists(key string) bool {
	s.mu.RLock()
	it, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return false
	}
	if it.ExpireAt > 0 && time.Now().UnixNano() > it.ExpireAt {
		return false
	}
	return true
}

// TTL 返回键剩余生存时间。如果不存在或不过期，返回 (0, false) 或 (0, true) 分别表示不存在/不过期
func (s *KVStore[V]) TTL(key string) (time.Duration, bool) {
	s.mu.RLock()
	it, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		return 0, false
	}
	if it.ExpireAt == 0 {
		// 存在但不过期
		return 0, true
	}
	now := time.Now().UnixNano()
	if now > it.ExpireAt {
		return 0, false
	}
	return time.Duration(it.ExpireAt - now), true
}

// Keys 返回所有未过期的键（顺序不保证）
func (s *KVStore[V]) Keys() []string {
	now := time.Now().UnixNano()
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k, v := range s.data {
		if v.ExpireAt == 0 || v.ExpireAt > now {
			keys = append(keys, k)
		}
	}
	return keys
}

// Save 导出到磁盘（对外暴露的保存方法）
// 如果 filePath == ""（memory-only），则无操作并返回 nil。
func (s *KVStore[V]) Save() error {
	return s.save()
}

// ---------------------------
// 持久化实现（内部）
// ---------------------------

func (s *KVStore[V]) load() error {
	if s.filePath == "" {
		return nil
	}
	f, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(f) == 0 {
		return nil
	}
	tmp := make(map[string]item[V])
	if err := json.Unmarshal(f, &tmp); err != nil {
		return err
	}
	s.mu.Lock()
	s.data = tmp
	// load 后认为与磁盘一致，dirty = false
	s.dirty = false
	s.mu.Unlock()
	return nil
}

// save 在内部执行实际保存（原子写入）
func (s *KVStore[V]) save() error {
	// memory-only 模式跳过
	if s.filePath == "" {
		return nil
	}

	s.mu.Lock()
	// 如果不脏则不保存
	if !s.dirty {
		s.mu.Unlock()
		return nil
	}

	// 在持锁的情况下 Marshal，保证一致性与避免并发 map 读写 panic
	var data []byte
	var err error
	if s.pretty {
		data, err = json.MarshalIndent(s.data, "", "  ")
	} else {
		data, err = json.Marshal(s.data)
	}
	if err != nil {
		// 保持 dirty 原样（以便下次重试）
		s.mu.Unlock()
		return err
	}

	// reset dirty only after successful marshal to avoid losing changes
	s.dirty = false
	s.mu.Unlock()

	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o644); err != nil {
		// 写失败，恢复 dirty 标记以便下次重试
		s.mu.Lock()
		s.dirty = true
		s.mu.Unlock()
		return err
	}
	// 尝试重命名，覆盖目标文件（原子）
	if err := os.Rename(tmpFile, s.filePath); err != nil {
		// 重命名失败也恢复 dirty 标记
		s.mu.Lock()
		s.dirty = true
		s.mu.Unlock()
		return err
	}
	return nil
}

// ---------------------------
// 后台清理与保存循环
// ---------------------------

func (s *KVStore[V]) backgroundLoop() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.saveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			// 在退出前做一次清理和保存尝试
			s.cleanupLocked(time.Now().UnixNano())
			_ = s.save()
			return
		case <-ticker.C:
			now := time.Now().UnixNano()
			s.cleanupLocked(now)
			_ = s.save()
		}
	}
}

// cleanupLocked 在外部无需加锁的情况下调用（内部会加锁），清理过期项并设置 dirty
func (s *KVStore[V]) cleanupLocked(now int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	changed := false
	for k, v := range s.data {
		if v.ExpireAt > 0 && now > v.ExpireAt {
			delete(s.data, k)
			changed = true
		}
	}
	if changed {
		s.dirty = true
	}
}
