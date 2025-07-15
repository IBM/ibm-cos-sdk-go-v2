module github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager

go 1.23.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.36.5
	github.com/IBM/ibm-cos-sdk-go-v2/config v1.29.14
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67
	github.com/IBM/ibm-cos-sdk-go-v2/service/s3 v1.79.3
	github.com/aws/smithy-go v1.22.2
)

require (
	github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/checksum v1.7.1 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/IBM/ibm-cos-sdk-go-v2/config => ../../../config/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/ini => ../../../internal/ini/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/v4a => ../../../internal/v4a/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/checksum => ../../../service/internal/checksum/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/s3 => ../../../service/s3/

