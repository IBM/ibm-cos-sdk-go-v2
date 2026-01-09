module github.com/IBM/ibm-cos-sdk-go-v2/service/s3/internal/configtesting

go 1.24.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2/config v1.29.14
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared v1.18.15
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/config => ../../../../config/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/ini => ../../../../internal/ini/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding => ../../../../service/internal/accept-encoding/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../../service/internal/presigned-url/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared => ../../../../service/internal/s3shared/
