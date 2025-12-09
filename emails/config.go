package emails

type EmailConfig struct {
	Name      string            `json:"name,omitempty"`
	Server    string            `json:"server,omitempty"`
	User      string            `json:"user,omitempty"`
	Password  string            `json:"password,omitempty"`
	Receivers map[string]string `json:"receivers,omitempty"`
}
