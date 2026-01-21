package errorx

import (
	"fmt"
	"sync"
)

var (
	mu    sync.RWMutex
	specs = make(map[int]CodeSpec)
)

func Register(spec CodeSpec) error {
	if spec.Code < 0 {
		return fmt.Errorf("invalid code %d", spec.Code)
	}
	if spec.MessageKey == "" {
		return fmt.Errorf("message key required for code %d", spec.Code)
	}
	if spec.HttpStatus < 100 || spec.HttpStatus >= 600 {
		return fmt.Errorf("invalid http status %d", spec.HttpStatus)
	}
	if spec.Version == "" {
		return fmt.Errorf("version required for code %d", spec.Code)
	}

	mu.Lock()
	defer mu.Unlock()

	if old, ok := specs[spec.Code]; ok {
		if old == spec {
			return nil
		}
		return fmt.Errorf("code %d already registered", spec.Code)
	}

	specs[spec.Code] = spec
	return nil
}

func RegisterMust(spec CodeSpec) {
	if err := Register(spec); err != nil {
		panic(err)
	}
}

func GetSpec(code int) (CodeSpec, bool) {
	mu.RLock()
	defer mu.RUnlock()
	spec, ok := specs[code]
	return spec, ok
}

func GetSpecMust(code int) CodeSpec {
	spec, ok := GetSpec(code)
	if !ok {
		panic(fmt.Sprintf("error code %d not registered", code))
	}
	return spec
}
