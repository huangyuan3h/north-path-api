module north-path.it-t.xyz/auth/google

go 1.22.0

replace north-path.it-t.xyz/utils => ../../../utils

replace north-path.it-t.xyz/auth/db => ../db

replace north-path.it-t.xyz/user => ../../user

require (
	north-path.it-t.xyz/user v0.0.0-00010101000000-000000000000
	north-path.it-t.xyz/utils v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.27.0
	github.com/aws/aws-sdk-go-v2/config v1.27.12
	github.com/aws/aws-sdk-go-v2/service/s3 v1.54.2
	github.com/google/uuid v1.6.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.12 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.13.15 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.32.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.20.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.23.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.7 // indirect
	github.com/aws/smithy-go v1.20.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
)
