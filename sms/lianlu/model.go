package lianlu

import "fmt"

type Request interface {
	Uri() string
	SignParam() string
	SetSignature(string)
}

type TemplateSmsSendRequest struct {
	*baseParam
	TemplateId       string   `json:"TemplateId"`
	Type             string   `json:"Type"`
	Tag              string   `json:"Tag,omitempty"`
	TaskTime         string   `json:"TaskTime,omitempty"`
	PhoneNumberSet   []string `json:"PhoneNumberSet"`
	TemplateParamSet []string `json:"TemplateParamSet"`
}

func (r *TemplateSmsSendRequest) Uri() string {
	return "sms/trade/template/send"
}

func (r *TemplateSmsSendRequest) SignParam() string {
	optional := ""
	if r.Tag != "" {
		optional += fmt.Sprintf("&Tag=%s", r.Tag)
	}
	if r.TaskTime != "" {
		optional += fmt.Sprintf("&TaskTime=%s", r.TaskTime)
	}
	return fmt.Sprintf("AppId=%s&MchId=%s&SignType=%s%s&TemplateId=%s&TimeStamp=%s&Type=%s&Version=%s&key=%s", r.AppId, r.MchId, r.SignType,
		optional, r.TemplateId, r.TimeStamp, r.Type, r.Version, r.Key)
}

func NewTemplateSmsSendResponse() *TemplateSmsSendResponse {
	return &TemplateSmsSendResponse{}
}

type TemplateSmsSendResponse struct {
	TaskId    string `json:"taskId,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Message   string `json:"message,omitempty"`
}

type TemplateQueryRequest struct {
	*baseParam
	TemplateId string `json:"TemplateId"`
}

func (r *TemplateQueryRequest) Uri() string {
	return "sms/product/template/getById"
}

func (r *TemplateQueryRequest) SignParam() string {
	return fmt.Sprintf("AppId=%s&MchId=%s&SignType=%s&TemplateId=%s&TimeStamp=%s&Version=%s&key=%s", r.AppId, r.MchId, r.SignType,
		r.TemplateId, r.TimeStamp, r.Version, r.Key)
}

type MessageStatusPullRequest struct {
	*baseParam
	TaskId string `json:"TaskId"`
}

func (r *MessageStatusPullRequest) Uri() string {
	return "sms/trade/report"
}

func (r *MessageStatusPullRequest) SignParam() string {
	return fmt.Sprintf("AppId=%s&MchId=%s&SignType=%s&TaskId=%s&TimeStamp=%s&Version=%s&key=%s", r.AppId, r.MchId, r.SignType,
		r.TaskId, r.TimeStamp, r.Version, r.Key)
}

type baseParam struct {
	AppId     string `json:"AppId"`
	MchId     string `json:"MchId"`
	Version   string `json:"Version"`
	TimeStamp string `json:"TimeStamp"`
	SignType  string `json:"SignType"`
	Signature string `json:"Signature"`
	Key       string `json:"-"`
}

func (b *baseParam) SetSignature(signature string) {
	b.Signature = signature
}
