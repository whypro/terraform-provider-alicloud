package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeVpcAttribute invokes the vpc.DescribeVpcAttribute API synchronously
// api document: https://help.aliyun.com/api/vpc/describevpcattribute.html
func (client *Client) DescribeVpcAttribute(request *DescribeVpcAttributeRequest) (response *DescribeVpcAttributeResponse, err error) {
	response = CreateDescribeVpcAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeVpcAttributeWithChan invokes the vpc.DescribeVpcAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/describevpcattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeVpcAttributeWithChan(request *DescribeVpcAttributeRequest) (<-chan *DescribeVpcAttributeResponse, <-chan error) {
	responseChan := make(chan *DescribeVpcAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeVpcAttribute(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// DescribeVpcAttributeWithCallback invokes the vpc.DescribeVpcAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/describevpcattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeVpcAttributeWithCallback(request *DescribeVpcAttributeRequest, callback func(response *DescribeVpcAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeVpcAttributeResponse
		var err error
		defer close(result)
		response, err = client.DescribeVpcAttribute(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// DescribeVpcAttributeRequest is the request struct for api DescribeVpcAttribute
type DescribeVpcAttributeRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	VpcId                string           `position:"Query" name:"VpcId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	IsDefault            requests.Boolean `position:"Query" name:"IsDefault"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeVpcAttributeResponse is the response struct for api DescribeVpcAttribute
type DescribeVpcAttributeResponse struct {
	*responses.BaseResponse
	RequestId          string                               `json:"RequestId" xml:"RequestId"`
	VpcId              string                               `json:"VpcId" xml:"VpcId"`
	RegionId           string                               `json:"RegionId" xml:"RegionId"`
	Status             string                               `json:"Status" xml:"Status"`
	VpcName            string                               `json:"VpcName" xml:"VpcName"`
	CreationTime       string                               `json:"CreationTime" xml:"CreationTime"`
	CidrBlock          string                               `json:"CidrBlock" xml:"CidrBlock"`
	Ipv6CidrBlock      string                               `json:"Ipv6CidrBlock" xml:"Ipv6CidrBlock"`
	VRouterId          string                               `json:"VRouterId" xml:"VRouterId"`
	Description        string                               `json:"Description" xml:"Description"`
	IsDefault          bool                                 `json:"IsDefault" xml:"IsDefault"`
	ClassicLinkEnabled bool                                 `json:"ClassicLinkEnabled" xml:"ClassicLinkEnabled"`
	ResourceGroupId    string                               `json:"ResourceGroupId" xml:"ResourceGroupId"`
	VSwitchIds         VSwitchIdsInDescribeVpcAttribute     `json:"VSwitchIds" xml:"VSwitchIds"`
	UserCidrs          UserCidrsInDescribeVpcAttribute      `json:"UserCidrs" xml:"UserCidrs"`
	AssociatedCens     AssociatedCensInDescribeVpcAttribute `json:"AssociatedCens" xml:"AssociatedCens"`
	CloudResources     CloudResourcesInDescribeVpcAttribute `json:"CloudResources" xml:"CloudResources"`
}

// CreateDescribeVpcAttributeRequest creates a request to invoke DescribeVpcAttribute API
func CreateDescribeVpcAttributeRequest() (request *DescribeVpcAttributeRequest) {
	request = &DescribeVpcAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "DescribeVpcAttribute", "vpc", "openAPI")
	return
}

// CreateDescribeVpcAttributeResponse creates a response to parse from DescribeVpcAttribute response
func CreateDescribeVpcAttributeResponse() (response *DescribeVpcAttributeResponse) {
	response = &DescribeVpcAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
