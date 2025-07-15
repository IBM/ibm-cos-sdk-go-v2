module github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources/configtesting

go 1.23.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2/config v1.29.14
	github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources v1.3.34
)

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.36.5 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/config => ../../../config/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/ini => ../../../internal/ini/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

