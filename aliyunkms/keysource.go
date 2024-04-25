package aliyunkms

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	kms20160120 "github.com/alibabacloud-go/kms-20160120/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/getsops/sops/v3/logging"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	// ENVKmsSopsConfigFile is a JSON file config for credentials.Config
	ENVKmsSopsConfigFile = "ALIBABA_CLOUD_KMS_SOPS_CONFIG_FILE"

	ENVKmsSopsRegionId = "ALIBABACLOUD_KMS_SOPS_REGION_ID"
	ENVRegionId        = "ALIBABACLOUD_REGION_ID"

	ENVKmsSopsEndPoint = "ALIBABACLOUD_KMS_SOPS_ENDPOINT"

	// KeyTypeIdentifier is the string used to identify a Aliyun KMS MasterKey.
	KeyTypeIdentifier = "aliyun_kms"

	// RotationTtl is the duration after which a MasterKey requires rotation.
	RotationTtl = time.Hour * 24 * 30 * 6
)

var (
	// Aliyun KMS Encryption error
	encryptDataKeyError = errors.New("EncryptDataKeyError")
	// Aliyun KMS Decryption error
	decryptDataKeyError = errors.New("DecryptDataKeyError")

	// log is the global logger for any AlibabaCloud KMS MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("ALIYUNKMS")
}

// MasterKey is an Alibaba Cloud KMS key used to encrypt and decrypt SOPS' data key using
// Alibaba Cloud SDK for Go.
type MasterKey struct {
	// Arn associated with the Alibaba Cloud KMS key.
	Arn string
	// EncryptedKey stores the data key in it's encrypted form.
	EncryptedKey string
	// CreationDate is when this MasterKey was created.
	CreationDate time.Time
	// Reference: https://help.aliyun.com/zh/kms/support/encryptioncontext
	EncryptionContext map[string]*string
}

// NewMasterKey creates a new MasterKey from an ARN, setting
// the creation date to the current date.
func NewMasterKey(arn string, context map[string]*string) *MasterKey {
	return &MasterKey{
		Arn:               arn,
		EncryptionContext: context,
		CreationDate:      time.Now().UTC(),
	}
}

func MasterKeysFromArnString(arn string, context map[string]*string) []*MasterKey {
	var keys []*MasterKey
	if arn == "" {
		return keys
	}
	for _, k := range strings.Split(arn, ",") {
		keys = append(keys, NewMasterKey(k, context))
	}
	return keys
}

func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

func (key *MasterKey) Encrypt(dataKey []byte) error {
	client, err := key.createClient()
	if err != nil {
		return err
	}
	encryptRequest := &kms20160120.EncryptRequest{
		KeyId:     tea.String(key.Arn),
		Plaintext: tea.String(base64.StdEncoding.EncodeToString(dataKey)),
	}
	if len(key.EncryptionContext) > 0 {
		encryptRequest.SetEncryptionContext(key.getEncryptionContext())
	}
	encryptResponse, err := client.Encrypt(encryptRequest)
	if err != nil {
		return err
	}
	if *encryptResponse.StatusCode != 200 {
		log.WithField("arn", key.Arn).Info("failed to encrypt via aliyunkms, status code: %d, response: %s", *encryptResponse.StatusCode, encryptResponse.String())
		return encryptDataKeyError
	}
	encryptedDataKey := encryptResponse.Body.CiphertextBlob
	key.SetEncryptedDataKey([]byte(*encryptedDataKey))
	return nil
}

func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

func (key *MasterKey) Decrypt() ([]byte, error) {
	client, err := key.createClient()
	if err != nil {
		return nil, err
	}
	decryptRequest := &kms20160120.DecryptRequest{
		CiphertextBlob: tea.String(key.EncryptedKey),
	}
	if len(key.EncryptionContext) > 0 {
		decryptRequest.SetEncryptionContext(key.getEncryptionContext())
	}
	decryptResponse, err := client.Decrypt(decryptRequest)
	if err != nil {
		return nil, err
	}
	if *decryptResponse.StatusCode != 200 {
		log.WithField("arn", key.Arn).Info("failed to decrypt via aliyunkms, status code: %d, response: %s", *decryptResponse.StatusCode, decryptResponse.String())
		return nil, decryptDataKeyError
	}
	return base64.StdEncoding.DecodeString(*decryptResponse.Body.Plaintext)
}

func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > RotationTtl
}

func (key *MasterKey) ToString() string {
	return key.Arn
}

func (key *MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["arn"] = key.Arn
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	if len(key.EncryptionContext) > 0 {
		out["context"] = key.getEncryptionContext()
	}
	return out
}

func (key *MasterKey) createClient() (*kms20160120.Client, error) {
	var credentialConfig *credentials.Config = nil
	kmsSopsConfigFile, ok := os.LookupEnv(ENVKmsSopsConfigFile)
	if ok {
		// if ALIBABA_CLOUD_KMS_SOPS_CONFIG_FILE is configured in ENV, load config file
		log.Infof("found kms sops config file: %s", kmsSopsConfigFile)
		kmsSopsConfigContent, err := os.ReadFile(kmsSopsConfigFile)
		if err != nil {
			log.WithField("arn", key.Arn).Info("read kms sops config file: %s, error: %s", kmsSopsConfigFile, err.Error())
			return nil, err
		}
		var tempCredentialConfig credentials.Config
		unmarshalErr := json.Unmarshal(kmsSopsConfigContent, &tempCredentialConfig)
		if unmarshalErr != nil {
			log.WithField("arn", key.Arn).Info("parse kms sops config file: %s, error: %s", kmsSopsConfigFile, unmarshalErr.Error())
			return nil, unmarshalErr
		}
	}
	// if ALIBABA_CLOUD_KMS_SOPS_CONFIG_FILE is NOT configured in ENV:
	// lookup credential provider by default chain:
	// []Provider{providerEnv, providerOIDC, providerProfile, providerInstance}
	// Reference: https://github.com/aliyun/credentials-go#credential-provider-chain
	return key.createClientWithConfig(credentialConfig)
}

func (key *MasterKey) createClientWithConfig(credentialConfig *credentials.Config) (*kms20160120.Client, error) {
	credentialProvider, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		log.WithField("arn", key.Arn).Info("new aliyun kms credential error", err.Error())
		return nil, err
	}

	regionId := key.getRegionId()
	endPoint, _ := os.LookupEnv(ENVKmsSopsEndPoint)

	openapiConfig := &openapi.Config{
		Credential: credentialProvider,
	}
	if regionId != "" {
		openapiConfig.RegionId = tea.String(regionId)
	}
	if endPoint != "" {
		openapiConfig.Endpoint = tea.String(endPoint)
	}
	client, err := kms20160120.NewClient(openapiConfig)
	if err != nil {
		log.WithField("arn", key.Arn).Info("create aliyun kms client error: %s", err.Error())
	}
	return client, err
}

func (key *MasterKey) getRegionId() string {
	regionId := parseRegionIdFromArn(key.Arn)
	if regionId == "" {
		if kmsOpsRegionId, ok := os.LookupEnv(ENVKmsSopsRegionId); ok {
			regionId = kmsOpsRegionId
		}
	}
	if regionId == "" {
		if aliRegionId, ok := os.LookupEnv(ENVRegionId); ok {
			regionId = aliRegionId
		}
	}
	return regionId
}

func (key *MasterKey) getEncryptionContext() map[string]interface{} {
	encryptionContext := make(map[string]interface{})
	for k, v := range key.EncryptionContext {
		encryptionContext[k] = v
	}
	return encryptionContext
}

func parseRegionIdFromArn(keyId string) string {
	// ARN: acs:kms:cn-hangzhou:1192853035118460:key/key-hzz64a3dbd1prbfsnnvpe
	kmsArnPattern := regexp.MustCompile("acs:kms:([\\w-]+):\\d*:.*")
	matchGroups := kmsArnPattern.FindStringSubmatch(keyId)
	if len(matchGroups) > 1 {
		return matchGroups[1]
	} else {
		return ""
	}
}
