package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCloudSsoDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudSsoDirectoryCreate,
		Read:   resourceAlicloudCloudSsoDirectoryRead,
		Update: resourceAlicloudCloudSsoDirectoryUpdate,
		Delete: resourceAlicloudCloudSsoDirectoryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"directory_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mfa_authentication_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
			"saml_identity_provider_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encoded_metadata_document": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							Computed:  true,
						},
						"sso_status": {
							Computed:     true,
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
						},
					},
				},
				ForceNew: true,
			},
			"scim_synchronization_status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
		},
	}
}

func resourceAlicloudCloudSsoDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDirectory"
	request := make(map[string]interface{})
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("directory_name"); ok {
		request["DirectoryName"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_directory", action, AlibabaCloudSdkGoERROR)
	}
	responseDirectory := response["Directory"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseDirectory["DirectoryId"]))

	return resourceAlicloudCloudSsoDirectoryUpdate(d, meta)
}
func resourceAlicloudCloudSsoDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}
	getDirectoryObject, err := cloudssoService.DescribeCloudSsoDirectory(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_directory cloudssoService.DescribeCloudSsoDirectory Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("directory_name", getDirectoryObject["DirectoryName"])
	d.Set("mfa_authentication_status", getDirectoryObject["MFAAuthenticationStatus"])
	d.Set("scim_synchronization_status", getDirectoryObject["SCIMSynchronizationStatus"])

	if SAMLIdentityProviderConfiguration, ok := getDirectoryObject["SAMLIdentityProviderConfiguration"]; ok && len(SAMLIdentityProviderConfiguration.(map[string]interface{})) > 0 {
		SAMLIdentityProviderConfigurationSli := make([]map[string]interface{}, 0)
		SAMLIdentityProviderConfigurationMap := make(map[string]interface{})
		SAMLIdentityProviderConfigurationMap["sso_status"] = SAMLIdentityProviderConfiguration.(map[string]interface{})["SSOStatus"]
		if v, ok := SAMLIdentityProviderConfiguration.(map[string]interface{})["EncodedMetadataDocument"]; ok && SAMLIdentityProviderConfiguration.(map[string]interface{})["SSOStatus"].(string) == "Enabled" {
			SAMLIdentityProviderConfigurationMap["encoded_metadata_document"] = v
		}
		SAMLIdentityProviderConfigurationSli = append(SAMLIdentityProviderConfigurationSli, SAMLIdentityProviderConfigurationMap)
		d.Set("saml_identity_provider_configuration", SAMLIdentityProviderConfigurationSli)
	}

	return nil
}
func resourceAlicloudCloudSsoDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DirectoryId": d.Id(),
	}
	if d.HasChange("mfa_authentication_status") {
		update = true
		if v, ok := d.GetOk("mfa_authentication_status"); ok {
			request["MFAAuthenticationStatus"] = v
		}
	}
	if update {
		action := "SetMFAAuthenticationStatus"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("mfa_authentication_status")
	}
	update = false
	setSCIMSynchronizationStatusReq := map[string]interface{}{
		"DirectoryId": d.Id(),
	}
	if d.HasChange("scim_synchronization_status") {
		update = true
		if v, ok := d.GetOk("scim_synchronization_status"); ok {
			setSCIMSynchronizationStatusReq["SCIMSynchronizationStatus"] = v
		}
	}
	if update {
		action := "SetSCIMSynchronizationStatus"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, setSCIMSynchronizationStatusReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setSCIMSynchronizationStatusReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("scim_synchronization_status")
	}
	update = false
	updateDirectoryReq := map[string]interface{}{
		"DirectoryId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("directory_name") {
		update = true
		if v, ok := d.GetOk("directory_name"); ok {
			updateDirectoryReq["NewDirectoryName"] = v
		}
	}
	if update {
		action := "UpdateDirectory"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, updateDirectoryReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateDirectoryReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("directory_name")
	}
	update = false
	setExternalSAMLIdentityProviderReq := map[string]interface{}{
		"DirectoryId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("saml_identity_provider_configuration") {
		update = true
		if v, ok := d.GetOk("saml_identity_provider_configuration"); ok {
			for _, setExternalSAMLIdentityProvider := range v.(*schema.Set).List() {
				setExternalSAMLIdentityProviderArg := setExternalSAMLIdentityProvider.(map[string]interface{})
				setExternalSAMLIdentityProviderReq["SSOStatus"] = setExternalSAMLIdentityProviderArg["sso_status"]
				if v, ok := setExternalSAMLIdentityProviderArg["encoded_metadata_document"]; ok && v != "" {
					setExternalSAMLIdentityProviderReq["EncodedMetadataDocument"] = setExternalSAMLIdentityProviderArg["encoded_metadata_document"]
				}
			}
		}
	}
	if update {
		action := "SetExternalSAMLIdentityProvider"
		conn, err := client.NewCloudssoClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, setExternalSAMLIdentityProviderReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setExternalSAMLIdentityProviderReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	d.Partial(false)
	return resourceAlicloudCloudSsoDirectoryRead(d, meta)
}
func resourceAlicloudCloudSsoDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDirectory"
	var response map[string]interface{}
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DirectoryId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Directory"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
