package i18n

import "sync"

type Translator interface {
	Translate(key Key, locale Locale) (string, bool)
}

var (
	global Translator
	mu     sync.RWMutex
)

func init() {
	global = NewRegistry()
}

// SetTranslator 设置全局翻译器（通常用于接入第三方 i18n 系统）
func SetTranslator(t Translator) {
	if t == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	global = t
}

// GetTranslator 获取全局翻译器
func GetTranslator() Translator {
	mu.RLock()
	defer mu.RUnlock()
	return global
}

// Translate 翻译入口（统一出口）
func Translate(key Key, locale Locale) string {
	mu.RLock()
	defer mu.RUnlock()

	if msg, ok := global.Translate(key, locale); ok {
		return msg
	}
	return key.String()
}
