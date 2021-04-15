package xsms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xlog"
)

type (
	SmsResponse = dysmsapi.SendSmsResponse
	SmsRequest  = dysmsapi.SendSmsRequest
	Client      struct {
		SMS          *dysmsapi.Client
		signName     string
		templateCode string
	}
)

func (i *smsInvoker) newSMSClient(o *options) *Client {
	c, err := dysmsapi.NewClientWithAccessKey(o.Area, o.AccessKeyId, o.AccessSecret)
	if err != nil {
		xlog.Panic("Application Starting",
			xlog.FieldComponentName("XInvoker"),
			xlog.FieldMethod("XInvoker.XSms.NewSMSClient"),
			xlog.FieldDescription("New SMSClient error"),
			xlog.FieldErr(err),
		)
	}
	return &Client{SMS: c, signName: o.SignName, templateCode: o.TemplateCode}
}

func (i *smsInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)
	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newSMSOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}

func (ali *Client) Send(req *SmsRequest) (*SmsResponse, error) {
	if req.RpcRequest == nil {
		req.RpcRequest = new(requests.RpcRequest)
	}
	if req.TemplateCode == "" {
		req.TemplateCode = ali.templateCode
	}
	if req.SignName == "" {
		req.SignName = ali.signName
	}
	req.InitWithApiInfo("coder2z_sms_api", "2017-05-25", "SendSms", "coder2z_sms", "openAPI")
	rep, err := ali.SMS.SendSms(req)
	if err != nil {
		return nil, err
	}
	return rep, nil
}
