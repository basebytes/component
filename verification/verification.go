package verification

import (
	"errors"
	"sync"

	"github.com/basebytes/expireMap"
	"github.com/basebytes/tools"
)

var (
	codes  map[string]*expireMap.Map[string]
	config *Config
	lock   sync.RWMutex
)

func Init(config *Config) {
	Reload(config)
}

func Reload(cfg *Config) {
	if cfg == nil {
		cfg = NewConfig()
		_ = cfg.Init()
	}
	_codes := reload(cfg)
	lock.Lock()
	defer lock.Unlock()
	codes = _codes
	config = cfg
}

func reload(cfg *Config) (codes map[string]*expireMap.Map[string]) {
	codes = make(map[string]*expireMap.Map[string])
	if cfg.Mode == modeSingle {
		codes[defaultKey] = expireMap.New[string](8, cfg.Default.Expire)
	} else if cfg.Verifies != nil {
		for key, _cfg := range cfg.Verifies {
			codes[key] = expireMap.New[string](8, _cfg.Expire)
		}
	}
	return
}

func Code(key, identity string) (code string, err error) {
	if _codes, ok := get(key); ok {
		cfg := config.getConfig(key)
		code, _ = _codes.SetWithExpire(identity, tools.RandomNumChar(cfg.Size), cfg.Expire)
	} else {
		err = verifyNotFoundErr
	}
	return
}

func Verify(key, identity, code string) (err error) {
	if _codes, ok := get(key); !ok {
		err = verifyNotFoundErr
	} else if save, OK := _codes.Get(identity); !OK {
		err = codeExpiredErr
	} else if save != code {
		err = verifyFailedErr
	} else {
		_codes.Del(identity)
	}
	return
}

func Exist(key, identity string) bool {
	if _codes, ok := get(key); ok {
		return _codes.Exist(identity)
	}
	return false
}

func Del(key, identity string) {
	if _codes, ok := get(key); ok {
		_codes.Del(identity)
	}
}

func Disable(key string) bool {
	lock.RLock()
	defer lock.RUnlock()
	return config.Disable(key)
}

func CanSkip(key, identity string) bool {
	lock.RLock()
	defer lock.RUnlock()
	return config.CanSkip(key, identity)
}

func VerifyFlag(key string) byte {
	lock.RLock()
	defer lock.RUnlock()
	return config.VerifyFlag(key)
}

func get(key string) (_codes *expireMap.Map[string], ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	if config.Mode == modeSingle {
		_codes, ok = codes[defaultKey]
	} else {
		_codes, ok = codes[key]
	}
	return
}

const defaultKey = "default"

var (
	codeExpiredErr    = errors.New("验证码已过期")
	verifyFailedErr   = errors.New("验证码错误")
	verifyNotFoundErr = errors.New("验证码不存在")
)
