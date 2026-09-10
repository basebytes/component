package captcha

import (
	"errors"
	"sync"
)

var (
	captcha map[string]*Captcha
	config  *Config
	lock    sync.RWMutex
)

func Init(cfg *Config) {
	Reload(cfg)
}

func Reload(cfg *Config) {
	if cfg == nil {
		cfg = NewConfig()
		_ = cfg.Init()
	}
	capt := reload(cfg)
	lock.Lock()
	defer lock.Unlock()
	captcha = capt
	config = cfg
}

func Generate(key string) (id string, b64s string, err error) {
	if capt, ok := get(key); ok {
		id, b64s, err = capt.Generate()
	} else {
		err = captchaNotFoundErr
	}
	return
}

func Verify(key, id, value string) (valid bool) {
	if capt, ok := get(key); ok {
		valid = capt.Verify(id, value)
	}
	return
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

func CheckFlag(key string) byte {
	lock.RLock()
	defer lock.RUnlock()
	return config.CheckFlag(key)
}

func get(key string) (capt *Captcha, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	if config.Mode == modeSingle {
		capt, ok = captcha[defaultKey]
	} else {
		capt, ok = captcha[key]
	}
	return
}

func reload(cfg *Config) (capt map[string]*Captcha) {
	capt = make(map[string]*Captcha)
	if cfg.Mode == modeSingle {
		capt[defaultKey] = NewCaptcha(cfg.Default).Init()
	} else if cfg.Captcha != nil {
		for key, _cfg := range cfg.Captcha {
			capt[key] = NewCaptcha(_cfg).Init()
		}
	}
	return
}

const defaultKey = "default"

var captchaNotFoundErr = errors.New("captcha not found")
