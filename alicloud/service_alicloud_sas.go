package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type SasService struct {
	client *connectivity.AliyunClient
}

func (s *SasService) DescribeAllGroups(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewSasClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeAllGroups"
	request := map[string]interface{}{}
	idExist := false
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Groups", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Groups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SecurityCenter", id)), NotFoundWithResponse, response)
	}
	for _, vv := range v.([]interface{}) {
		if fmt.Sprint(vv.(map[string]interface{})["GroupId"]) == id {
			idExist = true
			return vv.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("SecurityCenter", id)), NotFoundWithResponse, response)
	}
	return object, nil
}
