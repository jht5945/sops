# sops

## Important

Tested binaries:
* sops-darwin-arm64
* sops-linux-amd64

Weakness of sops, although sops encrypts boolean, for the GCM encryption method, boolean encryption is broken. 
The length of value (true or false)'s ciphertext is different.
So it is very easy to know boolean's ciphertext is true or false. 

## Build

Install `sops`:
```shell
make install
```

Build multiple platform binaries:
```shell
make build-all
```

If `keyservice.proto` is modified, `*.pb.go` files should be regenerated using command `protoc`: 
```shell
protoc \
     --go_out . --go_opt=Mkeyservice/keyservice.proto=/keyservice \
     --go-grpc_out require_unimplemented_servers=false:. --go-grpc_opt=Mkeyservice/keyservice.proto=/keyservice \
     keyservice/keyservice.proto
```
> protoc --version : libprotoc 3.21.9


## Usage

### Aliyun credential

#### AK/SK

```shell
$ export ALIBABA_CLOUD_ACCESS_KEY_ID=LT**********************
$ export ALIBABA_CLOUD_ACCESS_KEY_SECRET=P5****************************
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf --encryption-context App:App1,Env:Prod a.json
```

#### STS

```shell
$ export ALIBABA_CLOUD_ACCESS_KEY_ID=STS.*************************
$ export ALIBABA_CLOUD_ACCESS_KEY_SECRET=B1****************************
$ export ALIBABA_CLOUD_SECURITY_TOKEN=CA****************************
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf --encryption-context App:App1,Env:Prod a.json
```

OR Legacy STS environment configuration

```shell
$ export ALICLOUD_ACCESS_KEY=STS.*************************
$ export ALICLOUD_SECRET_KEY=B1****************************
$ export ALICLOUD_SECURITY_TOKEN=CA****************************
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf --encryption-context App:App1,Env:Prod a.json
```

#### Credentials file

Use credentials file set ENV `ALIBABA_CLOUD_CREDENTIALS_FILE` or default file: `~/.alibabacloud/credentials`:
```ini
[default]
type = access_key
access_key_id = LT**********************
access_key_secret = P5****************************
```

```shell
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf --encryption-context App:App1,Env:Prod a.json
```

#### ECS RAM role

Attach an instance RAM role to ECS first: https://help.aliyun.com/zh/ecs/user-guide/attach-an-instance-ram-role-to-an-ecs-instance/
```shell
$ export ALIBABA_CLOUD_ECS_METADATA=Test001hatter
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf --encryption-context App:App1,Env:Prod a.json
```

### sops encrypt

Prepare JSON file `example.json`:
```json
{
  "name": "hatter"
}
```

Encrypt via Aliyun KMS, single KMS key:
```shell
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00 example.json
{
	"name": "ENC[AES256_GCM,data:eFwcv+bL,iv:Nn4Wj7l3TyOR+/jXWQplEr3xDeshL1ZJPmQEADNzohA=,tag:wI7IzpKDvOJHxY0tFjK82A==,type:str]",
	"sops": {
		"kms": null,
		"aliyun_kms": [
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
				"created_at": "2024-04-23T01:37:13Z",
				"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVkIFyF017wMiJZbq5iXp/1ADSbdKw3fJ8eANGhFjQktTDh0gu1WpTgVNWQ/EHGY89mQ7iPnkKmj40gySWHyG0cgROrxY5eoBIQ"
			}
		],
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T01:37:14Z",
		"mac": "ENC[AES256_GCM,data:voau0ImPgSIugJGcIo0AmOPTzP5BG1LshlIRbJV2Iv7ixrf1AO9PJbhJV91nP6/mnaPmnxGGIpoLGwIE//0shluH+tBq8tRQryvAihtsXLQiQ+0OSMc0hYRHvCP0YyMy5LTWw55jVCmpJZd0WK35gfEonbfk6d3ap9InWEJKrRw=,iv:7ynvP7Opx4k0uODfRKF/RJzCbDXbT8kcjyUQzqFTp8o=,tag:nBeWY8eWkDMHfJKxUcha8Q==,type:str]",
		"pgp": null,
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

Encrypt via Aliyun KMS, multiple KMS key:
```shell
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00,acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf example.json
{
	"name": "ENC[AES256_GCM,data:FlJ0WN+k,iv:lDdsKUz2i0SKTImKIlqzdcLmtDMrx8KbF/vNBUrJoF8=,tag:xmBuSsJNwxt2p7N3tFmljA==,type:str]",
	"sops": {
		"kms": null,
		"aliyun_kms": [
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
				"created_at": "2024-04-23T01:39:42Z",
				"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVkzH55wZZeAje6HmCMvAZDwu0NKDUJERYo5oIXNzXaYGctLEK2LKC5EPc1bGfPb5+dF8PAg7fo7cMqobhZ1dmOR06fVt1tKWhj"
			},
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf",
				"created_at": "2024-04-23T01:39:42Z",
				"enc": "ZTY0NGQxYWEtZjZjNi00MzA4LThmNjctMTE5YjMzYzIxMjNi7bJn8T9b/7ZNZ3HhUNzP5+TEYcISfyg5twXkwqM1rVHfXhgectEqkAdT9xpl/DLT0siYfbMKF3//sRATTpDn38A1ykatMTVM"
			}
		],
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T01:39:43Z",
		"mac": "ENC[AES256_GCM,data:GnJEuTC3JNNXaxcfO+dUXDXKj6ZlG0BigMOmF+usPfFqptKEROyaaBG72eFx/Qs9SqEM8O6su08/nxF30mBzP/E+X7kUKioMZpfJTStfyUYnsaJUSG21/IQlwirju1SjMZwKl3Xd+iFvwXin21ytiraLd1eqmyUifHBrLzuEFB4=,iv:lIN/JCPvZf1bJv+JbwgjMeBV//CrqzALlgqHm1mDwyc=,tag:XhJsNPhVO2FK2RoX9GWmpQ==,type:str]",
		"pgp": null,
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

Encrypt via Aliyun KMS and PGP:
```shell
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00 --pgp FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4 example.json 
{
	"name": "ENC[AES256_GCM,data:jyPhZyRL,iv:RksyrpKAjLWDu/aA/Ub4OHquqXIbSwBGc/gwZJRMMmE=,tag:3WpZLybw80lUEbdZEqs0Dg==,type:str]",
	"sops": {
		"kms": null,
		"aliyun_kms": [
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
				"created_at": "2024-04-23T01:41:09Z",
				"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVkM11UYfp5UYxdZtq34XlvKsTEgTnORK/Y5PLdqLSIsPB74ivu+3zw2UEkY0bnesoHkmT2xDEddG07VO/PMwSSi8w0a27PEhxN"
			}
		],
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T01:41:09Z",
		"mac": "ENC[AES256_GCM,data:gxU+QHR8B8xVZg2sDRxuK2LmeOb0O0wKGqyPp7y7KBXOhUIN6ivMwuinGX7e84NNS4JV2S7s2mwJGW79BRvhHktyEvhNLDtiPMwmA1KYqYQUmucKEXK3HoTNv0gg7hVp/6lPKzHGvHrL2IsAyx82AFElHsIfkq/8rPrRV1XEVyY=,iv:acsw6YpLSZN+MoMnZJXguNkLXprDhfCNSfxH0Nx76BY=,tag:YDaxh+8hTaQmNcvVguPBAw==,type:str]",
		"pgp": [
			{
				"created_at": "2024-04-23T01:41:09Z",
				"enc": "-----BEGIN PGP MESSAGE-----\n\nhQEMAyUpShfNkFB/AQf/TDg1CQckoh7i9pFw/rK6H8X5uYgshvKrwt0dgUdwnqcU\nZThER1dwb9trYD2LaEXpHbOKPNjHh6wS/Fbr+Jgiy9JSpl9UkspAQbhlr2mybADS\nEPTHU5NdWMIffVMl9LdzdkiE9+HlX07CLCzEdMmNcWdjWc3/4IEROtks9I2o6kJs\nr3dEFVRVa5Nd1TdlLt+Ggv4Sn1m/Luygj2aPXGzDGKIyQJ1wnhqimj3P3lR3RDnq\nhfK4Cj2uBucJQ2f+URLhJVGPZ2n73AXrgnhsZgHzXJJXNFLmj/pxBV7jNtaI1KAR\nNpt9aEITObZRrjs8Xl6+nwCPhJV+A1MM8yxpXDdwBNJcAQfA3Mg1FamapZZ6pUTw\nR/fcao7j55Izws0oXzLRbT+nDtggG3/M2xWou3iK1vl+9/Hk4lNsqY+zU4dnEZlx\n6vUK0AMfJXE2KdIb9SbUha3iDRoDhhppJ9RigWU=\n=fcnF\n-----END PGP MESSAGE-----",
				"fp": "FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"
			}
		],
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

Encrypt via Aliyun KMS with encryption context:
```shell
$ sops encrypt --aliyun-kms acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00 --encryption-context App:Test,Env:Prod example.json
{
	"name": "ENC[AES256_GCM,data:FVRTgjK7,iv:gORMAKch9lGdRN26s7+wV0UZZNnNGCMOwUhzR+B8T4U=,tag:ZJHTWsWerAgsTdpbooSv7g==,type:str]",
	"sops": {
		"kms": null,
		"aliyun_kms": [
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
				"context": {
					"App": "Test",
					"Env": "Prod"
				},
				"created_at": "2024-04-23T02:31:20Z",
				"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVk/nBNxkwQxOYTqfS2xgMAJUOtFIH3BKVx9MdO4Z3wLoFjsVvR6ReGef/wawy54gfE8rzIz7USeQRZ3HuRfNfsmQA1yP3o7cLZ"
			}
		],
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T02:31:20Z",
		"mac": "ENC[AES256_GCM,data:p95GJtOg4TbisNT5v2KSQQxSlWYHnQTXHkyFr7iLsuzFrsxqK8yBSYpyDU/ZQ9rl8z9wStSukF6MTr93Lfz3r9ePenwVrE2txxyIT0TuSjbNGMOefFuwZL+WpEEm62PV9AtrXmDXs7ju0ZmU9vJj4Yc9JjBgwaP0dKF1ldvhNCM=,iv:vc47YC0RDsdeGKv4Ag0Jez2CTs7VU3LAvSPik8/z8X4=,tag:q7lfQd0uKEHUpGm4RIvv1w==,type:str]",
		"pgp": null,
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

### sops decrypt

```shell
$ sops decrypt example_encrypted.json 
{
	"name": "hatter"
}
```

### sops edit

```shell
$ sops example_encrypted.json
```

or

```shell
$ sops edit example_encrypted.json
```

### using config file `.sops.yaml`

Edit `.sops.yaml` file:
```yaml
creation_rules:
  - path_regex: .*prod.*
    aliyun_kms: >-
      acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00,
      acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf
    pgp: >-
      FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4
  - path_regex: .*dev.*
    aliyun_kms: acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00
```

Encrypt `app1.prod.json` (match `.*prod.*` rule):
```shell
$ sops encrypt app1.prod.json
{
	"env": "ENC[AES256_GCM,data:Upappg==,iv:3wgbLko2QPsq8Shr+4x4maTj51SrfgWm6w+2IE5KqaQ=,tag:jpncu+FSfsG8xkFAkE328w==,type:str]",
	"sops": {
		"kms": null,
		"aliyun_kms": [
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
				"created_at": "2024-04-23T02:44:45Z",
				"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVk4ldI1mlQ7kil+8bRnmTczW5V2UvZGpLtlIaruArjNyyssZZ2YZ2BGdUbrkH4W9uyH4QL28OQfi7wGg6wkyjzwHQYc4AZic+Z"
			},
			{
				"arn": " acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf",
				"created_at": "2024-04-23T02:44:45Z",
				"enc": "ZTY0NGQxYWEtZjZjNi00MzA4LThmNjctMTE5YjMzYzIxMjNiKKy+JHFcKkTjBGpd12XFQDHQA0hQdoNurga5GgG0wK1aE9igWfglRjAcv4fpluZ2SqwsHSAA/evOjKG9k05WGro7w61+BdH+"
			}
		],
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T02:44:45Z",
		"mac": "ENC[AES256_GCM,data:C//AWizm16QWMin9WkuF2R3V4D4VKzjagwC7EgfotO0qzv3EEmPaxS/976/sFTXtDRB2nmjNAEn6hU5el3ZIXrLxSFlo3eLrdcJlW0nWhAEJYjhmMRh231t35NjX7frJSB61lK0gPUqUOUPywhNhsrKX7yMeHZt3+PsQV3ztV84=,iv:xcwozziKb1Q03tdBfd301pdhdyteHGJlrs0ieG15A0Q=,tag:GdiRzAuyUzMO6SHTfidKfQ==,type:str]",
		"pgp": [
			{
				"created_at": "2024-04-23T02:44:45Z",
				"enc": "-----BEGIN PGP MESSAGE-----\n\nhQEMAyUpShfNkFB/AQf+O6umAQ5BUf2utQim7yOnzVocH2nS7j9qm2shwtlSJQr8\n5QYP6aGAR/KfhR38VaNnv8WZxnzawzEB99vcXhJj2uPGj/atD9zmBftlPmF8rgWj\n7bzp/fEx3VAT6iPUlXILfx9JyhAEp3Vtp+UttAw3j6zFnz+La2FXYvhByxg4EdTk\nzjFbbUOwWTxUYrcRwdpLt/qobWTc0lXyXxsRoJga9LrhK05fb24PMMldjWOZHVo0\nSmyt5BQ5ZIXP1yrWf5Iu5xYVUTY2zFCU6/BnNA9m8bOZ6trZyCCuNs/Mkz8okVmv\nHsD6xRxznkWYFqDqd1bMACXwcMNUwviSvR3wyaFEltJcAT8kugPkzf+LqgNVGLaH\nVUhRYpT3XovUm26w4zDugKTdeZ1IqE+OmLQyzzR8eHWHVpuAuDESHW2YI+ImOWif\ndUGpeioXPgk8k0kSuHAlq1b2t9zjXYHDOWncNww=\n=3clz\n-----END PGP MESSAGE-----",
				"fp": "FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"
			}
		],
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

Encrypt `app1.dev.json` (match `.*dev.*` rule):
```shell
sops encrypt app1.dev.json 
{
	"env": "ENC[AES256_GCM,data:N1ex+w==,iv:btgI4/VDBQiUsQFP0ILLErm7oCCDvuYl7BU8j1Q9bj8=,tag:sFf+B67KY4zV94CInMyvew==,type:str]",
	"sops": {
		"kms": null,
		"aliyun_kms": [
			{
				"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
				"created_at": "2024-04-23T02:46:21Z",
				"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVkdiygG879TZHj4O7S64RqVVtd5WCPRr26c4ygGX7gLaWofE1WH4ilUAtKEsE/grda3SyT+mfCtpdfck17kwcxk0PQ9urRY3g/"
			}
		],
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T02:46:22Z",
		"mac": "ENC[AES256_GCM,data:nhKEjWISx9K7cTCnAQMhvMlR3gW28wUELIL72CSQWmQmrfV2pUiPHsAZsrJMtdh2WQx9QU+82hYJd+ZXce9cov0/beCDW7Lq75HhqlpJg94CjDh5LgUAkkOXT6uhrei/Q7M01nf2IIzBqCKFJk3nkWX5ImPXYTqHD4Hd0CS6xF8=,iv:q3hAjfikiz61dGGJyYYRtC0ywtKX827ppPf5Vezpj7I=,tag:mpdv3unjncPCQobDNUGxvQ==,type:str]",
		"pgp": null,
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

### using config file `.sops.yaml` (Using key group)

Edit `.sops.yaml` file:
```yaml
creation_rules:
    - path_regex: .*prod.*
      shamir_threshold: 3
      key_groups:
        - aliyun_kms:
          - arn: acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00
            context:
              Env: Prod
        - aliyun_kms:
          - arn: acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf
            context:
              Env: Prod
        - pgp:
          - FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4
        - pgp:
          - D7229043384BCC60326C6FB9D8720D957C3D3074
    - path_regex: .*dev.*
      shamir_threshold: 2
      key_groups:
        - aliyun_kms:
          - arn: acs:kms:cn-hangzhou:1021806970344813:key/faf15783-0000-0000-0000-eb7a1ef56b00
            context:
              Env: Test
        - aliyun_kms:
          - arn: acs:kms:cn-hangzhou:1021806970344813:key/054a65da-0000-0000-0000-2b1ee3e071cf
            context:
              Env: Test
        - pgp:
          - FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4
        - pgp:
          - D7229043384BCC60326C6FB9D8720D957C3D3074
```

Encrypt `app1.prod.json` (match `.*prod.*` rule):
```shell
$  sops encrypt app1.prod.json
{
	"Env": "ENC[AES256_GCM,data:NcyoDg==,iv:SOg8OXmMcuAl+RnP3KJH1czX8iIsCNEm2wtnpNlFsRE=,tag:kcSxhIZ7dHMlUbFlqJf9cg==,type:str]",
	"sops": {
		"shamir_threshold": 3,
		"key_groups": [
			{
				"aliyun_kms": [
					{
						"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
						"context": {
							"Env": "Prod"
						},
						"created_at": "2024-04-23T02:59:34Z",
						"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVkm40oS9IFyxsUb+u8ECu97bxFrs5J6QJMgCBqwfe5qnpMMx1iex5GPKRhCfjmF4R+4qbuYx8u2lW+zgaHm/6LCeRFIGG/VfzU"
					}
				],
				"hc_vault": null,
				"age": null
			},
			{
				"aliyun_kms": [
					{
						"arn": "acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf",
						"context": {
							"Env": "Prod"
						},
						"created_at": "2024-04-23T02:59:34Z",
						"enc": "ZTY0NGQxYWEtZjZjNi00MzA4LThmNjctMTE5YjMzYzIxMjNiKIyWoVnKjbBk5/i6x4zmF8Lhs09WII+0zo86/kzVxm5Jefq/rBimvSF37Hak1SgQu4er4J46JIOafJmUpJFiAiMQdBLmJ6Jf"
					}
				],
				"hc_vault": null,
				"age": null
			},
			{
				"pgp": [
					{
						"created_at": "2024-04-23T02:59:34Z",
						"enc": "-----BEGIN PGP MESSAGE-----\n\nhQEMAyUpShfNkFB/AQgAoInus1gb/Qs+VjRCdu0Txkri0JrVPfs7d/MP5sIkNuxC\nXPQTnZ4wsGIDubyPwgOSwOKHPkQAOF5/0aX68isH/sV1fyMSDL5b6P3cCzqX7y+z\nCFwjlYmPLfjkCq+AegeWX/MN2GvVMkDnZbj/yRBSTUevlqHZJnJ+YQMoVDT4RJzc\neGSTKeC5wOFoOqLiWOiwjkjRfovnviLWPgjtN5gcxY+TYk2Y4wJDfEnYYT2PP7vr\nxjKmkRYFILESpjPbqeI5iawRzgEYFqN18eqxyiO4lhHN2kXKOlTUoUAKxxgfw1Qz\n+hpboxQiND6ImATanSbJkG1fd6XXPmwmDrwJ5a3NRNJdAayqy42COMB6CV53KXon\n1hIPZ4A//CuhPgPCcah7sTokY4hNZspnAJ+X5VWc6d/tx0Sy8sgpsC0DTFm2C/Fm\nhiYJ7gZwuds6kSo7/r8MHbPiuTxXKbzBGdQPiizH\n=RUdz\n-----END PGP MESSAGE-----",
						"fp": "FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"
					}
				],
				"hc_vault": null,
				"age": null
			},
			{
				"pgp": [
					{
						"created_at": "2024-04-23T02:59:34Z",
						"enc": "-----BEGIN PGP MESSAGE-----\n\nhIwDXFUltYFwV4MBA/wLNPt3B7R5PXtmxujES52EBGN1dyGkN98TrTnaBtSgWyBb\notVB6PcKv8CnZ6kKcdZ+LAK/g6PTtcmu/bf9fC7qsKJnMxJ28Ebs4yMsM7LX/HeR\nAqnpux7xFbFLno64TS4xJbm4FbnH+gwnk41UhONl2te7zWyrH+AifwcYWOotfNJf\nAVUN+jEUey8f1ZaOavcLmwbrwT05YdqYtzQPZHReA7lY4keI9Azr5XrBcBqP0zrc\nWcNlAOfKtH/+3nVsOyG90Bdsrlv5EQvOOR2NP8G1SycGOcLP4WnpZHFzcVnxoso=\n=FsPv\n-----END PGP MESSAGE-----",
						"fp": "D7229043384BCC60326C6FB9D8720D957C3D3074"
					}
				],
				"hc_vault": null,
				"age": null
			}
		],
		"kms": null,
		"aliyun_kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T02:59:35Z",
		"mac": "ENC[AES256_GCM,data:uk95rAnQR0oFlZWqKNqhxdFWH2o/EMoNkujPOin1h9TtfaJbAgajUm/KbQaV+AcLU/pBblYOf0n1Bz9XBOdOhYS0jQttJRJLWxNAFVBB1OsxV/kgrfOGwri81YZVC3FcdzpkBlGTMRdu8XPchSgDmjN5j7MQYqFWLO1t5O4sfiw=,iv:uaNiYPi5aDUa6pBjgB0QAn1jd+x5/0jA0tBn6nLryWU=,tag:G1kLKEGLxdVDO3XLB6IoXQ==,type:str]",
		"pgp": null,
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```

Encrypt `app1.dev.json` (match `.*dev.*` rule):
```shell
$ sops encrypt app1.dev.json
{
	"Env": "ENC[AES256_GCM,data:LucQ6Q==,iv:WgVM0olJCLnvU5VWyW2teCu+ryNGfa4tjm/zvAWkyTs=,tag:5ICorMdj5jQTBa+F7oUM+w==,type:str]",
	"sops": {
		"shamir_threshold": 2,
		"key_groups": [
			{
				"aliyun_kms": [
					{
						"arn": "acs:kms:cn-hangzhou:1012345678901234:key/faf15783-0000-0000-0000-eb7a1ef56b00",
						"context": {
							"Env": "Test"
						},
						"created_at": "2024-04-23T03:00:48Z",
						"enc": "MzJkMzQ1YTMtN2I2NC00MTY0LWFkZTgtNjQwZGJlODgwZjVkfH3czDEhiI+ciMIlzy3+CzaIcCn1OLugxLPK35TdP55YVBZuyJgkVDzTjb9Da7928huxHNalGCK0PmpyIYSfV3iSdgD9CMru"
					}
				],
				"hc_vault": null,
				"age": null
			},
			{
				"aliyun_kms": [
					{
						"arn": "acs:kms:cn-hangzhou:1012345678901234:key/054a65da-0000-0000-0000-2b1ee3e071cf",
						"context": {
							"Env": "Test"
						},
						"created_at": "2024-04-23T03:00:48Z",
						"enc": "ZTY0NGQxYWEtZjZjNi00MzA4LThmNjctMTE5YjMzYzIxMjNi/fytFxz64hq8SdTcANizTC23hZxR+rpX2vEh3/8VkHuFuFBXXd95j/8XvkLTFjtepj/jy6O0hYyBW4Y4E3T5Kptl7sSlL7lX"
					}
				],
				"hc_vault": null,
				"age": null
			},
			{
				"pgp": [
					{
						"created_at": "2024-04-23T03:00:48Z",
						"enc": "-----BEGIN PGP MESSAGE-----\n\nhQEMAyUpShfNkFB/AQf7Bi3WU2URBvMxlBz9ypfsNU3Dg53XK3y5ciBwrTMKM+iT\nGy1TGAUAO6JLfP1T4lKwL/uv7ecs0s0knQoP0H3MjgfUlOjzhPJ3NgtKOKNUw8hM\n34Mo80m4kzdmrfNF3KZrkumQ8F+iyoilWRN/EXfHQLcjiwk5dj60p0d3Z/S/O/jW\nEJ4cJeJMVbAmyFMOJkOsMLyk32i6Q1H6Ukz9G6LPbijbrjS+BGADy+fZZtYAUwnV\nGtEKDZMLv5ztUjJRPio/xgJn4wclldkM2j4T1A1vkBGvW5y+NRETEj3R0X9SNk+3\nMltLy5qOPolk3xm3nXbeGqFqjOnQ5XZ7HxY4VpUv9dJfAVUuVE/a6/N8sOvEFImo\nzi2ba1D5PSbRiloSYvjlwREiE31PPVLi58Wp53dCjKkqQEhP4JQFWEZBKTXtzTFh\no76WuriUYsiA/D4PYPar/7ptUCHpQAUcYLCS7rbnxnQ=\n=9eMC\n-----END PGP MESSAGE-----",
						"fp": "FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"
					}
				],
				"hc_vault": null,
				"age": null
			},
			{
				"pgp": [
					{
						"created_at": "2024-04-23T03:00:48Z",
						"enc": "-----BEGIN PGP MESSAGE-----\n\nhIwDXFUltYFwV4MBBACJ27DQe3aP9BD120MMaS1TaTi729Ki3kd/bfnBvZ0AHaAD\ndRQh6RnC7SZ8CzvMHiFapBqEpvbXEtiXiEc/O9E+uyMerx2NEIQXu49F+Fl7j6aF\n+5dV0kb641Vnzg+pmKDuoR+jg4ODbex/Or5qNaV2Vd5oi/IzldLICnz2a1eZ1dJf\nAe//56XMmAQeuSci7YXclufjK8SIjndq1+n6QcJuneq3DSSWZJs06oIaewcCFUc9\nvsADFsAnoXQ0BhyRHH+NZ2AI863B166EIQqFDGXBV8t/DRI/1jybjN4fXVCKZAs=\n=ytLO\n-----END PGP MESSAGE-----",
						"fp": "D7229043384BCC60326C6FB9D8720D957C3D3074"
					}
				],
				"hc_vault": null,
				"age": null
			}
		],
		"kms": null,
		"aliyun_kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": null,
		"lastmodified": "2024-04-23T03:00:49Z",
		"mac": "ENC[AES256_GCM,data:30VbA0eyWIVsOi6b5rUutHVYcykcHBEame1m4woNUI2nEyyKuHAssPfriTGvor6Io14xAfVrcbtkKF67FBBnreuCQwauUXQuxim2qcNioisjqA1nZKZLHbiSchA1MQPdOOolh81ekcfvy0lq4NRaiouSYR4vSQHmYZamvYJMxyY=,iv:BTXP/4XEtsZAgSctTt2xOxXT6rfeJ0fTt5WLj+9XCs8=,tag:RzBEVVF8ocNbS7tCpctwjQ==,type:str]",
		"pgp": null,
		"unencrypted_suffix": "_unencrypted",
		"version": "3.8.1-alibaba-cloud-kms-r1"
	}
}
```
