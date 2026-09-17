package verification

import (
	"fmt"
	"strings"
	"time"
)

const (
	defaultSize       = 6
	defaultCodeExpire = time.Minute * 5
)

func NewConfig() *Config {
	return &Config{
		Mode: modeSingle,
		Default: &VerifyConfig{
			Disable: true,
		},
	}
}

type Config struct {
	Mode     string                   `json:"mode,omitempty"`
	Default  *VerifyConfig            `json:"default,omitempty"`
	Verifies map[string]*VerifyConfig `json:"verifies,omitempty"`
}

func (c *Config) Init() (err error) {
	if c.Mode != modelMulti {
		c.Mode = modeSingle
	}
	if c.Mode == modeSingle {
		if c.Default == nil {
			c.Default = defaultVerifyConfig
		} else {
			c.Default.merge(defaultVerifyConfig)
		}
		c.Default.Init()
		if !c.Default.Disable && c.Default.VerifyFlag == VerifyNone {
			err = fmt.Errorf("default verification is enabled but verifyFlag is not set")
		}
	} else {
		cfgs := make(map[string]*VerifyConfig)
		for key, cfg := range c.Verifies {
			cfgs[strings.ToUpper(key)] = cfg.merge(defaultVerifyConfig).Init()
			if !cfg.Disable && cfg.VerifyFlag == VerifyNone {
				err = fmt.Errorf("%s verification is enabled but verifyFlag is not set", key)
			}
		}
		c.Verifies = cfgs
	}
	return
}

func (c *Config) CanSkip(key, identity string) (ok bool) {
	if c.Mode == modeSingle {
		ok = c.Default.CanSkip(identity)
	} else if cfg, OK := c.Verifies[key]; OK {
		ok = cfg.CanSkip(identity)
	}
	return
}

func (c *Config) Disable(key string) (ok bool) {
	if c.Mode == modeSingle {
		ok = c.Default.Disable
	} else if cfg, OK := c.Verifies[key]; OK {
		ok = cfg.Disable
	} else {
		ok = true
	}
	return
}

func (c *Config) VerifyFlag(key string) (flag byte) {
	if c.Mode == modeSingle {
		flag = c.Default.VerifyFlag
	} else if cfg, OK := c.Verifies[key]; OK {
		flag = cfg.VerifyFlag
	} else {
		flag = VerifyNone
	}
	return
}

func (c *Config) getConfig(key string) *VerifyConfig {
	if c.Mode == modeSingle {
		return c.Default
	}
	return c.Verifies[key]
}

type VerifyConfig struct {
	Disable    bool          `json:"disable,omitempty"`
	Size       int           `json:"size,omitempty"`
	Expire     time.Duration `json:"expire,omitempty"`
	VerifyFlag byte          `json:"verifyFlag,omitempty"`
	Skip       []string      `json:"skip,omitempty"`
	skips      map[string]struct{}
}

func (c *VerifyConfig) Init() *VerifyConfig {
	if len(c.Skip) > 0 {
		c.skips = make(map[string]struct{}, len(c.Skip))
		for _, a := range c.Skip {
			c.skips[a] = struct{}{}
		}
		c.Skip = nil
	}
	c.VerifyFlag &= VerifyAll
	return c
}

func (c *VerifyConfig) CanSkip(identity string) (ok bool) {
	if c.skips != nil {
		_, ok = c.skips[identity]
	}
	return
}

func (c *VerifyConfig) Disabled() bool {
	return c.Disable
}

func (c *VerifyConfig) merge(defaults *VerifyConfig) *VerifyConfig {
	if c.Size <= 0 {
		c.Size = defaults.Size
	}
	if c.Expire <= 0 {
		c.Expire = defaults.Expire
	}
	if c.VerifyFlag < 0 {
		c.VerifyFlag = defaults.VerifyFlag
	}
	return c
}

var defaultVerifyConfig = &VerifyConfig{
	Size:       defaultSize,
	Expire:     defaultCodeExpire,
	VerifyFlag: VerifyAll,
}

const (
	modeSingle = "single"
	modelMulti = "multi"
)

const (
	VerifyNone byte = iota
	VerifyAccount
	VerifyPhone
	VerifyAll
)
