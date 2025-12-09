package server_config

func NewConfig(server, configType string) *Config {
	return &Config{Server: server, Type: configType}
}

type Config struct {
	Id          int64  `gorm:"column:id;<-:create" json:"id,omitempty"`
	Server      string `gorm:"column:server" json:"server,omitempty"`
	Type        string `gorm:"column:type" json:"type,omitempty"`
	Content     []byte `gorm:"column:content" json:"content,omitempty"`
	LastContent []byte `gorm:"column:last_content" json:"lastContent,omitempty"`
}

func (s *Config) WithContent(content []byte) *Config {
	s.Content = content
	return s
}

func (s *Config) TableName() string {
	return "server_config"
}
