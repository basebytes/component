package captcha

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/mojocn/base64Captcha"
)

const (
	defaultHeight          = 80
	defaultWidth           = 240
	defaultLength          = 4
	defaultGCLimit         = 1024
	defaultFonts           = "wqy-microhei.ttc"
	defaultSource          = "1234567890qwertyuioplkjhgfdsazxcvbnm"
	defaultShowLineOptions = base64Captcha.OptionShowHollowLine
)

func NewConfig() *Config {
	return &Config{
		Mode: modeSingle,
		Default: &ImgConfig{
			Disable: true,
		},
	}
}

type Config struct {
	Mode    string                `json:"mode,omitempty"`
	Default *ImgConfig            `json:"default,omitempty"`
	Captcha map[string]*ImgConfig `json:"captcha,omitempty"`
}

func (c *Config) Init() (err error) {
	if c.Mode != modelMulti {
		c.Mode = modeSingle
	}
	if c.Mode == modeSingle {
		if c.Default == nil {
			c.Default = defaultConfig
		} else {
			c.Default.merge(defaultConfig)
		}
		c.Default.Init()
		if !c.Default.Disable && c.Default.CheckFlag == CheckNone {
			err = fmt.Errorf("default captcha is enabled but checkFlag is not set")
		}
	} else {
		cfgs := make(map[string]*ImgConfig)
		for key, cfg := range c.Captcha {
			cfgs[strings.ToUpper(key)] = cfg.merge(defaultConfig).Init()
			if !cfg.Disable && cfg.CheckFlag == CheckNone {
				err = fmt.Errorf("%s captcha is enabled but checkFlag is not set", key)
			}
		}
		c.Captcha = cfgs
	}
	return
}

func (c *Config) CanSkip(key, identity string) (ok bool) {
	if c.Mode == modeSingle {
		ok = c.Default.CanSkip(identity)
	} else if cfg, OK := c.Captcha[key]; OK {
		ok = cfg.CanSkip(identity)
	}
	return
}

func (c *Config) Disable(key string) (ok bool) {
	if c.Mode == modeSingle {
		ok = c.Default.Disable
	} else if cfg, OK := c.Captcha[key]; OK {
		ok = cfg.Disable
	} else {
		ok = false
	}
	return
}

func (c *Config) CheckFlag(key string) (flag byte) {
	if c.Mode == modeSingle {
		flag = c.Default.CheckFlag
	} else if cfg, OK := c.Captcha[key]; OK {
		flag = cfg.CheckFlag
	} else {
		flag = CheckAll
	}
	return
}

type ImgConfig struct {
	Disable         bool          `json:"disable,omitempty"`
	Height          int           `json:"height,omitempty"`
	Width           int           `json:"width,omitempty"`
	NoiseCount      int           `json:"noiseCount,omitempty"`
	ShowLineOptions int           `json:"showLineOptions,omitempty"`
	Length          int           `json:"length,omitempty"`
	Source          string        `json:"source,omitempty"`
	BgColor         *color.RGBA   `json:"bgColor,omitempty"`
	Fonts           string        `json:"fonts,omitempty"`
	GCLimit         int           `json:"gcLimit,omitempty"`
	Expiration      time.Duration `json:"expiration,omitempty"`
	CheckFlag       byte          `json:"checkFlag,omitempty"`
	Skip            []string      `json:"skip,omitempty"`
	skips           map[string]struct{}
}

func (c *ImgConfig) Init() *ImgConfig {
	if len(c.Skip) > 0 {
		c.skips = make(map[string]struct{})
		for _, skip := range c.Skip {
			c.skips[skip] = struct{}{}
		}
		c.Skip = nil
	}
	c.CheckFlag &= CheckAll
	return c
}

func (c *ImgConfig) CanSkip(identity string) (ok bool) {
	if c.skips != nil {
		_, ok = c.skips[identity]
	}
	return
}

func (c *ImgConfig) merge(defaults *ImgConfig) *ImgConfig {
	if c.Height <= 0 {
		c.Height = defaults.Height
	}
	if c.Width <= 0 {
		c.Width = defaults.Width
	}
	if c.NoiseCount < 0 {
		c.NoiseCount = defaults.NoiseCount
	}
	if c.ShowLineOptions <= 0 {
		c.ShowLineOptions = defaults.ShowLineOptions
	}
	if c.Length <= 0 {
		c.Length = defaults.Length
	}
	if c.Source == "" {
		c.Source = defaults.Source
	}
	if c.BgColor == nil {
		c.BgColor = defaults.BgColor
	}
	if c.Fonts == "" {
		c.Fonts = defaults.Fonts
	}
	if c.GCLimit <= 0 {
		c.GCLimit = defaults.GCLimit
	}
	if c.Expiration <= 0 {
		c.Expiration = defaults.Expiration
	}
	return c
}

var defaultConfig = &ImgConfig{
	Height:          defaultHeight,
	Width:           defaultWidth,
	ShowLineOptions: defaultShowLineOptions,
	Length:          defaultLength,
	Source:          defaultSource,
	BgColor:         &color.RGBA{R: 255, G: 255, B: 255, A: 255},
	Fonts:           defaultFonts,
	GCLimit:         defaultGCLimit,
	Expiration:      time.Minute,
}

const (
	modeSingle = "single"
	modelMulti = "multi"
)

const (
	CheckNone = iota
	CheckAccount
	CheckPhone
	CheckAll
)
