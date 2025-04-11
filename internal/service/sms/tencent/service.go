package tencent

import (
	"context"
	"fmt"

	smsSvc "github.com/l0slakers/webook/internal/service/sms"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	Client   *sms.Client
	appId    *string
	signName *string
}

func NewService(Client *sms.Client, appId, signature string) *Service {
	return &Service{
		Client:   Client,
		appId:    &appId,
		signName: &signature,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, phone ...string) (err error) {
	request := sms.NewSendSmsRequest()

	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = common.StringPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(phone)

	response, err := s.Client.SendSms(request)
	// 处理异常
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return
	}

	// 对于批量发送的情况这里的处理并不好
	// 批量发送的情况要检索出发送失败的短信，而不是统一返回错误
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			err = smsSvc.SendSmsError(*status.Code, *status.Message)
			return
		}
	}

	return nil
}
