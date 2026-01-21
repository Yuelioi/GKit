package i18n

import "sync"

// Registry 是一个可增量、可组合的翻译器
type Registry struct {
	mu       sync.RWMutex
	messages map[Locale]map[Key]string
}

func NewRegistry() *Registry {
	return &Registry{
		messages: make(map[Locale]map[Key]string),
	}
}

// Register 注册/覆盖一个 key
func (r *Registry) Register(key Key, locale Locale, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.messages[locale]; !ok {
		r.messages[locale] = make(map[Key]string)
	}
	r.messages[locale][key] = message
}

// RegisterBatch 批量注册
func (r *Registry) RegisterBatch(locale Locale, msgs map[Key]string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.messages[locale]; !ok {
		r.messages[locale] = make(map[Key]string)
	}
	for k, v := range msgs {
		r.messages[locale][k] = v
	}
}

// Delete 删除单个 key
func (r *Registry) Delete(key Key, locale Locale) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if m, ok := r.messages[locale]; ok {
		delete(m, key)
	}
}

// DeleteLocale 删除整个 locale
func (r *Registry) DeleteLocale(locale Locale) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.messages, locale)
}

// Translate 实现 Translator 接口
func (r *Registry) Translate(key Key, locale Locale) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if m, ok := r.messages[locale]; ok {
		v, ok := m[key]
		return v, ok
	}
	return "", false
}
