package lianlu

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/basebytes/component/sms"
	"github.com/basebytes/tools"
)

const defaultEndPoint = "apis.shlianlu.com"

const (
	defaultApiVersion         = "1.1.0"
	defaultProtocol           = "https"
	defaultSignatureAlgorithm = tools.SignatureAlgorithmMD5
)

func NewClient(config *sms.Config) (sms.Client, error) {
	c := &Client{
		host:       defaultEndPoint,
		protocol:   defaultProtocol,
		apiVersion: defaultApiVersion,
		signType:   defaultSignatureAlgorithm,
	}
	return c, c.Init(config)
}

type Client struct {
	host       string
	protocol   string
	apiVersion string
	appId      string
	mchId      string
	key        string
	signType   string
	reqUrl     string
	header     *http.Header
	templates  map[string]*sms.TemplateConfig
}

func (c *Client) Init(config *sms.Config) (err error) {
	if config.Id == "" || config.MchId == "" || config.Secret == "" {
		err = fmt.Errorf("missing required config")
	}
	if err == nil {
		c.appId = config.Id
		c.mchId = config.MchId
		c.key = config.Secret
		c.reqUrl = fmt.Sprintf("%s://%s", c.protocol, c.host)
		c.templates = config.Templates
		c.header = header
	}
	return
}

func (c *Client) SendTemplateSms(templateName, receivers string, params map[string]string, result any) (err error) {
	if receivers == "" {
		err = fmt.Errorf("receiver not found")
	} else if temps, ok := c.templates[templateName]; !ok || temps.Code == "" {
		err = fmt.Errorf("template %s not found", templateName)
	} else {
		err = c.doPost(c.newTemplateSmsSendRequest(templateName, receivers, params), &result)
	}
	return
}

func (c *Client) GetTemplate(templateId string, result any) (err error) {
	if templateId == "" {
		err = fmt.Errorf("templateId not found")
	} else {
		err = c.doPost(c.newTemplateQueryRequest(templateId), &result)
	}
	return
}

func (c *Client) PullMessageStatus(taskId string, result any) (err error) {
	if taskId == "" {
		err = fmt.Errorf("taskId not found")
	} else {
		err = c.doPost(c.newMessageStatusPullRequest(taskId), &result)
	}
	return
}

func (c *Client) doPost(req Request, result any) error {
	c.sign(req)
	response, err := tools.DoRequest(http.MethodPost, c.reqUrl, req.Uri(), nil, c.header, bytes.NewBuffer(tools.EncodeBytes(req)))
	if err == nil {
		var code int
		if code, err = tools.ParseResponse(response, &result); err == nil && code != http.StatusOK {
			err = fmt.Errorf("sms request [%s] failed [%d]", req.Uri(), code)
		}
	}
	return err
}

func (c *Client) sign(r Request) {
	r.SetSignature(strings.ToUpper(hex.EncodeToString(tools.Sign([]byte(r.SignParam()), c.signType))))
}

func (c *Client) newTemplateSmsSendRequest(templateName, receivers string, params map[string]string) *TemplateSmsSendRequest {
	return &TemplateSmsSendRequest{
		baseParam:        c.newBaseParam(),
		TemplateId:       c.templates[templateName].Code,
		Type:             typeTemplateSmsSend,
		PhoneNumberSet:   strings.Split(receivers, ","),
		TemplateParamSet: strings.Split(params[paramParams], ","),
		Tag:              params[paramTag],
		TaskTime:         params[paramTaskTime],
	}
}

func (c *Client) newTemplateQueryRequest(templateId string) *TemplateQueryRequest {
	return &TemplateQueryRequest{
		baseParam:  c.newBaseParam(),
		TemplateId: templateId,
	}
}

func (c *Client) newMessageStatusPullRequest(taskId string) *MessageStatusPullRequest {
	return &MessageStatusPullRequest{
		baseParam: c.newBaseParam(),
		TaskId:    taskId,
	}
}

func (c *Client) newBaseParam() *baseParam {
	return &baseParam{
		AppId:     c.appId,
		MchId:     c.mchId,
		SignType:  c.signType,
		TimeStamp: fmt.Sprintf("%d", tools.GetCurrentMillSeconds()),
		Version:   c.apiVersion,
		Key:       c.key,
	}
}

const (
	typeTemplateSmsSend = "3"
	paramParams         = "params"
	paramTag            = "tag"
	paramTaskTime       = "taskTime"
)

const (
	ErrorCodeSuccess      = "00"
	ErrorCodeInvalidParam = "01"
	ErrorCodeSignErr      = "02"
	ErrorCodeWrongMchId   = "03"
)

var header = &http.Header{
	"Accept":       []string{"application/json"},
	"Content-Type": []string{"application/json;charset=utf-8"},
}

func init() {
	sms.RegisterNewClientFunc(sms.ClientNameLianLu, NewClient)
}
