module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/integrationtest

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v0.0.1
	github.com/IBM/ibm-cos-sdk-go-v2/config v1.29.14
	github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager v1.17.75
	github.com/IBM/ibm-cos-sdk-go-v2/service/s3 v1.79.3
	github.com/aws/smithy-go v1.22.2
)

go 1.24.0

toolchain go1.24.4

replace github.com/IBM/ibm-cos-sdk-go-v2/service/codestar => ../../../service/codestar/

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/IBM/ibm-cos-sdk-go-v2/config => ../../../config/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager => ../../../feature/s3/manager/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/ini => ../../../internal/ini/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/v4a => ../../../internal/v4a/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/s3 => ../../../service/s3/
