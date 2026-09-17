package sms

type Config struct {
	Name      string                     `json:"name,omitempty"`
	Id        string                     `json:"id,omitempty"`
	Secret    string                     `json:"secret,omitempty"`
	MchId     string                     `json:"mchId,omitempty"`
	Templates map[string]*TemplateConfig `json:"templates,omitempty"`
}

func (c *Config) GetTemplateConfig(name string) *TemplateConfig {
	if c.Templates != nil {
		return c.Templates[name]
	}
	return nil
}

type TemplateConfig struct {
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}
