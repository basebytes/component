package captcha

import (
	"strings"

	"github.com/mojocn/base64Captcha"
)

func NewCaptcha(cfg *ImgConfig) *Captcha {
	return &Captcha{
		driver: (&base64Captcha.DriverString{
			Height:          cfg.Height,
			Width:           cfg.Width,
			NoiseCount:      cfg.NoiseCount,
			ShowLineOptions: cfg.ShowLineOptions,
			Length:          cfg.Length,
			Source:          cfg.Source,
			BgColor:         cfg.BgColor,
			Fonts:           strings.Split(cfg.Fonts, ","),
		}).ConvertFonts(),
		store: base64Captcha.NewMemoryStore(cfg.GCLimit, cfg.Expiration),
	}
}

type Captcha struct {
	driver  base64Captcha.Driver
	store   base64Captcha.Store
	captcha *base64Captcha.Captcha
}

func (c *Captcha) Init() *Captcha {
	c.captcha = base64Captcha.NewCaptcha(c.driver, c.store)
	return c
}

func (c *Captcha) Generate() (id string, b64s string, err error) {
	var answer string
	if id, b64s, answer, err = c.captcha.Generate(); err == nil {
		err = c.store.Set(id, answer)
		//data, _ := base64.StdEncoding.DecodeString(strings.Split(b64s, ",")[1])
		//_ = os.WriteFile("captcha.png", data, 0644)
		//log.Println("Captcha:", answer)
	}
	return
}

func (c *Captcha) Verify(id, value string) bool {
	return c.store.Verify(id, strings.ToLower(value), true)
}
