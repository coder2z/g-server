/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/6 18:05
 **/
package xsms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/myxy99/component/xcfg"
)

type (
	SmsResponse = dysmsapi.SendSmsResponse
	SmsRequest  = dysmsapi.SendSmsRequest
	client      struct {
		SMS *dysmsapi.Client
	}
)

func (i *smsInvoker) newSMSClient(o *options) *client {
	c, err := dysmsapi.NewClientWithAccessKey(o.Area, o.AccessKeyId, o.AccessSecret)
	if err != nil {
		panic(err)
	}
	return &client{c}
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

func (ali *client) Send(req *SmsRequest) (*SmsResponse, error) {
	if req.RpcRequest == nil {
		req.RpcRequest = new(requests.RpcRequest)
	}
	req.InitWithApiInfo("Dysmsapi", "2017-05-25", "SendSms", "dysms", "openAPI")
	rep, err := ali.SMS.SendSms(req)
	if err != nil {
		return nil, err
	}
	return rep, nil
}
