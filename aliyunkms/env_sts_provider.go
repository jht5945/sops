package aliyunkms

import (
	"github.com/aliyun/credentials-go/credentials"
	"os"

	"github.com/alibabacloud-go/tea/tea"
)

type envStsProvider struct{}

var providerEnvSts = new(envStsProvider)

const (
	// EnvVarAccessKeyId is a name of ALIBABA_CLOUD_ACCESS_KEY_ID
	EnvVarAccessKeyId = "ALIBABA_CLOUD_ACCESS_KEY_ID"
	// EnvVarAccessKeySecret is a name of ALIBABA_CLOUD_ACCESS_KEY_SECRET
	EnvVarAccessKeySecret = "ALIBABA_CLOUD_ACCESS_KEY_SECRET"
	// EnvVarSecurityToken is a name of ALIBABA_CLOUD_SECURITY_TOKEN
	EnvVarSecurityToken = "ALIBABA_CLOUD_SECURITY_TOKEN"

	// EnvVarAliCloudAccessKeyId is a name of ALICLOUD_ACCESS_KEY
	EnvVarAliCloudAccessKeyId = "ALICLOUD_ACCESS_KEY"
	// EnvVarAliCloudAccessKeySecret is a name of ALICLOUD_SECRET_KEY
	EnvVarAliCloudAccessKeySecret = "ALICLOUD_SECRET_KEY"
	// EnvVarAliCloudSecurityToken is a name of ALICLOUD_SECURITY_TOKEN
	EnvVarAliCloudSecurityToken = "ALICLOUD_SECURITY_TOKEN"
)

func (p *envStsProvider) resolve() (*credentials.Config, error) {
	config1 := internalResolve(EnvVarAliCloudAccessKeyId, EnvVarAliCloudAccessKeySecret, EnvVarAliCloudSecurityToken)
	if config1 != nil {
		return config1, nil
	}
	config2 := internalResolve(EnvVarAccessKeyId, EnvVarAccessKeySecret, EnvVarSecurityToken)
	if config2 != nil {
		return config2, nil
	}
	return nil, nil
}

func internalResolve(envVarAccessKeyId, envVarAccessKeySecret, envVarSecurityToken string) *credentials.Config {
	accessKeyId, ok1 := os.LookupEnv(envVarAccessKeyId)
	accessKeySecret, ok2 := os.LookupEnv(envVarAccessKeySecret)
	securityToken, ok3 := os.LookupEnv(envVarSecurityToken)
	if ok1 && ok2 && ok3 {
		return &credentials.Config{
			Type:            tea.String("sts"),
			AccessKeyId:     tea.String(accessKeyId),
			AccessKeySecret: tea.String(accessKeySecret),
			SecurityToken:   tea.String(securityToken),
		}
	}
	return nil
}
