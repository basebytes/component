package sms

type Client interface {
	SendTemplateSms(templateName string, receivers string, params map[string]string, result any) error
	GetTemplate(templateId string, result any) error
	PullMessageStatus(taskId string, result any) error
}

const (
	ClientNameLianLu = "lianlu"
)

type newClient func(*Config) (Client, error)
