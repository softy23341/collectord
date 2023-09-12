package internal

import (
	"encoding/json"
	"errors"

	"git.softndit.com/collector/backend/npusher"

	"github.com/BurntSushi/toml"
	"github.com/inconshreveable/log15"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func init() {
	RegisterProvider("aws", NewAwsProvider)
}

type awsProviderKey struct {
	regular string
	sandbox string
}
type awsProviderKeysMap map[string]awsProviderKey

var awsProviderKeys = awsProviderKeysMap{
	"apns": awsProviderKey{
		regular: "APNS",
		sandbox: "APNS_SANDBOX",
	},
	"gcm": awsProviderKey{
		regular: "GCM",
		sandbox: "GCM",
	},
}

// ErrUnsupportedType TBD
var ErrUnsupportedType = errors.New("unsupported notification type")

// ErrUnregisteredType TBD
var ErrUnregisteredType = errors.New("unregistered in config notification type")

type awsConfig struct {
	Region          string            `toml:"region"`
	PlatformArns    map[string]string `toml:"platform_arns"`
	Storage         toml.Primitive    `toml:"storage"`
	AccessKeyID     string            `toml:"aws_access_key_id"`
	SecretAccessKey string            `toml:"aws_secret_access_key"`
}

func (c *awsConfig) parse(config *toml.Primitive) error {
	return toml.PrimitiveDecode(*config, c)
}

// AWSProvider TBD
type awsProvider struct {
	keys         awsProviderKeysMap
	config       *awsConfig
	snsClient    *sns.SNS
	tokenStorage tokenStorager
	log          log15.Logger
}

// NewAwsProvider TBD
func NewAwsProvider(ctx *ProviderCtx) (Provider, error) {
	config := ctx.Config

	awsConfig := &awsConfig{}
	if err := awsConfig.parse(config); err != nil {
		return nil, err
	}

	svc := sns.New(
		session.New(&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AccessKeyID,
				awsConfig.SecretAccessKey,
				"",
			)}))

	storage, err := createTokenStorage(ctx.Log, &awsConfig.Storage)

	if err != nil {
		return nil, err
	}

	ctx.Log.Debug("Build new aws provider", "arn", awsConfig)
	return &awsProvider{
		keys:         awsProviderKeys,
		log:          ctx.Log,
		config:       awsConfig,
		snsClient:    svc,
		tokenStorage: storage,
	}, nil
}

// SupportedTypes TBD
func (p *awsProvider) SupportedTypes() (types []string) {
	types = make([]string, 0, len(p.keys))
	for t := range p.keys {
		types = append(types, t)
	}
	return
}

// Send TBD
func (p *awsProvider) Send(task *npusher.NotificationTask) error {
	p.log.Debug("New push notify task", "token", task.Token)

	pushService, err := p.getPushServiceKey(task)
	if err != nil {
		return err
	}

	arn, err := p.arnByDeviceToken(task.Token, pushService)
	if err != nil {
		return err
	}

	data, err := p.buildRequest(task)
	if err != nil {
		return err
	}

	p.log.Debug("Payload", "token", task.Token, "data", data)

	resp, err := p.sendRequest(arn, data)

	if err != nil {
		p.log.Error("Fail send push", "token", task.Token, "err", err)
		return err
	}
	p.log.Debug("Successful send push", "token", task.Token, "resp", resp)

	return nil
}

func (p *awsProvider) getPushServiceKey(task *npusher.NotificationTask) (string, error) {
	pushService, ok := p.keys[task.Type]

	if !ok {
		return "", ErrUnsupportedType
	}

	var payloadKey string
	if task.Sandbox {
		payloadKey = pushService.sandbox
	} else {
		payloadKey = pushService.regular
	}

	return payloadKey, nil
}

func (p *awsProvider) arnByDeviceToken(deviceToken, pushService string) (string, error) {
	arn, err := p.getArnByDeviceToken(deviceToken, pushService)
	if err != nil {
		return "", err
	}

	arnInfo, err := p.arnInfo(arn)
	if err != nil {
		return "", err
	}

	if *arnInfo.Attributes["Enabled"] != "true" {
		err := p.setArnInfo(arn, map[string]*string{
			"Enabled": aws.String("true"),
		})
		if err != nil {
			return "", err
		}
	}

	return arn, nil
}

func (p *awsProvider) getArnByDeviceToken(deviceToken, pushService string) (string, error) {
	if arnToken, err := p.tokenStorage.get(deviceToken); arnToken != "" {
		p.log.Debug("get arn from storage",
			"targetDevice", deviceToken,
			"arnToken", arnToken)
		return arnToken, nil
	} else if err != nil {
		return "", err
	}

	arnToken, err := p.registerEndpoint(deviceToken, pushService)
	if err != nil {
		p.log.Error("Can't get arn by device",
			"token", deviceToken,
			"err", err)
		return "", err
	}

	p.tokenStorage.set(deviceToken, arnToken)

	return arnToken, nil
}

func (p *awsProvider) registerEndpoint(deviceToken, pushService string) (string, error) {
	appArn, err := p.pushServiceAppArn(pushService)
	if err != nil {
		return "", nil
	}

	params := &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(appArn),
		Token:                  aws.String(deviceToken),
	}

	resp, err := p.snsClient.CreatePlatformEndpoint(params)
	if err != nil || resp == nil {
		p.log.Error("Fail at register endpoint",
			"targetDevice", deviceToken,
			"err", err)

		return "", err
	}

	p.log.Debug("Successful register endpoint",
		"targetDevice", deviceToken,
		"resp", resp)

	return *resp.EndpointArn, nil
}

func (p *awsProvider) pushServiceAppArn(pushService string) (string, error) {
	pushServiceAppArn, ok := p.config.PlatformArns[pushService]

	if pushServiceAppArn == "" || !ok {
		return "", ErrUnregisteredType
	}
	return pushServiceAppArn, nil
}

func (p *awsProvider) arnInfo(arn string) (*sns.GetEndpointAttributesOutput, error) {
	params := &sns.GetEndpointAttributesInput{
		EndpointArn: aws.String(arn), // Required
	}
	respAttr, err := p.snsClient.GetEndpointAttributes(params)
	if err != nil || respAttr == nil {
		p.log.Error("Can't get info by arn", "arn", arn)
	}

	return respAttr, err
}

func (p *awsProvider) setArnInfo(arn string, attributes map[string]*string) error {
	p.log.Debug("set attributes", "arn", arn)

	params := &sns.SetEndpointAttributesInput{
		Attributes:  attributes,
		EndpointArn: aws.String(arn), // Required
	}
	_, err := p.snsClient.SetEndpointAttributes(params)
	if err != nil {
		p.log.Error("Can't set attributes", "arn", arn, "attr", attributes)
	}
	return err
}

func (p *awsProvider) buildRequest(task *npusher.NotificationTask) (string, error) {
	defaultMessage := task.DefaultMessage

	pushService, _ := p.getPushServiceKey(task)

	var message = map[string]interface{}{
		pushService: string(task.Payload),
		"default":   defaultMessage,
	}

	data, err := json.Marshal(message)

	return string(data), err
}

func (p *awsProvider) sendRequest(targetArn, data string) (*sns.PublishOutput, error) {
	params := &sns.PublishInput{
		Message:          aws.String(data),
		TargetArn:        aws.String(targetArn),
		MessageStructure: aws.String("json"),
	}

	resp, err := p.snsClient.Publish(params)
	return resp, err
}
